package main

import (
	"bufio"
	"errors"
	"fmt"
	"goim-pro/internal/app/grpc"
	"goim-pro/pkg/logs"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

var logger = logs.GetLogger("INFO")
var server *grpc.GRPCServer

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := errors.New(fmt.Sprint(r))
			logger.Panicf("uncheck exception: %v", err)
		}
	}()

	cpuProfile := os.Getenv("cpu_profile")
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			logger.Fatalf("Failed to create CPU profiling file due to error - %s", err.Error())
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	server = grpc.NewServer()
	server.InitServer()
	server.StartGRPCServer()

	exitChan := make(chan int)
	go signalHandler(exitChan)
	// go commandHandler(exitChan)

	code := <-exitChan
	server.GracefulStopGRPCServer()

	memProfile := os.Getenv("mem_profile")
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			logger.Fatalf("Failed to create memory profiling file due to error - %s", err.Error())
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			logger.Fatalf("Failed to write memory profiling data to file due to error - %s", err.Error())
		}
		_ = f.Close()
	}

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
