package main

import (
	"bufio"
	"errors"
	"fmt"
	"goim-pro/internal/app/grpc"
	"goim-pro/pkg/logs"
	"os"
	"os/signal"
	"syscall"
)

var logger = logs.GetLogger("INFO")
var server *grpc.GRPCServer

func main() {
	server = grpc.NewServer()
	server.InitServer()
	server.StartGRPCServer()

	defer func() {
		if r := recover(); r != nil {
			err := errors.New(fmt.Sprint(r))
			logger.Panicf("uncheck exception: %v", err)
		}
	}()

	exitChan := make(chan int)
	go signalHandler(exitChan)
	go commandHandler(exitChan)

	code := <-exitChan
	server.GracefulStopGRPCServer()
	os.Exit(code)
}

func signalHandler(exitChan chan int) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			exitChan <- 0
		default:
			logger.Debug("Unknown signal.")
			exitChan <- 1
		}
	}
}

func commandHandler(exitChan chan int) {
	for {
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		str := string(char)
		switch str {
		case "q":
			logger.Infoln("server is starting to disconnect...")
			server.GracefulStopGRPCServer()
			logger.Infoln("server has been gracefully disconnected!")
			exitChan <- 0
		default:
			logger.Info("server continue to listen...")
		}
	}
}
