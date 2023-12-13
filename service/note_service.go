package service

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator"
	"simple-note/exception"
	"simple-note/helper"
	"simple-note/model/domain"
	"simple-note/model/web"
	"simple-note/repository"
)

type NoteService struct {
	NoteRepository repository.NoteRepositoryInterface
	DB             *sql.DB
	Validate       *validator.Validate
}

// belajar wire bind
func NewNoteService(noteRepository repository.NoteRepositoryInterface, DB *sql.DB, validate *validator.Validate) *NoteService {
	return &NoteService{
		NoteRepository: noteRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *NoteService) Create(ctx context.Context, request web.NoteCreateRequest) web.NoteResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	note := domain.Note{
		Title:   request.Title,
		Content: request.Content,
	}

	note = service.NoteRepository.Store(ctx, tx, note)

	return helper.ToNoteResponse(note)
}

func (service *NoteService) Update(ctx context.Context, request web.NoteUpdateRequest) web.NoteResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	note, err := service.NoteRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	note.Title = request.Title
	note.Content = request.Content

	note = service.NoteRepository.Update(ctx, tx, note)

	return helper.ToNoteResponse(note)
}

func (service *NoteService) Delete(ctx context.Context, noteId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	note, err := service.NoteRepository.FindById(ctx, tx, noteId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.NoteRepository.Delete(ctx, tx, note)
}

func (service *NoteService) FindById(ctx context.Context, noteId int) web.NoteResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	note, err := service.NoteRepository.FindById(ctx, tx, noteId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToNoteResponse(note)
}

func (service *NoteService) FindAll(ctx context.Context) []web.NoteResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	notes := service.NoteRepository.FindAll(ctx, tx)

	return helper.ToNoteResponses(notes)
}
