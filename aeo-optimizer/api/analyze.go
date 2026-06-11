package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Request body structure
type RequestPayload struct {
	Text    string `json:"text"`
	Keyword string `json:"keyword"`
}

// Response structure
type AnalysisResult struct {
	WordCount     int      `json:"wordCount"`
	AeoScore      int      `json:"aeoScore"`
	Entities      []string `json:"entities"`
	SchemaSnippet string   `json:"schemaSnippet"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST request", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Lỗi đọc dữ liệu", http.StatusBadRequest)
		return
	}

	text := strings.ToLower(payload.Text)
	wordCount := len(strings.Fields(text))

	// Logic cơ bản: Tính điểm AEO dựa trên độ dài và việc chứa keyword
	aeoScore := 50
	if wordCount > 100 {
		aeoScore += 20
	}
	if payload.Keyword != "" && strings.Contains(text, strings.ToLower(payload.Keyword)) {
		aeoScore += 30
	}

	// Trích xuất Entities (Giả lập việc tách từ khóa quan trọng)
	entities := []string{payload.Keyword, "AI Optimization", "Generative Search"}

	// Tự động tạo FAQ Schema Markup (Rất quan trọng cho AEO)
	schema := `{
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [{
    "@type": "Question",
    "name": "Khái niệm chính về ` + payload.Keyword + ` là gì?",
    "acceptedAnswer": {
      "@type": "Answer",
      "text": "Đoạn văn bản cung cấp thông tin chi tiết về ` + payload.Keyword + `, được tối ưu hóa cho các công cụ trả lời tự động."
    }
  }]
}`

	result := AnalysisResult{
		WordCount:     wordCount,
		AeoScore:      aeoScore,
		Entities:      entities,
		SchemaSnippet: schema,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}