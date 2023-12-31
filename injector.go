//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-playground/validator"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"github.com/milhamap/simple-note/app"
	"github.com/milhamap/simple-note/controller"
	"github.com/milhamap/simple-note/middleware"
	"github.com/milhamap/simple-note/repository"
	"github.com/milhamap/simple-note/service"
	"net/http"
)

var noteSet = wire.NewSet(
	repository.NewNoteRepository,
	wire.Bind(new(repository.NoteRepositoryInterface), new(*repository.NoteRepository)),
	service.NewNoteService,
	wire.Bind(new(service.NoteServiceInterface), new(*service.NoteService)),
	controller.NewNoteController,
	wire.Bind(new(controller.NoteControllerInterface), new(*controller.NoteController)),
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewDB,
		validator.New,
		noteSet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		app.NewServer,
	)
	return nil
}
