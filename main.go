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

// convexRequest and convexResponse are used to marshal and unmarshal 
// the request and response bodies when communicating with the Convex API. 
// The Args field in convexRequest is of type any to allow for flexibility 
// in the arguments passed to different endpoints. The Value and ErrorData 
// fields in convexResponse are of type json.RawMessage to allow for deferred 
// parsing of the response data, which can be useful when
// the structure of the response is not known in advance.
type convexRequest struct {
	Path string      `json:"path"`
	Args any `json:"args"`
}

// convexResponse represents the structure of the response returned by the Convex API.
type convexResponse struct {
	Status    string          `json:"status"`
	Value     json.RawMessage `json:"value"`
	ErrorData json.RawMessage `json:"errorData"`
}

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

	// Define the root route that returns a simple JSON response to confirm that the server is running.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"msg": "hello world"})
	})

	app.Get("/todos", func(c *fiber.Ctx) error {
		result, err := callConvex(convexURL, convexAdminKey, "query", "todoActions:list", fiber.Map{})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(result)
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

		result, err := callConvex(convexURL, convexAdminKey, "mutation", "todoActions:add", fiber.Map{"body": todo.Body})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusCreated).Send(result)
	})

	app.Patch("/update/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := callConvex(convexURL, convexAdminKey, "mutation", "todoActions:setCompleted", fiber.Map{"id": id})
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(result)
	})
	
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
