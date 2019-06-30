#!/usr/bin/env bash
dep ensure -add golang.org/x/net/context \
    github.com/envoyproxy/go-control-plane/envoy/api/v2 \
    github.com/envoyproxy/go-control-plane/envoy/api/v2 \
    github.com/envoyproxy/go-control-plane/envoy/api/v2/auth \
    github.com/envoyproxy/go-control-plane/envoy/api/v2/core \
    github.com/envoyproxy/go-control-plane/envoy/api/v2/listener \
    github.com/envoyproxy/go-control-plane/envoy/api/v2/route \
    github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2 \
    github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2 \
    github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2 \
    github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2 \
    github.com/envoyproxy/go-control-plane/pkg/cache \
    github.com/envoyproxy/go-control-plane/pkg/server \
    github.com/envoyproxy/go-control-plane/pkg/util \
    github.com/sirupsen/logrus \