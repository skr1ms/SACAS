package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"codegrader-backend/internal/config"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService interface {
	AnalyzeCode(code, fileType string) (int, string, error)
	CheckForPlagiarism(code, fileType string, existingSubmissions []string) (bool, string, error)
}

type openAIService struct {
	client *openai.Client
}

func NewOpenAIService(cfg *config.Config) OpenAIService {
	if cfg.OpenAI.APIKey == "" {
		log.Printf("WARNING: OpenAI API key is not set")
	}
	client := openai.NewClient(cfg.OpenAI.APIKey)
	return &openAIService{client: client}
}

func (s *openAIService) AnalyzeCode(code, fileType string) (int, string, error) {
	language := getLanguageName(fileType)

	log.Printf("Starting OpenAI analysis for %s code", language)

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
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   500,
			Temperature: 0.7,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return 0, "", fmt.Errorf("failed to analyze code with OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		log.Printf("OpenAI returned no choices")
		return 0, "", fmt.Errorf("no response from OpenAI")
	}

	response := resp.Choices[0].Message.Content
	log.Printf("OpenAI analysis completed successfully")

	grade, feedback := parseGPTResponse(response)

	return grade, feedback, nil
}

func (s *openAIService) CheckForPlagiarism(code, fileType string, existingSubmissions []string) (bool, string, error) {
	if len(existingSubmissions) == 0 {
		return false, "", nil
	}

	language := getLanguageName(fileType)
	log.Printf("Starting plagiarism check for %s code against %d existing submissions", language, len(existingSubmissions))

	existingCode := strings.Join(existingSubmissions, "\n\n--- NEXT SUBMISSION ---\n\n")

	prompt := fmt.Sprintf(`Проанализируй следующий код на языке %s на предмет плагиата.

Новый код для проверки:
%s

Существующие решения для сравнения:
%s

Определи, является ли новый код копией или очень похожим на одно из существующих решений.
Учитывай:
1. Идентичность или почти идентичность структуры кода
2. Одинаковые названия переменных и функций
3. Идентичную логику решения
4. Минимальные изменения (переименование переменных, изменение комментариев)

Если код является плагиатом или очень похож на существующее решение, ответь:
ПЛАГИАТ: Да
Объяснение: [объяснение почему это плагиат]

Если код оригинальный, ответь:
ПЛАГИАТ: Нет
Объяснение: Код имеет оригинальную структуру и подход к решению`, language, code, existingCode)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   800,
			Temperature: 0.3,
		},
	)

	if err != nil {
		log.Printf("OpenAI plagiarism check error: %v", err)
		return false, "", fmt.Errorf("failed to check plagiarism with OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		log.Printf("OpenAI returned no choices for plagiarism check")
		return false, "", fmt.Errorf("no response from OpenAI for plagiarism check")
	}

	response := resp.Choices[0].Message.Content
	log.Printf("Plagiarism check completed successfully")

	isPlagiarism, explanation := parsePlagiarismResponse(response)
	return isPlagiarism, explanation, nil
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

func parsePlagiarismResponse(response string) (bool, string) {
	lines := strings.Split(response, "\n")
	isPlagiarism := false
	explanation := response

	for _, line := range lines {
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "плагиат:") {
			if strings.Contains(lowerLine, "да") {
				isPlagiarism = true
			}
			break
		}
	}

	return isPlagiarism, explanation
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
