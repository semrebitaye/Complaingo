package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

// send prompt to OpenAI and receive response
func CallOpenAI(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "nil", fmt.Errorf("missing OPENAI_API_KEY")
	}

	reqBody := RequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	// if api error occur
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var response ResponseBody
	if err := json.Unmarshal(body, &response); err != nil {
		return "", nil
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return response.Choices[0].Message.Content, nil
}
