package main

import (
	"context"
	"github.com/ducknightii/grpc-demo/pb"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
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
	pb.RegisterGreeterServer(s, &server{})
	go startTrace()
	log.Println("server start...")
	err = s.Serve(lis)
	panic(err)
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (res *pb.HelloReply, err error) {

	p, ok := peer.FromContext(ctx)

	res = new(pb.HelloReply)
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
	log.Println("trace start...")
	grpclog.Infoln("Trace listen on 30001")

}
