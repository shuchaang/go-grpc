package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server)SayHi(ctx context.Context,r *pb.SearchRequest)(*pb.SearchResponse,error){
	fmt.Println("收到消息:"+r.Request)
	return &pb.SearchResponse{
		Response: r.GetRequest()+"Server",
	},nil
}
func (s *server)SayHi1(server pb.HelloService_SayHi1Server)error{
	for{
		req, err := server.Recv()
		if err!=nil{
			server.SendAndClose(&pb.SearchResponse{
				Response: "close",
			})
			break
		}
		fmt.Println(req.Request)
	}
	return nil
}

func (s *server) SayHi2(req *pb.SearchRequest, server pb.HelloService_SayHi2Server) error {
	name:=req.Request
	for i:=0;i<10;i++{
		time.Sleep(1*time.Second)
		server.Send(&pb.SearchResponse{Response: fmt.Sprintf("%s%d", name, i)})
	}
	return nil
}

func (s *server) SayHi3(server pb.HelloService_SayHi3Server) error {
	data := make(chan string)
	go func() {
		for  {
			recv, err := server.Recv()
			if recv==nil{
				continue
			}
			if err!=nil{
				//这里注意 不放信息会死锁
				fmt.Println("rev error ",err)
				data<-"EOF"
			}
			data<-recv.Request
		}
	}()
	for{
		req:= <-data
		fmt.Println("rev:"+req)
		if req=="EOF"{
			break
		}
		server.Send(&pb.SearchResponse{Response: "rev:" + req})
	}
	return nil
}


func main(){
	done:=make(chan int)
	go regGrpc()
	go httpGtw()
	<-done
}

func regGrpc(){
	listen, err := net.Listen("tcp", ":50051")
	if err!=nil{
		panic(err)
	}
	s:= grpc.NewServer()
	pb.RegisterHelloServiceServer(s,&server{})
	log.Printf("server listen:%s",listen.Addr())
	if err=s.Serve(listen);err!=nil{
		log.Fatalf("start server error %v",err)
	}
}

func httpGtw(){
	fmt.Println("启动gtw")
	conn, err := grpc.DialContext(context.Background(),
		"127.0.0.1:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err!=nil{
		panic(err)
	}

	mux := runtime.NewServeMux()
	s := &http.Server{
		Handler: mux,
		Addr:    ":8088",
	}
	err = pb.RegisterHelloServiceHandler(context.Background(), mux, conn)
	if err!=nil{
		panic(err)
	}

	err = s.ListenAndServe()
	fmt.Println(err)
}