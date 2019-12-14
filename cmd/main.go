package main

import (
	"bufio"
	"flag"
	"fmt"
	"goim-pro/internal/app/grpc"
	"os"
)

func main() {
	flag.Parse()

	_ = grpc.New()

MetaExist:
	for {
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		switch char {
		case 'q':
			fmt.Println("server disconnecting...")
			break MetaExist
		default:
			fmt.Println("server continue to listen...")
		}
	}

}
