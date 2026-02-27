package jobs

// import (
// 	"encoding/json"
// 	"bitbucket/api_format/controllers"
// 	"bitbucket/api_format/db"
// 	"bitbucket/api_format/models"

// 	"github.com/robfig/cron/v3"
// 	"gorm.io/datatypes"
// )

// func Similarity(a, b string) float64 {
//     dist := levenshtein(a, b)
//     maxLen := float64(max(len(a), len(b)))
//     return 1.0 - float64(dist)/maxLen
// }

// func levenshtein(a, b string) int {
//     la := len(a)
//     lb := len(b)

//     d := make([][]int, la+1)
//     for i := range d {
//         d[i] = make([]int, lb+1)
//     }

//     for i := 0; i <= la; i++ {
//         d[i][0] = i
//     }
//     for j := 0; j <= lb; j++ {
//         d[0][j] = j
//     }

//     for i := 1; i <= la; i++ {
//         for j := 1; j <= lb; j++ {
//             cost := 0
//             if a[i-1] != b[j-1] {
//                 cost = 1
//             }

//             d[i][j] = min(
//                 d[i-1][j]+1,
//                 d[i][j-1]+1,
//                 d[i-1][j-1]+cost,
//             )
//         }
//     }

//     return d[la][lb]
// }

// func min(a, b, c int) int {
//     if a < b && a < c {
//         return a
//     }
//     if b < c {
//         return b
//     }
//     return c
// }

// func max(a, b int) int {
//     if a > b {
//         return a
//     }
//     return b
// }

// func decodeEmbedding(j datatypes.JSON) []float64 {
// 	var v []float64
// 	_ = json.Unmarshal(j, &v)
// 	return v
// }

// func AutoLearn(tenantInfoID int) {

// 	var logs []models.ChatLog

// 	// Fetch unlearned logs to learn from:
// 	// Only process rows that haven't been learned yet (learned=false)
// 	db.DB.
// 		Where(`
// 			tenant_info_id = ?
// 			AND learned = false
// 			AND confidence >= 0.75
// 			AND answer NOT LIKE 'I don''t have an answer yet%'
// 		`, tenantInfoID).
// 		Find(&logs)

// 	if len(logs) == 0 {
// 		return
// 	}

// 	var kbItems []models.KnowledgeBases
// 	db.DB.Where("tenant_info_id = ?", tenantInfoID).Find(&kbItems)

// 	for _, log := range logs {

// 		qEmbed, err := controllers.GetEmbedding(log.Question)
// 		if err != nil {
// 			continue
// 		}

// 		var bestKB *models.KnowledgeBases
// 		bestScore := 0.0

// 		for _, kb := range kbItems {

// 			kbEmbed := decodeEmbedding(kb.Embedding)
// 			if len(kbEmbed) == 0 {
// 				continue
// 			}

// 			score := controllers.CosineSimilarity(qEmbed, kbEmbed)
// 			if score > bestScore {
// 				bestScore = score
// 				bestKB = &kb
// 			}
// 		}

// 		// Learn from feedback-validated answers (learned=true, confidence=1.0)
// 		if log.Learned && log.Confidence == 1.0 && log.Feedback != "" {
// 			// User explicitly marked this as relevant - create/reinforce KB entry with feedback context
// 			newKB := models.KnowledgeBases{
// 				TenantInfoID: tenantInfoID,
// 				Question:     log.Question,
// 				Answer:       log.Answer, // Use user-validated answer
// 			}

// 			db.DB.Create(&newKB)
// 			db.DB.Model(&log).Update("learned", true)
// 			continue
// 		}

// 		// Learn from high-confidence auto-answers (confidence >= 0.75)
// 		// Only learn if very close to an existing KB
// 		if bestScore >= 0.85 && bestKB != nil {

// 			newKB := models.KnowledgeBases{
// 				TenantInfoID: tenantInfoID,
// 				Question:     log.Question,
// 				Answer:       bestKB.Answer, // reuse correct answer
// 			}

// 			db.DB.Create(&newKB)
// 		}

// 		// Mark processed regardless
// 		db.DB.Model(&log).Update("learned", true)
// 	}
// }

// func StartCronJobs() {
// 	c := cron.New()

// 	c.AddFunc("*/05 * * * *", func() {

// 		var tenants []models.TenantInfo
// 		db.DB.Find(&tenants)

// 		for _, t := range tenants {
// 			go AutoLearn(t.TenantInfoID)
// 		}
// 	})

// 	c.Start()
// }
