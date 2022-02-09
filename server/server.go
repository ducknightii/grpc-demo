package main

import (
	"context"
	"github.com/ducknightii/grpc-demo/pb/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", ":30000")
	if err != nil {
		log.Fatal("tcp listen err: ", err)
	}

	grpc.EnableTracing = true
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	go startTrace()
	errCh := make(chan error)
	go func() {
		// grpc
		log.Println("server start  on http://0.0.0.0:30000")
		errCh <- s.Serve(lis)
	}()
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:30000",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = helloworld.RegisterGreeterHandler(context.Background(), gwmux, conn)
	gwServer := &http.Server{
		Addr:    ":30003",
		Handler: gwmux,
	}

	go func() {
		log.Println("Serving gRPC-Gateway on http://0.0.0.0:30003")
		errCh <- gwServer.ListenAndServe()
	}()

	select {
	case err = <-errCh:
		panic(err)
	}
}

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (res *helloworld.HelloReply, err error) {

	p, ok := peer.FromContext(ctx)

	res = new(helloworld.HelloReply)
	res.Message = "Hello " + req.Name
	gap := 20 * time.Duration(time.Now().Nanosecond()/(1000*1000*10))
	time.Sleep(gap * time.Millisecond)
	if ok {
		log.Printf("request info: %s, sleep: %dms\n", p.Addr, gap)
	} else {
		log.Println("peer>FromContext empty")
	}

	return
}

func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	go http.ListenAndServe(":30001", nil)
	log.Println("trace start on http://0.0.0.0:30001")
}
