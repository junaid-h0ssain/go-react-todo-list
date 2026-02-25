package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// ConvexID is the Convex document ID type (a string).
type ConvexID = string

// Todo mirrors the shape returned by Convex for a todos document.
type Todo struct {
	ID        ConvexID `json:"_id"`
	Body      string   `json:"body"`
	Completed bool     `json:"completed"`
}

// convexClient holds the deployment URL and deploy key used to talk to Convex.
type convexClient struct {
	url       string // e.g. https://happy-animal-123.convex.cloud
	deployKey string
}

// convexRequest is the JSON body sent to /api/run/{functionIdentifier}.
// The function path goes in the URL, so only args and format are in the body.
type convexRequest struct {
	Args   map[string]any `json:"args"`
	Format string         `json:"format"`
}

// convexResponse is the envelope returned by all Convex HTTP API calls.
type convexResponse struct {
	Status       string          `json:"status"`
	Value        json.RawMessage `json:"value"`
	ErrorMessage string          `json:"errorMessage"`
}

// call sends a POST to /api/run/{functionIdentifier} and returns the raw JSON
// value on success, or an error on failure.
// path must use "file:export" format e.g. "todos:list" — it is converted to
// the URL segment "todos/list" automatically.
func (c *convexClient) call(path string, args map[string]any) (json.RawMessage, error) {
	body, err := json.Marshal(convexRequest{
		Args:   args,
		Format: "json",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// "todos:list" -> "/api/run/todos/list"
	urlPath := strings.ReplaceAll(path, ":", "/")
	endpoint := fmt.Sprintf("%s/api/run/%s", c.url, urlPath)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.deployKey != "" {
		req.Header.Set("Authorization", "Convex "+c.deployKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var cr convexResponse
	if err := json.Unmarshal(raw, &cr); err != nil {
		return nil, fmt.Errorf("unmarshal response (raw: %s): %w", string(raw), err)
	}
	if cr.Status != "success" {
		return nil, fmt.Errorf("convex error: %s", cr.ErrorMessage)
	}
	return cr.Value, nil
}

func (c *convexClient) query(path string, args map[string]any) (json.RawMessage, error) {
	return c.call(path, args)
}

func (c *convexClient) mutation(path string, args map[string]any) (json.RawMessage, error) {
	return c.call(path, args)
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	convexURL := os.Getenv("CONVEX_URL")
	convexDeployKey := os.Getenv("CONVEX_DEPLOY_KEY")

	if convexURL == "" {
		log.Fatal("CONVEX_URL is required in .env")
	}
	convexURL = strings.TrimRight(convexURL, "/")

	db := &convexClient{url: convexURL, deployKey: convexDeployKey}

	app := fiber.New()

	// GET / — health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"msg": "hello world"})
	})

	// GET /todos — list all todos
	app.Get("/todos", func(c *fiber.Ctx) error {
		raw, err := db.query("todos:list", map[string]any{})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var todos []Todo
		if err := json.Unmarshal(raw, &todos); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(todos)
	})

	// POST /add — create a new todo
	app.Post("/add", func(c *fiber.Ctx) error {
		var body struct {
			Body string `json:"body"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if body.Body == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Todo Body is required"})
		}

		raw, err := db.mutation("todos:add", map[string]any{"body": body.Body})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var todo Todo
		if err := json.Unmarshal(raw, &todo); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).JSON(todo)
	})

	// PATCH /update/:id — mark a todo as completed
	app.Patch("/update/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		raw, err := db.mutation("todos:setCompleted", map[string]any{"id": id})
		if err != nil {
			// Surface "Todo not found" as 404
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		var todo Todo
		if err := json.Unmarshal(raw, &todo); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(todo)
	})

	// DELETE /delete/:id — remove a todo
	app.Delete("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := db.mutation("todos:remove", map[string]any{"id": id})
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"success": "true"})
	})

	log.Fatal(app.Listen(":" + port))
}
