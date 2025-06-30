package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/api"
	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	grpcPort = os.Getenv("GRPC_PORT")
	promPort = os.Getenv("PROM_PORT")
)

func main() {
	if grpcPort == "" {
		log.Fatal("GRPC_PORT environment variable must me set to run the server.")
	}

	allEmoji := emoji.NewAllEmoji()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatal(err.Error())
	}

	errs := make(chan error, 1)

	if promPort != "" {
		// Start prometheus server
		go func() {
			log.Printf("Starting prom metrics on PROM_PORT=[%s]", promPort)
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(fmt.Sprintf(":%s", promPort), nil)
			errs <- err
		}()
	}

	// Start grpc server
	go func() {
		grpc_prometheus.EnableHandlingTimeHistogram()
		grpcServer := grpc.NewServer(
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		)
		api.NewGrpServer(grpcServer, allEmoji)
		log.Printf("Starting grpc server on GRPC_PORT=[%s]", grpcPort)
		err := grpcServer.Serve(lis)
		errs <- err
	}()

	// Catch shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT)
		s := <-sig
		errs <- fmt.Errorf("caught signal %v", s)
	}()

	log.Fatal(<-errs)
}
