package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("hello worlds")
	app := fiber.New()

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

	log.Fatal(app.Listen(":8080"))
}
