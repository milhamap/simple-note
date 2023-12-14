package helper

import (
	"github.com/milhamap/simple-note/model/domain"
	"github.com/milhamap/simple-note/model/web"
)

func ToNoteResponse(note domain.Note) web.NoteResponse {
	return web.NoteResponse{
		Id:      note.Id,
		Title:   note.Title,
		Content: note.Content,
	}
}

func ToNoteResponses(notes []domain.Note) []web.NoteResponse {
	var noteResponses []web.NoteResponse
	for _, note := range notes {
		noteResponses = append(noteResponses, ToNoteResponse(note))
	}
	return noteResponses
}
