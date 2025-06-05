package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"quiz-app/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

type Answer struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
	Image    *string  `json:"image_path"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "host=db user=postgres password=postgres dbname=quiz port=5432 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	handlers.SetDB(db)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/", "./public")

	app.Get("/questions", getQuestions)
	app.Post("/questions", createQuestion)
	app.Put("/questions/:id", updateQuestion)
	app.Delete("/questions/:id", deleteQuestion)
	app.Post("/upload/:id", uploadImage)
	app.Post("/answer/:id", handlers.SubmitAnswer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server is starting on port :%s", port)
	log.Fatal(app.Listen(":" + port))
}

func getQuestions(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, question, image_path FROM questions")
	if err != nil {
		log.Printf("DB Error getting questions: %v", err)
		return c.Status(500).SendString("Internal server error.")
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		var image sql.NullString
		if err := rows.Scan(&q.ID, &q.Question, &image); err != nil {
			log.Printf("DB Error scanning question: %v", err)
			return c.Status(500).SendString("Internal server error.")
		}
		if image.Valid {
			q.Image = &image.String
		} else {
			q.Image = nil
		}

		answerRows, err := db.Query("SELECT id, text, is_correct FROM answers WHERE question_id=$1", q.ID)
		if err != nil {
			log.Printf("DB Error getting answers for question %d: %v", q.ID, err)
			return c.Status(500).SendString("Internal server error.")
		}
		for answerRows.Next() {
			var a Answer
			if err := answerRows.Scan(&a.ID, &a.Text, &a.IsCorrect); err != nil {
				log.Printf("DB Error scanning answer for question %d: %v", q.ID, err)
				return c.Status(500).SendString("Internal server error.")
			}
			q.Answers = append(q.Answers, a)
		}
		answerRows.Close()

		questions = append(questions, q)
	}
	return c.JSON(questions)
}

func createQuestion(c *fiber.Ctx) error {
	var input struct {
		Question string   `json:"question"`
		Answers  []Answer `json:"answers"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Printf("BodyParser error (createQuestion): %v", err)
		return c.Status(400).SendString(err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("DB Error starting transaction (createQuestion): %v", err)
		return c.Status(500).SendString("Internal server error.")
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow("INSERT INTO questions (question) VALUES ($1) RETURNING id", input.Question).Scan(&id)
	if err != nil {
		log.Printf("DB Error inserting question: %v", err)
		return c.Status(500).SendString(err.Error())
	}

	for _, a := range input.Answers {
		_, err := tx.Exec("INSERT INTO answers (question_id, text, is_correct) VALUES ($1, $2, $3)",
			id, a.Text, a.IsCorrect)
		if err != nil {
			log.Printf("DB Error inserting answer for qID %d: %v", id, err)
			return c.Status(500).SendString(err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("DB Error committing transaction (createQuestion): %v", err)
		return c.Status(500).SendString("Internal server error.")
	}

	return c.JSON(fiber.Map{"id": id})
}

func updateQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	var input struct {
		Question string   `json:"question"`
		Answers  []Answer `json:"answers"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Printf("BodyParser error (updateQuestion): %v", err)
		return c.Status(400).SendString(err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("DB Error starting transaction (updateQuestion): %v", err)
		return c.Status(500).SendString("Internal server error.")
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE questions SET question=$1 WHERE id=$2", input.Question, id)
	if err != nil {
		log.Printf("DB Error updating question %s: %v", id, err)
		return c.Status(500).SendString(err.Error())
	}

	_, err = tx.Exec("DELETE FROM answers WHERE question_id=$1", id)
	if err != nil {
		log.Printf("DB Error deleting old answers for qID %s: %v", id, err)
	}

	for _, a := range input.Answers {
		_, err := tx.Exec("INSERT INTO answers (question_id, text, is_correct) VALUES ($1, $2, $3)",
			id, a.Text, a.IsCorrect)
		if err != nil {
			log.Printf("DB Error inserting new answer for qID %s: %v", id, err)
			return c.Status(500).SendString(err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("DB Error committing transaction (updateQuestion): %v", err)
		return c.Status(500).SendString("Internal server error.")
	}

	return c.SendStatus(fiber.StatusOK)
}

func deleteQuestion(c *fiber.Ctx) error {
	id := c.Params("id")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("DB Error starting transaction (deleteQuestion): %v", err)
		return c.Status(500).SendString("Internal server error.")
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM questions WHERE id=$1", id)
	if err != nil {
		log.Printf("DB Error deleting question %s: %v", id, err)
		return c.Status(500).SendString(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("DB Error committing transaction (deleteQuestion): %v", err)
		return c.Status(500).SendString("Internal server error.")
	}

	return c.SendStatus(fiber.StatusOK)
}

func uploadImage(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("FormFile error (uploadImage): %v", err)
		return c.Status(400).SendString(err.Error())
	}

	uploadDir := "./public/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			log.Printf("Failed to create upload directory %s: %v", uploadDir, err)
			return c.Status(500).SendString("Failed to create upload directory.")
		}
	}

	savePath := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		log.Printf("Failed to save file %s: %v", savePath, err)
		return c.Status(500).SendString("Failed to save file.")
	}

	relativePath := "/uploads/" + file.Filename
	_, err = db.Exec("UPDATE questions SET image_path = $1 WHERE id = $2", relativePath, id)
	if err != nil {
		log.Printf("DB Error updating image_path for qID %s: %v", id, err)
		os.Remove(savePath)
		return c.Status(500).SendString("Failed to update image path in database.")
	}
	return c.SendStatus(fiber.StatusOK)
}
