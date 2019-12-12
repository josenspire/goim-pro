package main

import "fmt"

func main() {
	fmt.Print("Hello Go IM PRO!!")

	s := grpc.NewServer()
	lis, _ := net.Listen("tcp", "localhost:50051")
	// error handling omitted
	s.Serve(lis)
}
