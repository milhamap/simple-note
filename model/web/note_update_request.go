package web

type NoteUpdateRequest struct {
	Id      int    `validate:"required" json:"id"`
	Title   string `validate:"required,max=100|min=1" json:"title"`
	Content string `validate:"required,max=255|min=1" json:"content"`
}
