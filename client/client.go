package main

import (
	"context"
	"fmt"
	"github.com/ducknightii/grpc-demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

var client pb.GreeterClient

func main() {
	conn, err := grpc.Dial("127.0.0.1:30000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("dail err: ", err)
	}
	defer conn.Close()

	client = pb.NewGreeterClient(conn)

	go func() {
		// /debug/pprof/
		http.ListenAndServe(":30002", nil)
	}()

	time.Sleep(time.Second * 300)

	fmt.Println("request start...")
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sayHello(fmt.Sprintf("%s-%d", "World", i))
		}(i)
		time.Sleep(time.Second * 10)
	}
	wg.Wait()

	time.Sleep(time.Second * 200)

	return
}

func sayHello(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Printf("time: %s, client request: %s\n", time.Now(), name)

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Printf("time: %s, grpc.SayHello err: %+v\n", time.Now(), err)
		return
	}

	log.Printf("time: %s, client receive: %s\n", time.Now(), r.Message)
}
