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

import (2
	"context"
	"log"
	"os"
	"time"
	"strconv"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	//"google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	//address     = "localhost:50051"
	address     = "localhost:50051"
	defaultName = "world"
	iterations  = 100
	num_thraeds = 10
)

/*
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
*/
func sendMessage(thread_name string, iterations int, finished chan bool) {
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
	for i:=0; i < iterations ;i++{
  	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	  defer cancel()
	  r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	  if err != nil {
			//log.Fatalf("could not greet: %v", err)
			log.Printf("could not greet: %v", err)
	  }
	  log.Printf("Greeting: %s", r.GetMessage())
  }
  finished <- true
}

func main() {
	log.Printf("Starting...")
	finished := make(chan bool)
  for i := 0; i < num_thraeds; i++{
    go sendMessage("Thread" + strconv.Itoa(i), iterations, finished)
	}
	<- finished
	log.Printf("Done...")
}