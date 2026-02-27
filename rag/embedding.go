package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

func EmbedText(ctx context.Context, model string, text string) ([]float32, error) {
	// Call Ollama CLI to get embeddings
	cmd := exec.CommandContext(ctx, "ollama", "embed", model, text, "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	var resp struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse embedding: %w", err)
	}

	return resp.Embedding, nil
}
