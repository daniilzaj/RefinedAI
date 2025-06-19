package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PromptRequest struct {
	Input string `json:"input"`
}

type DeepSeekRequest struct {
	Model    string        `json:"model"`
	Messages []MessageItem `json:"messages"`
}

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Ошибка при загрузке .env файла")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/generate", func(c echo.Context) error {
		req := new(PromptRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
		}

		prompt, err := callDeepSeek(req.Input)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to process prompt"})
		}

		return c.JSON(http.StatusOK, echo.Map{"prompt": prompt})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func callDeepSeek(input string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API ключ не найден")
	}

	apiURL := "https://openrouter.ai/api/v1/chat/completions"

	template := `Ты — профессиональный инженер промптов и эксперт по взаимодействию с языковыми моделями.

Цель: Преобразовать входной неструктурированный пользовательский промпт в улучшенный, ясный и высокоэффективный формат, пригодный для подачи в ИИ. Используй формат AUTOMAT или аналогичный, если уместно.

...
Теперь улучши этот промпт:
"%s"`

	finalPrompt := fmt.Sprintf(template, input)

	reqData := DeepSeekRequest{
		Model: "deepseek/deepseek-r1-0528-qwen3-8b:free",
		Messages: []MessageItem{
			{Role: "user", Content: finalPrompt},
		},
	}

	jsonBytes, _ := json.Marshal(reqData)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var parsed DeepSeekResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) > 0 {
		return parsed.Choices[0].Message.Content, nil
	}

	return "", nil
}
