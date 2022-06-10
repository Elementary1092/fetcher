package domain

import (
	"strconv"
	"strings"
)

type Post struct {
	UserId int32  `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (p *Post) String() string {
	var str strings.Builder
	str.WriteString("Post {")
	str.WriteString("; user_id: " + strconv.Itoa(int(p.UserId)))
	str.WriteString("; title: " + shorten(p.Title))
	str.WriteString("; body: " + shorten(p.Body) + "}")

	return str.String()
}

func shorten(s string) string {
	if len(s) > 20 {
		return s[0:20] + "..."
	}

	return s
}
