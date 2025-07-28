package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
)

// Summarize builds a prompt and queries the LLM
func Summarize(question string, notes []string) (string, error) {
	InitOpenAIClient()
	if initErr != nil {
		return "", initErr
	}

	contextText := buildContextFromNotes(notes)
	fullPrompt := fmt.Sprintf(`You are an AI assistant built into a personal note-taking app. You only respond using the user's notes, documents, and other saved content. Your job is to help the user recall, summarize, and answer questions based only on their own notes.

If a question cannot be answered from the user's data, say: "I couldn't find anything in your notes related to that." Do not use external knowledge or make up information.

Be concise, clear, and helpful. Use bullet points or formatting where useful. If the question is vague, ask follow-up questions to clarify what the user wants.

Never hallucinate facts. If a topic is mentioned in the notes but incomplete, clearly say so.

User's Notes:
%s

Question:
%s

Answer:`, contextText, question)

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: "gpt-4.1-nano",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fullPrompt),
		},
	})
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// buildContextFromNotes creates a string representation of user notes
func buildContextFromNotes(notes []string) string {
	var sb strings.Builder
	for i, note := range notes {
		sb.WriteString(fmt.Sprintf("Note %d:\n%s\n\n", i+1, note))
	}
	return sb.String()
}
