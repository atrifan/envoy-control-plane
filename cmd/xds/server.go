package xds

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	v1Handler "github.com/atrifan/envoy-plane/pkg/api/handler/rest/v1"
	xdshandler "github.com/atrifan/envoy-plane/pkg/api/handler/xds"
	v1 "github.com/atrifan/envoy-plane/pkg/api/service/v1"
	myals "github.com/atrifan/envoy-plane/pkg/api/util/acesslogs"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	accesslog "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	debug       bool
	onlyLogging bool

	localhost = "127.0.0.1"

	port        uint
	gatewayPort uint
	alsPort     uint
	httpRestPort    uint
	grpcRestPort	uint

	mode string

	version int32

	config cache.SnapshotCache
)

const (
	XdsCluster = "xds_cluster"
	Ads        = "ads"
	Xds        = "xds"
	Rest       = "rest"
)

func init() {
	flag.BoolVar(&debug, "debug", true, "Use debug logging")
	flag.BoolVar(&onlyLogging, "onlyLogging", false, "Only demo AccessLogging Service")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&gatewayPort, "gateway", 18001, "Management server port for HTTP gateway")
	flag.UintVar(&httpRestPort, "rest-port", 8082, "HTTP rest port to bind")
	flag.UintVar(&grpcRestPort, "grpc-rest-port", 8081, "gRPC rest port to bind")
	flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "ads", Ads, "Management server type (ads, xds, rest)")
}

// RunAccessLogServer starts an accesslog service.
func RunAccessLogServer(ctx context.Context, als *myals.AccessLogService, port uint) {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	accesslog.RegisterAccessLogServiceServer(grpcServer, als)
	log.WithFields(log.Fields{"port": port}).Info("access log server listening")

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

const grpcMaxConcurrentStreams = 1000000

// RunManagementServer starts an xDS server at the given port.
func RunManagementServer(ctx context.Context, server xds.Server, port uint) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	// register services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	log.WithFields(log.Fields{"port": port}).Info("management server listening")
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

// RunManagementGateway starts an HTTP gateway to an xDS server.
func RunManagementGateway(ctx context.Context, srv xds.Server, port uint) {
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &xds.HTTPGateway{Server: srv}}
	log.WithFields(log.Fields{"port": port}).Info("gateway listening HTTP/1.1")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	<- ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		log.Error(err)
	}
}

func RunRestServicesGrpc(ctx context.Context, grpcPort uint, v1API v1Handler.ClusterServiceServer) {
	server := grpc.NewServer()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d",grpcPort))
	if err != nil {
		log.Fatalf("Error grpc listenting %s", err)
		return
	}

	// register service
	v1Handler.RegisterClusterServiceServer(server, v1API)
	log.WithFields(log.Fields{"port": grpcPort}).Info("started grpc rest server")

	// graceful shutdown
	go func() {
		if err := server.Serve(listen); err != nil {
			log.Fatalf("Grpc rest server failed: %s", err)
		}
	}()

	<- ctx.Done()

	server.GracefulStop()
	log.Println("stoped gRPC server...")
}

func RunRestServicesHttp(ctx context.Context, grpcPort uint, httpPort uint) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1Handler.RegisterClusterServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d",grpcPort), opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
		return
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mux,
	}

	log.WithFields(log.Fields{"port": httpPort}).Info("started http rest server")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to start HTTP gateway ListenAndServe(): %s", err)
		}
	}()

	<- ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown rest http: %s", err)
	}
	log.Fatalf("Close http rest server")
}

func InitXds(ctx context.Context) {
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	signal := make(chan struct{})
	cb := &xdshandler.Callbacks{
		Signal:   signal,
		Fetches:  0,
		Requests: 0,
	}
	config = cache.NewSnapshotCache(mode == Ads, xdshandler.Hasher{}, xdshandler.Logger{})

	srv := xds.NewServer(config, cb)

	//als := &accesslogs.AccessLogService{}
	als := &myals.AccessLogService{}
	go RunAccessLogServer(ctx, als, alsPort)

	if onlyLogging {
		cc := make(chan struct{})
		<-cc
		os.Exit(0)
	}

	// start the xDS server
	go RunManagementServer(ctx, srv, port)
	go RunManagementGateway(ctx, srv, gatewayPort)

	//start rest server
	v1API := v1.NewToDoServiceServer(config)
	go RunRestServicesGrpc(ctx, grpcRestPort, v1API)
	go RunRestServicesHttp(ctx, grpcRestPort, httpRestPort)

	<-signal

	als.Dump(func(s string) { log.Debug(s) })
	cb.Report()

	_cacheInit()

}

func _cacheInit() {
	for {
		atomic.AddInt32(&version, 1)
		nodeId := config.GetStatusKeys()[0]

		var clusterName = "service_bbc"
		var remoteHost = "www.bbc.com"
		var sni = "www.bbc.com"
		log.Infof(">>>>>>>>>>>>>>>>>>> creating cluster " + clusterName)

		//c := []cache.Resource{resource.MakeCluster(resource.Ads, clusterName)}

		h := &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address:  remoteHost,
				Protocol: core.TCP,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: uint32(443),
				},
			},
		}}

		c := []cache.Resource{
			&v2.Cluster{
				Name:            clusterName,
				ConnectTimeout:  2 * time.Second,
				ClusterDiscoveryType: &v2.Cluster_Type{v2.Cluster_LOGICAL_DNS},
				DnsLookupFamily: v2.Cluster_V4_ONLY,
				LbPolicy:        v2.Cluster_ROUND_ROBIN,
				Hosts:           []*core.Address{h},
				TlsContext: &auth.UpstreamTlsContext{
					Sni: sni,
				},
			},
		}

		// =================================================================================
		var listenerName = "listener_0"
		var targetHost = "www.bbc.com"
		var targetRegex = "/api"
		var virtualHostName = "local_service"
		var routeConfigName = "local_route"

		log.Infof(">>>>>>>>>>>>>>>>>>> creating listener " + listenerName)

		v := route.VirtualHost{
			Name:    virtualHostName,
			Domains: []string{"*"},

			Routes: []route.Route{{
				Match: route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: targetRegex,
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						HostRewriteSpecifier: &route.RouteAction_HostRewrite{
							HostRewrite: targetHost,
						},
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
						PrefixRewrite: "/robots.txt",
					},
				},
			}}}

		manager := &hcm.HttpConnectionManager{
			CodecType:  hcm.AUTO,
			StatPrefix: "ingress_http",
			RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
				RouteConfig: &v2.RouteConfiguration{
					Name:         routeConfigName,
					VirtualHosts: []route.VirtualHost{v},
				},
			},
			HttpFilters: []*hcm.HttpFilter{{
				Name: util.Router,
			}},
		}
		pbst, err := util.MessageToStruct(manager)
		if err != nil {
			panic(err)
		}

		var l = []cache.Resource{
			&v2.Listener{
				Name: listenerName,
				Address: core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Protocol: core.TCP,
							Address:  localhost,
							PortSpecifier: &core.SocketAddress_PortValue{
								PortValue: 9191,
							},
						},
					},
				},
				FilterChains: []listener.FilterChain{{
					Filters: []listener.Filter{{
						Name:   util.HTTPConnectionManager,
						ConfigType: &listener.Filter_Config{pbst},
					}},
				}},
			}}

		// =================================================================================

		log.Infof(">>>>>>>>>>>>>>>>>>> creating snapshot Version " + fmt.Sprint(version))
		snap := cache.NewSnapshot(fmt.Sprint(version), nil, c, nil, l)


		config.SetSnapshot(nodeId, snap)

		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')

	}
}

