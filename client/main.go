package main

import (
	"context"
	"fmt"
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"time"
)

func main(){
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err!=nil{
		panic(err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)


	//method 1
	//req, err := client.SayHi(context.Background(), &pb.SearchRequest{
	//	Request: "scccc",
	//})
	//
	//response := req.GetResponse()
	//log.Println(response)

	//method 2
	//c,err:=client.SayHi1(context.Background())
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//for i:=0;i<10;i++{
	//	c.Send(&pb.SearchRequest{Request: strconv.Itoa(i)})
	//	time.Sleep(1*time.Second)
	//}
	//recv, err := c.CloseAndRecv()
	//fmt.Println(recv)


	//method 3
	//out, err := client.SayHi2(context.Background(), &pb.SearchRequest{Request: "shuchang"})
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//for{
	//	recv, err := out.Recv()
	//	if err!=nil{
	//		fmt.Println(err)
	//		break
	//	}
	//	fmt.Println(recv.Response)
	//}

	//method4

	c, err := client.SayHi3(context.Background())
	if err!=nil{
		fmt.Println(err)
	}
	go func() {
		for{
			recv, err := c.Recv()
			if  err!=nil{
				fmt.Println(err)
				break
			}
			fmt.Println(recv.Response)
		}
	}()
	go func() {
		for i:=0;i<10;i++{
			err := c.Send(&pb.SearchRequest{Request: fmt.Sprintf("shuchang%d", i)})
			time.Sleep(time.Second*1)
			if err!=nil{
				break
			}
		}
	}()

	time.Sleep(20*time.Second)
}
