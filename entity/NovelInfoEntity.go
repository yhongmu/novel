package entity

type NovelInfoEntity struct {
	ID string			`json:"id"`
	Title string		`json:"title"`
	Author string		`json:"author"`
	Category string		`json:"category"`
	Status string		`json:"status"`
	Description string	`json:"description"`
	Update string		`json:"update"`
	BookImg string		`json:"book_img,omitempty"`
	Source string		`json:"source,omitempty"`
	SourceURL string	`json:"source_url,omitempty"`
}
