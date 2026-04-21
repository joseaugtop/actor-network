package model

type Movie struct {
	Id    int64    `json:"id"`
	Title string   `json:"title"`
	Cast  []string `json:"cast"`
}
