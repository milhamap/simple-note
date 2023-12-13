package service

import (
	"context"
	"simple-note/model/web"
)

type NoteServiceInterface interface {
	Create(ctx context.Context, request web.NoteCreateRequest) web.NoteResponse
	Update(ctx context.Context, request web.NoteUpdateRequest) web.NoteResponse
	Delete(ctx context.Context, categoryId int)
	FindById(ctx context.Context, categoryId int) web.NoteResponse
	FindAll(ctx context.Context) []web.NoteResponse
}
