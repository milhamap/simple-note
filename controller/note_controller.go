package controller

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"simple-note/helper"
	"simple-note/model/web"
	"simple-note/service"
	"strconv"
)

type NoteController struct {
	NoteService service.NoteServiceInterface
}

// belajar wire bind
func NewNoteController(serviceInterface service.NoteServiceInterface) *NoteController {
	return &NoteController{
		NoteService: serviceInterface,
	}
}

func (controller *NoteController) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	noteCreateRequest := web.NoteCreateRequest{}
	helper.ReadFromRequestBody(request, &noteCreateRequest)

	noteResponse := controller.NoteService.Create(request.Context(), noteCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   noteResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NoteController) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	noteUpdateRequest := web.NoteUpdateRequest{}
	helper.ReadFromRequestBody(request, &noteUpdateRequest)

	noteId := params.ByName("noteId")
	id, err := strconv.Atoi(noteId)
	helper.PanicIfError(err)
	noteUpdateRequest.Id = id

	noteResponse := controller.NoteService.Update(request.Context(), noteUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   noteResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NoteController) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	noteId := params.ByName("noteId")
	id, err := strconv.Atoi(noteId)
	helper.PanicIfError(err)

	controller.NoteService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NoteController) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	noteId := params.ByName("noteId")
	id, err := strconv.Atoi(noteId)
	helper.PanicIfError(err)

	noteResponse := controller.NoteService.FindById(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   noteResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *NoteController) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	noteResponses := controller.NoteService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   noteResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
