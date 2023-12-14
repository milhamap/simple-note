package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/milhamap/simple-note/controller"
	"github.com/milhamap/simple-note/exception"
)

func NewRouter(noteController controller.NoteControllerInterface) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/notes", noteController.FindAll)
	router.GET("/api/notes/:noteId", noteController.FindById)
	router.POST("/api/notes", noteController.Create)
	router.PUT("/api/notes/:noteId", noteController.Update)
	router.DELETE("/api/notes/:noteId", noteController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}
