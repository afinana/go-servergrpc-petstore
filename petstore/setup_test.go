package petstore

import (
	"context"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	testApp        *Application
	client         *mongo.Client
	mongoAvailable bool
)

// skipIfNoMongo is a test helper that skips integration tests if MongoDB is unreachable.
func skipIfNoMongo(t *testing.T) {
	t.Helper()
	if !mongoAvailable {
		t.Skip("Skipping integration test: MongoDB is not available")
	}
}

func TestMain(m *testing.M) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	if mongoDatabase == "" {
		mongoDatabase = "petstore_test"
	}

	infoLog := log.New(io.Discard, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	co := options.Client().ApplyURI(mongoURI).SetServerSelectionTimeout(2 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, co)
	if err == nil {
		if err = client.Ping(ctx, readpref.Primary()); err == nil {
			mongoAvailable = true
			log.Println("MongoDB connected for integration tests")
		}
	}

	if mongoAvailable {
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
	} else {
		log.Println("MongoDB is not available. Integration tests will be skipped; unit tests will run.")
	}

	code := m.Run()

	if mongoAvailable && client != nil {
		discCtx, discCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer discCancel()
		_ = client.Disconnect(discCtx)
	}

	os.Exit(code)
}
