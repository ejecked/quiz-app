package handlers

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

type Answer struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
	Image    *string  `json:"image"`
}

func SubmitAnswer(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(500).SendString("Database connection not initialized in handler.")
	}

	id := c.Params("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid answer ID.")
	}

	var correct bool
	err = db.QueryRow("SELECT is_correct FROM answers WHERE id=$1", aid).Scan(&correct)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).SendString("Ответ не найден.")
		}
		log.Printf("Database error querying answer ID %d: %v", aid, err)
		return c.Status(500).SendString("Internal server error.")
	}

	if correct {
		return c.SendString("Молодец! Это правильный ответ.")
	}
	return c.SendString("Это неправильный ответ.")
}
