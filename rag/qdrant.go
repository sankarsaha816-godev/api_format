// package rag

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	qdrant "github.com/qdrant/go-client/qdrant"
// 	"google.golang.org/protobuf/types/known/structpb"
// )

// var QClient *qdrant.Client

// // InitQdrant sets up an in-memory client (or adjust to remote)
// func InitQdrant() error {
// 	var err error
// 	// Create in-memory qdrant client
// 	QClient, err = qdrant.NewInMemoryClient()
// 	if err != nil {
// 		return fmt.Errorf("new in-memory client: %w", err)
// 	}
// 	log.Println("qdrant in-memory client ready")
// 	return nil
// }

// // EnsureCollection creates a collection if not exists; vectorSize is embedding size (model-dependent)
// func EnsureCollection(name string, vectorSize int) error {
// 	ctx := context.Background()
// 	// check existing collections
// 	_, err := QClient.GetCollections(ctx, &qdrant.GetCollectionsRequest{})
// 	if err != nil {
// 		// continue to create; library may return error depending on version
// 	}

// 	_, err = QClient.CreateCollection(ctx, &qdrant.CreateCollection{
// 		CollectionName: name,
// 		VectorsConfig: &qdrant.VectorsConfig{
// 			Config: &qdrant.VectorsConfig_Params{
// 				Params: &qdrant.VectorParams{
// 					Size:     int32(vectorSize),
// 					Distance: qdrant.Distance_Cosine,
// 				},
// 			},
// 		},
// 	})
// 	// ignore "already exists" errors if client lib returns such
// 	if err != nil {
// 		// some versions return error if exists; attempt to continue
// 		// you might check error string here
// 		log.Printf("create collection err (safe to ignore if exists): %v", err)
// 	}
// 	return nil
// }

// func InsertDoc(id string, text string, embedding []float32) error {
// 	ctx := context.Background()

// 	payload := map[string]*structpb.Value{
// 		"text": structpb.NewStringValue(text),
// 	}

// 	point := &qdrant.PointStruct{
// 		Id: &qdrant.PointId{
// 			PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
// 		},
// 		Vectors: &qdrant.Vectors{
// 			VectorsOptions: &qdrant.Vectors_Vector{
// 				Vector: &qdrant.Vector{
// 					Data: embedding,
// 				},
// 			},
// 		},
// 		Payload: payload,
// 	}

// 	_, err := QClient.UpsertPoints(ctx, &qdrant.UpsertPoints{
// 		CollectionName: "insight_kb",
// 		Points:         []*qdrant.PointStruct{point},
// 	})
// 	return err
// }

// func QueryDocs(query string, topK int, embedModel string) ([]string, error) {
// 	ctx := context.Background()
// 	embed, err := EmbedText(ctx, embedModel, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// wrap as qdrant.Vector
// 	vec := &qdrant.Vector{Data: embed}
// 	searchReq := &qdrant.SearchPoints{
// 		CollectionName: "insight_kb",
// 		Vector:         vec,
// 		Limit:          uint64(topK),
// 	}
// 	resp, err := QClient.Search(ctx, searchReq)
// 	if err != nil {
// 		return nil, fmt.Errorf("qdrant search: %w", err)
// 	}

// 	var results []string
// 	for _, r := range resp.Result {
// 		if v, ok := r.Payload["text"]; ok {
// 			results = append(results, v.GetStringValue())
// 		}
// 	}
// 	return results, nil
// }

package rag

// import (
// 	"context"
// 	"log"
// 	"strings"

// 	qdrant "github.com/qdrant/go-client/qdrant"
// 	"google.golang.org/protobuf/types/known/structpb"
// )

// var QClient *qdrant.Client

// // Init in-memory Qdrant
// func InitQdrant() error {
// 	var err error
// 	QClient, err = qdrant.NewInMemoryClient()
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("Qdrant in-memory client ready")
// 	return nil
// }

// // Create collection
// func EnsureCollection(name string, vectorSize int) error {
// 	ctx := context.Background()

// 	_, err := QClient.Collections.Create(ctx, &qdrant.CreateCollection{
// 		CollectionName: name,
// 		VectorsConfig: &qdrant.VectorsConfig{
// 			Config: &qdrant.VectorsConfig_Params{
// 				Params: &qdrant.VectorParams{
// 					Size:     uint64(vectorSize),
// 					Distance: qdrant.Distance_Cosine,
// 				},
// 			},
// 		},
// 	})

// 	// ignore "already exists"
// 	if err != nil && !isAlreadyExists(err) {
// 		return err
// 	}

// 	return nil
// }

// func isAlreadyExists(err error) bool {
// 	return err != nil && ( // qdrant usually returns these patterns
// 		strings.Contains(err.Error(), "exists") ||
// 			strings.Contains(err.Error(), "already"))
// }

// // Insert document + vector
// func InsertDoc(id string, text string, embedding []float32) error {
// 	ctx := context.Background()

// 	payload := map[string]*structpb.Value{
// 		"text": structpb.NewStringValue(text),
// 	}

// 	point := &qdrant.PointStruct{
// 		Id: &qdrant.PointId{
// 			PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
// 		},
// 		Vectors: &qdrant.Vectors{
// 			VectorsOptions: &qdrant.Vectors_Vector{
// 				Vector: embedding,
// 			},
// 		},
// 		Payload: payload,
// 	}

// 	_, err := QClient.Points.Upsert(ctx, &qdrant.UpsertPoints{
// 		CollectionName: "insight_kb",
// 		Points:         []*qdrant.PointStruct{point},
// 	})

// 	return err
// }

// // Query similar docs
// func QueryDocs(query string, topK int, embedModel string) ([]string, error) {
// 	ctx := context.Background()

// 	// Get embedding from Ollama or LLM
// 	embed, err := EmbedText(ctx, embedModel, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := QClient.Points.Search(ctx, &qdrant.SearchPoints{
// 		CollectionName: "insight_kb",
// 		Vector:         embed,
// 		Limit:          uint64(topK),
// 		WithPayload:    &qdrant.WithPayloadSelector{Enable: true},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var out []string
// 	for _, r := range resp.Result {
// 		val := r.Payload.AsMap()
// 		if txt, ok := val["text"].(string); ok {
// 			out = append(out, txt)
// 		}
// 	}

// 	return out, nil
// }
