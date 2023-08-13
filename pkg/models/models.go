package models

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Date string `json:"date"`
}

type Diary struct {
	Notes []Note `json:"notebooks"`
}

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
