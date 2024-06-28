package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//nolint:testableexamples
func Example() {
	type ShortenRequest struct {
		URL string `json:"url"`
	}

	// Создаем JSON-данные для запроса
	shortenRequest := ShortenRequest{URL: "http://ya.ru"}
	jsonBody, err := json.Marshal(shortenRequest)
	if err != nil {
		fmt.Println("Ошибка при маршализации JSON:", err)

		return
	}

	// Создаем POST-запрос
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"http://localhost:8080/api/shorten",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)

		return
	}

	// Устанавливаем заголовок Content-Type для JSON
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)

		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)

		return
	}

	fmt.Println("Ответ от сервера:", string(body))
}
