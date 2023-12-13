package repository

import (
	"context"
	"database/sql"
	"errors"
	"simple-note/helper"
	"simple-note/model/domain"
)

type NoteRepository struct {
}

// belajar wire bind
func NewNoteRepository() *NoteRepository {
	return &NoteRepository{}
}

func (repository *NoteRepository) Store(ctx context.Context, tx *sql.Tx, note domain.Note) domain.Note {
	SQL := "INSERT INTO note (title, content) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, SQL, note.Title, note.Content)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	note.Id = int(id)
	return note
}

func (repository *NoteRepository) Update(ctx context.Context, tx *sql.Tx, note domain.Note) domain.Note {
	SQL := "UPDATE note SET title = ?, content = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, note.Title, note.Content, note.Id)
	helper.PanicIfError(err)

	return note
}

func (repository *NoteRepository) Delete(ctx context.Context, tx *sql.Tx, note domain.Note) {
	SQL := "DELETE FROM note WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, note.Id)
	helper.PanicIfError(err)
}

func (repository *NoteRepository) FindById(ctx context.Context, tx *sql.Tx, noteId int) (domain.Note, error) {
	SQL := "SELECT id, title, content FROM note WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, noteId)
	helper.PanicIfError(err)
	defer rows.Close()

	note := domain.Note{}
	if rows.Next() {
		err := rows.Scan(&note.Id, &note.Title, &note.Content)
		helper.PanicIfError(err)
		return note, nil
	} else {
		return note, errors.New("note not found!")
	}
}

func (repository *NoteRepository) FindAll(ctx context.Context, tx *sql.Tx) []domain.Note {
	SQL := "SELECT id, title, content FROM note"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		note := domain.Note{}
		rows.Scan(&note.Id, &note.Title, &note.Content)
		notes = append(notes, note)
	}
	return notes
}
