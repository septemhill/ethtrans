package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

const grpcPort = 1117

type GRPCServer struct {
	server *grpc.Server
}

func (g *GRPCServer) Start() {
	go func() {
		listen, _ := net.Listen("tcp", ":"+strconv.FormatInt(grpcPort, 10))

		reflection.Register(g.server)

		if err := g.server.Serve(listen); err != nil {
			log.Fatal("server down")
		}
	}()
}

func (g *GRPCServer) Stop() {
	g.server.Stop()
}

func NewGRPCServer() *GRPCServer {
	g := &GRPCServer{
		server: grpc.NewServer(),
	}

	return g
}
