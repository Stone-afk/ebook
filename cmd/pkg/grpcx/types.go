package grpcx

import "google.golang.org/grpc"

type Register interface {
	Registry(server *grpc.Server)
}
