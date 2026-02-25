package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	type Todo struct {
		ID        int    `json:"id"`
		Completed bool   `json:"completed"`
		Body      string `json:"body"`
	}

	todoList := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"msg": "hello world"})
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(http.StatusBadRequest).JSON(
				fiber.Map{"error": "Todo Body is required"},
			)
		}
		todo.ID = len(todoList) + 1
		todoList = append(todoList, *todo)
		return c.Status(http.StatusCreated).JSON(todo)
	})

	app.Patch("/update/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todoList {
			if strconv.Itoa(todo.ID) == id {
				todoList[i].Completed = true
				return c.Status(http.StatusOK).JSON(todoList[i])
			}
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	})

	app.Delete("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todoList {
			if strconv.Itoa(todo.ID) == id {
				todoList = append(todoList[:i], todoList...)
				return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true"})
			}
		}

		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
