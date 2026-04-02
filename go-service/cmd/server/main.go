package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Dev66-66/LB10/go-service/internal/app"
	grpcserver "github.com/Dev66-66/LB10/go-service/internal/grpc"
	"github.com/Dev66-66/LB10/go-service/internal/grpc/pb"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

func main() {
	sharedStore := store.NewWorkoutStore()

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("msg=failed to listen addr=:50051 err=%v", err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterWorkoutServiceServer(grpcSrv, grpcserver.NewWorkoutGRPCServer(sharedStore))

	go func() {
		log.Println("msg=starting gRPC server addr=:50051")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("msg=gRPC server failed err=%v", err)
		}
	}()

	// HTTP server (blocking)
	log.Println("msg=starting HTTP server addr=:8080")
	httpApp := app.New(sharedStore)
	if err := httpApp.Run(":8080"); err != nil {
		log.Fatalf("msg=HTTP server failed err=%v", err)
	}
}
