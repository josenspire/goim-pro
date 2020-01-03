package main

import (
	"bufio"
	"flag"
	"fmt"
	"goim-pro/internal/app"
	"os"
)

func main() {
	flag.Parse()

	server := app.NewServer()
	server.StartGrpcServer()

	for {
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		switch char {
		case 'q':
			fmt.Println("server disconnecting...")
			server.GracefulStopGrpcServer()
		default:
			fmt.Println("server continue to listen...")
		}
	}
}
