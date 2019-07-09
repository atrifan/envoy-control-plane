package v1

import (
	"context"
	v1 "github.com/atrifan/envoy-plane/pkg/api/handler/rest/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {}

func (v1 *toDoServiceServer) ReadAll(context.Context, *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	panic("implement me")
}

func (v1 *toDoServiceServer) Create(context.Context, *v1.CreateRequest) (*v1.CreateResponse, error) {
	panic("implement me")
}

func (v1 *toDoServiceServer) Read(context.Context, *v1.ReadRequest) (*v1.ReadResponse, error) {
	panic("implement me")
}

func (v1 *toDoServiceServer) Update(context.Context, *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	panic("implement me")
}

func (v1 *toDoServiceServer) Delete(context.Context, *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	panic("implement me")
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer() v1.ToDoServiceServer {
	return &toDoServiceServer{}
}
