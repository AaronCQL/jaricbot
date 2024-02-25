package model

import (
	"github.com/google/generative-ai-go/genai"
)

type Role string

type Content struct {
	Role Role   `json:"role"`
	Text string `json:"text"`
}

type History struct {
	Contents []*Content `json:"contents"`
}

const (
	RoleUser  Role = "user"
	RoleModel Role = "model"
)

func NewHistory() *History {
	return &History{
		Contents: []*Content{},
	}
}

func (h *History) Append(role Role, text string) *History {
	if len(h.Contents) == 0 {
		h.Contents = append(h.Contents, &Content{
			Role: role,
			Text: text,
		})
		return h
	}

	last := h.Contents[len(h.Contents)-1]
	if last.Role != role {
		h.Contents = append(h.Contents, &Content{
			Role: role,
			Text: text,
		})
		return h
	}

	last.Text += "\n\n" + text

	return h
}

// Returns a tuple of the user message and all historical messages that is ready
// to be consumed by Gemini's genai package.
func (h *History) ToGeminiContents() (genai.Text, []*genai.Content) {
	contents := []*genai.Content{}
	if len(h.Contents) == 0 {
		return "", contents
	}
	// We assume the last content is always the user content
	for _, c := range h.Contents[:len(h.Contents)-1] {
		contents = append(contents, &genai.Content{
			Role:  string(c.Role),
			Parts: []genai.Part{genai.Text(c.Text)},
		})
	}
	return genai.Text(h.Contents[len(h.Contents)-1].Text), contents
}
