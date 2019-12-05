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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"time"
	"strconv"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"io/ioutil"
	"path/filepath"
  "gopkg.in/yaml.v2"
)


const (
	//address     = "localhost:50051"
	//defaultName = "world"
	//iterations  = 100
	//numThraeds = 10
	configFile    = "./config/config.yaml"
	maxRetries      = 10
)


func getConfig(configFile string) map[string]interface{}{
	log.Printf("Reading...%v\n", configFile)
	var config map[string]interface{}
	filename, _ := filepath.Abs(configFile)
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

func sendMessage(thread_name string, address string, iterations int, defaultName string, finished chan bool) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	name := thread_name
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	// Contact the server and print out its response.
	var counter, retry int = 0, 0;
	for {
		counter++
  	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	  defer cancel()
	  r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	  if err != nil {
			log.Printf("could not greet: %v", err)
			retry++
		} else {
			log.Printf("Greeting: %s", r.GetMessage())
			retry=0
	  }
		if iterations > 0 && iterations < counter  { break }
		if retry > maxRetries {
			log.Fatalf("Too many retries, quit...")
		}
	}
  finished <- true
}

func main() {
	log.Printf("Starting...")
	config := getConfig(configFile)
	address := config["client"].(map[interface{}]interface{})["address"].(string)
	iterations := config["client"].(map[interface{}]interface{})["iterations"].(int)
	defaultName := config["client"].(map[interface{}]interface{})["defaultName"].(string)
	numThraeds := config["client"].(map[interface{}]interface{})["numThraeds"].(int)
	finished := make(chan bool)
  for i := 0 ; i < numThraeds; i++ {
    go sendMessage("Thread" + strconv.Itoa(i), address, iterations, defaultName, finished)
	}
	<- finished
	log.Printf("Done...")
}