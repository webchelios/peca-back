package model

type Entry struct {
	ID      int    `json: "id"`
	Title   string `json: "title"`
	Content string `json: "content"`
	Image   string `json: "image"`
}
