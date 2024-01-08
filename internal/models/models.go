package models

type Project struct {
	Id          float64 `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Slug        string  `json:"slug"`
	CreatedAt   string  `json:"created_at"`
}

type HashTag struct {
	Id        float64 `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
}

type User struct {
	Id        float64 `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
}
