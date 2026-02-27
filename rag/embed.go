package rag

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type ollamaEmbeddingRequest struct {
// 	Model  string `json:"model"`
// 	Input  string `json:"input"`
// 	// Ollama may accept different shapes, adapt if necessary
// }

// type ollamaEmbeddingResponse struct {
// 	Embedding []float32 `json:"embedding"`
// }

// // EmbedText calls local Ollama embeddings endpoint.
// // Assumes Ollama daemon at http://localhost:11434
// func EmbedText(ctx context.Context, model, text string) ([]float32, error) {
// 	reqBody := ollamaEmbeddingRequest{
// 		Model: model,
// 		Input: text,
// 	}
// 	b, _ := json.Marshal(reqBody)

// 	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/embeddings", bytes.NewBuffer(b))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("ollama embeddings request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode >= 300 {
// 		return nil, fmt.Errorf("embeddings API returned status %d", resp.StatusCode)
// 	}

// 	var out struct {
// 		Embedding []float32 `json:"embedding"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
// 		return nil, fmt.Errorf("failed decode embedding: %w", err)
// 	}
// 	return out.Embedding, nil
// }
