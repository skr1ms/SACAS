package services

import (
	"context"
	"fmt"
	"strings"

	"codegrader-backend/internal/config"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService interface {
	AnalyzeCode(code, fileType string) (int, string, error)
}

type openAIService struct {
	client *openai.Client
}

func NewOpenAIService(cfg *config.Config) OpenAIService {
	client := openai.NewClient(cfg.OpenAI.APIKey)
	return &openAIService{client: client}
}

func (s *openAIService) AnalyzeCode(code, fileType string) (int, string, error) {
	language := getLanguageName(fileType)

	prompt := fmt.Sprintf(`Проанализируй следующий код на языке %s и оцени его по критериям:
1. Читаемость и структура
2. Соблюдение стиль-гайдов
3. Логика решения
4. Эффективность алгоритма

Выставь оценку от 3 до 5 баллов и дай краткие комментарии с рекомендациями.

Код:
%s

Ответ должен быть в формате:
Оценка: [число от 3 до 5]
Комментарии: [твои комментарии]`, language, code)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 500,
		},
	)

	if err != nil {
		return 0, "", err
	}

	response := resp.Choices[0].Message.Content
	grade, feedback := parseGPTResponse(response)

	return grade, feedback, nil
}

func parseGPTResponse(response string) (int, string) {
	lines := strings.Split(response, "\n")
	grade := 3
	feedback := response

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "оценка:") {
			if strings.Contains(line, "5") {
				grade = 5
			} else if strings.Contains(line, "4") {
				grade = 4
			} else if strings.Contains(line, "3") {
				grade = 3
			}
		}
	}

	return grade, feedback
}

func getLanguageName(fileType string) string {
	switch fileType {
	case ".cpp":
		return "C++"
	case ".java":
		return "Java"
	case ".js":
		return "JavaScript"
	case ".kt":
		return "Kotlin"
	case ".py":
		return "Python"
	default:
		return "Unknown"
	}
}
