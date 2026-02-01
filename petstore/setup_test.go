package petstore

import (
	"log"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

var (
	testApp *Application
	rdb     *redis.Client
	mr      *miniredis.Miniredis
)

func TestMain(m *testing.M) {
	// Setup miniredis
	var err error
	mr, err = miniredis.Run()
	if err != nil {
		log.Fatalf("Could not start miniredis: %s", err)
	}

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create redis client configuration to connect to miniredis
	rdb = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	infoLog.Printf("Miniredis started at %s", mr.Addr())

	testApp = NewLog(
		infoLog,
		errLog,
		&PetModel{
			Rdb: rdb,
		},
		&StoreModel{
			Rdb: rdb,
		},
		&UserModel{
			Rdb: rdb,
		},
	)

	// Run tests
	code := m.Run()

	// Teardown
	mr.Close()

	os.Exit(code)
}
