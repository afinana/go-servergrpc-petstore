/*
 * Swagger Petstore
 *
 * This is a sample server Petstore server.  You can find out more about Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For this sample, you can use the api key `special-key` to test the authorization filters.
 *
 * API version: 1.0.6
 * Contact: apiteam@swagger.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	api "middleland.net/swaggerapi/petstore"
)

func main() {

	// Define command-line flags
	serverAddr := flag.String("serverAddr", "localhost", "HTTP server network address")
	serverPort := flag.Int("serverPort", 8090, "HTTP server network port")
	redisURI := flag.String("redisURI", "redis://:@localhost:6379/", "Database hostname url")

	flag.Parse()

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Establish database connection
	opt, err := redis.ParseURL(*redisURI)
	if err != nil {
		errLog.Fatal(err)
	}
	client := redis.NewClient(opt)

	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		errLog.Fatal(err)
	}
	infoLog.Printf("Database connection established")

	app := api.NewLog(
		infoLog,
		errLog,
		&api.PetModel{
			C: client,
		},
		&api.StoreModel{
			C: client,
		},
		&api.UserModel{
			C: client,
		},
	)

	// Initialize a new http.Server struct.
	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	infoLog.Printf("Starting server on %s", serverURI)

	// Start listening in serverPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new Grpc Server
	server := grpc.NewServer()
	api.RegisterSwaggerPetstoreServiceServer(server, app)

	// Start serve request
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
