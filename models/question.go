package models

type Answer struct {
	ID        int    `json:"id"`         // Поле для ID ответа
	Text      string `json:"text"`       // Текст ответа
	IsCorrect bool   `json:"is_correct"` // Флаг, правильный ли ответ
}

type Question struct {
	ID       int      `json:"id"`        // ID вопроса
	Question string   `json:"question"`  // Текст вопроса
	ImageURL string   `json:"image_url"` // Путь к изображению
	Answers  []Answer `json:"answers"`   // Массив ответов для вопроса
}
