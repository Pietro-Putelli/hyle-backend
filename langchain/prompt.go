package langchain

import "strings"

/* Define interface for prompt methods */
type Prompt interface {
	Add(part string) Prompt
	GetPrompt() string
}

/* Define struct for prompt builder */
type PromptBuilder struct {
	parts []string
}

func (pb *PromptBuilder) Add(part string) Prompt {
	pb.parts = append(pb.parts, part)
	return pb
}

func (pb *PromptBuilder) GetPrompt() string {
	return strings.Join(pb.parts, " ")
}
