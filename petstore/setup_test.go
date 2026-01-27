package petstore

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testApp *Application
	client  *mongo.Client
)

func TestMain(m *testing.M) {
	// Setup
	mongoURI := "mongodb://localhost:27017"
	mongoDatabase := "petstore_test"

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create mongo client configuration
	co := options.Client().ApplyURI(mongoURI)

	// Establish database connection
	var err error
	client, err = mongo.NewClient(co)
	if err != nil {
		log.Fatalf("Failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to mongo: %v", err)
	}

	infoLog.Printf("Database connection established")
	testApp = NewLog(
		infoLog,
		errLog,
		&PetModel{
			C: client.Database(mongoDatabase).Collection("pets"),
		},
		&StoreModel{
			C: client.Database(mongoDatabase).Collection("stores"),
		},
		&UserModel{
			C: client.Database(mongoDatabase).Collection("users"),
		},
	)

	// Run tests
	code := m.Run()

	// Teardown
	if err = client.Disconnect(context.Background()); err != nil {
		log.Printf("Failed to disconnect from mongo: %v", err)
	}

	os.Exit(code)
}
