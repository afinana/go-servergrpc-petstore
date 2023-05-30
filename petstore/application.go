package petstore

import (
	"fmt"
	"log"
	"runtime/debug"
)

// grpcServer is used to implement SwaggerPetstoreServiceServer.
type Application struct {
	UnimplementedSwaggerPetstoreServiceServer
	errorLog *log.Logger
	infoLog  *log.Logger
	pets     *PetModel
	stores   *StoreModel
	users    *UserModel
}

func (app *Application) serverError(err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	fmt.Printf("%s/n", trace)
}

func NewLog(inLog *log.Logger, errLog *log.Logger,
	pets *PetModel, stores *StoreModel, users *UserModel) *Application {

	// Initialize a new instance of application containing the dependencies.
	app := &Application{errorLog: errLog, infoLog: inLog, pets: pets, stores: stores, users: users}
	return app

}
