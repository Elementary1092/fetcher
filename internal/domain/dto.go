package domain

type PostDTO struct {
	Id     int32  `json:"id"`
	UserId int32  `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (dto *PostDTO) FormPost() *Post {
	return &Post{
		UserId: dto.UserId,
		Title:  dto.Title,
		Body:   dto.Title,
	}
}
