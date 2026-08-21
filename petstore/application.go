package petstore

import (
	"fmt"
	"log"
	"runtime/debug"
)

// Application implements SwaggerPetstoreServiceServer with injected repository interfaces.
type Application struct {
	UnimplementedSwaggerPetstoreServiceServer
	errorLog *log.Logger
	infoLog  *log.Logger
	pets     PetRepository
	stores   StoreRepository
	users    UserRepository
}

func (app *Application) serverError(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
}

// NewLog initializes a new Application with logging and repository dependencies.
func NewLog(inLog *log.Logger, errLog *log.Logger,
	pets PetRepository, stores StoreRepository, users UserRepository) *Application {

	// Initialize a new instance of application containing the dependencies.
	app := &Application{errorLog: errLog, infoLog: inLog, pets: pets, stores: stores, users: users}
	return app
}
