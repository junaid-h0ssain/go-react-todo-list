package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        string `json:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

// convexRequest is the payload shape expected by Convex HTTP APIs:
//
//	{
//	  path: "module:function",
//	  args: { ... }
//	}
type convexRequest struct {
	Path string `json:"path"`
	Args any    `json:"args"`
}

// convexResponse is the common envelope from Convex.
// For success, data is usually in Value.
// For failure, ErrorData contains the Convex-side error details.
type convexResponse struct {
	Status    string          `json:"status"`
	Value     json.RawMessage `json:"value"`
	ErrorData json.RawMessage `json:"errorData"`
}

// callConvex is the single HTTP bridge from Fiber -> Convex.
// endpoint: "query" or "mutation"
// path: Convex function path, e.g. "todoActions:list"
func callConvex(url string, adminKey string, endpoint string, path string, args any) (json.RawMessage, error) {
	body, err := json.Marshal(convexRequest{Path: path, Args: args})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/%s", url, endpoint), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if adminKey != "" {
		req.Header.Set("Authorization", "Convex "+adminKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("convex returned status %d: %s", resp.StatusCode, string(respBody))
	}

	parsed := convexResponse{}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return respBody, nil
	}

	if parsed.Status == "error" {
		return nil, fmt.Errorf("convex error: %s", string(parsed.ErrorData))
	}

	if len(parsed.Value) > 0 {
		return parsed.Value, nil
	}

	return respBody, nil
}

func main() {
	// Quick setup:
	// - PORT: Fiber server port
	// - CONVEX_URL: Convex deployment URL
	// - CONVEX_ADMIN_KEY: optional, needed for protected Convex functions
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	convexURL := os.Getenv("CONVEX_URL")
	convexAdminKey := os.Getenv("CONVEX_ADMIN_KEY")

	if convexURL == "" {
		log.Fatal("CONVEX_URL is required")
	}

	// Health check route.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"msg": "hello world"})
	})

	// GET /todos -> Convex query todoActions:list
	app.Get("/todos", func(c *fiber.Ctx) error {
		result, err := callConvex(convexURL, convexAdminKey, "query", "todoActions:list", fiber.Map{})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(result)
	})

	// POST /add -> Convex mutation todoActions:add
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

		result, err := callConvex(convexURL, convexAdminKey, "mutation", "todoActions:add", fiber.Map{"body": todo.Body})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusCreated).Send(result)
	})

	// PATCH /update/:id -> Convex mutation todoActions:setCompleted
	app.Patch("/update/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := callConvex(convexURL, convexAdminKey, "mutation", "todoActions:setCompleted", fiber.Map{"id": id})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(result)
	})

	// PATCH /update/body/:id -> Convex mutation todoActions:setBody
	app.Patch("/update/body/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(http.StatusBadRequest).JSON(
				fiber.Map{"error": "Todo Body is required"},
			)
		}

		result, err := callConvex(convexURL, convexAdminKey, "mutation", "todoActions:setBody", fiber.Map{"id": id, "body": todo.Body})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(result)
	})

	// DELETE /delete/:id -> Convex mutation todoActions:remove
	app.Delete("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := callConvex(convexURL, convexAdminKey, "mutation", "todoActions:remove", fiber.Map{"id": id})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"success": true})
	})

	log.Fatal(app.Listen(":" + PORT))
}
