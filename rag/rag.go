package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
)

type KBEntry struct {
	ID        string
	Text      string
	Embedding []float32
}

var MemoryStore []*KBEntry

func AddDocument(ctx context.Context, id, text string, model string) error {
	embed, err := EmbedText(ctx, model, text)
	if err != nil {
		return err
	}

	MemoryStore = append(MemoryStore, &KBEntry{
		ID:        id,
		Text:      text,
		Embedding: embed,
	})
	return nil
}

type scoredDoc struct {
	Text  string
	Score float32
}

func cosineSim(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// Retrieve top K docs
func Query(ctx context.Context, queryText string, topK int, model string) ([]string, error) {
	qEmbed, err := EmbedText(ctx, model, queryText)
	if err != nil {
		return nil, err
	}

	scores := []scoredDoc{}
	for _, doc := range MemoryStore {
		score := cosineSim(qEmbed, doc.Embedding)
		scores = append(scores, scoredDoc{Text: doc.Text, Score: score})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	res := []string{}
	for i := 0; i < len(scores) && i < topK; i++ {
		res = append(res, scores[i].Text)
	}

	return res, nil
}

// ChatWithRAG: performs retrieval + LLM answer
func ChatWithRAG(ctx context.Context, question, embedModel, llmModel string) (string, []string, error) {
	//  Retrieve top 3 docs
	docs, err := Query(ctx, question, 3, embedModel)
	if err != nil {
		return "", nil, err
	}

	// Prepare prompt for LLM
	prompt := "Use the following context to answer the question.\n\n"
	for i, d := range docs {
		prompt += fmt.Sprintf("[%d] %s\n", i+1, d)
	}
	prompt += "\nQuestion: " + question + "\nAnswer:"

	// Call Ollama LLM
	cmd := exec.CommandContext(ctx, "ollama", "generate", llmModel, "--json")
	cmd.Stdin = stringReader(prompt)

	out, err := cmd.Output()
	if err != nil {
		return "", nil, fmt.Errorf("LLM call failed: %w", err)
	}

	var resp struct {
		Output string `json:"output"`
	}

	if err := json.Unmarshal(out, &resp); err != nil {
		return "", nil, fmt.Errorf("failed to parse LLM output: %w", err)
	}

	return resp.Output, docs, nil
}

// helper: convert string to reader
func stringReader(s string) *os.File {
	r, w, _ := os.Pipe()
	go func() {
		defer w.Close()
		w.Write([]byte(s))
	}()
	return r
}