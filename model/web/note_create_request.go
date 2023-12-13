package web

type NoteCreateRequest struct {
	Title   string `validate:"required,max=100|min=1" json:"title"`
	Content string `validate:"required,max=255|min=1" json:"content"`
}
