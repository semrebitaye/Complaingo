package handler

import (
	"crud_api/internal/middleware"
	"crud_api/internal/openai"
	"encoding/json"
	"net/http"
)

type AIRequest struct {
	Prompt string `json:"prompt"`
}

type AIResponse struct {
	Answer string `json:"answer"`
}

func AIChatHandler(w http.ResponseWriter, r *http.Request) {
	var req AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, err)
		return
	}

	answer, err := openai.CallOpenAI(req.Prompt)
	if err != nil {
		middleware.WriteError(w, err)
	}

	resp := AIResponse{Answer: answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
