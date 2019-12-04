/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"math/rand"
	"time"
	"strconv"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"net/http"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"contrib.go.opencensus.io/exporter/prometheus"
	"io/ioutil"
	"path/filepath"
  "gopkg.in/yaml.v2"
)

const (
	//port = ":50051"
	configFile    = "../config.yaml"
)


// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

//global variables
var waitMean, waitStd, waits float64


func getConfig(configFile string) map[string]interface{}{
	log.Printf("Reading...%v\n", configFile)
	var config map[string]interface{}
	filename, _ := filepath.Abs("../config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatalf("Error: %v\n",err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error: %v\n",err)
	}
	return config
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	waits := int((rand.NormFloat64() * waitStd + waitMean) * 1000)
	log.Printf("Received: %v (%v)", in.GetName(), strconv.Itoa(waits))
	time.Sleep( time.Duration(waits) * time.Millisecond )
	return &pb.HelloReply{Message: "Hello " + in.GetName()+ " , " + strconv.Itoa(waits)}, nil
}

func main() {
	log.Printf("Starting...")
	config := getConfig(configFile)
	port := config["server"].(map[interface{}]interface{})["port"].(int)
	waitMean = config["server"].(map[interface{}]interface{})["waitMean"].(float64)
	waitStd = config["server"].(map[interface{}]interface{})["waitStd"].(float64)
	exporterPort := config["server"].(map[interface{}]interface{})["exporterPort"].(int)

	lis, err := net.Listen("tcp",  ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pe, err := prometheus.NewExporter(prometheus.Options{
    Namespace: "demo",
  })
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	} 	else { log.Printf("Prometheus exporter has started.") }
	// start prometheous expoter
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(":" + strconv.Itoa(exporterPort), mux); err != nil {
			log.Fatalf("Failed to run Prometheus /metrics endpoint: %v", err)
		}
	}()
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register ocgrpc server views: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}), )
	//s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}	else { log.Printf("gRPC listener has started.") }


}
