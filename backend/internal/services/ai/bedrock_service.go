package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/research-paper-analyzer/backend/internal/config"
	"github.com/research-paper-analyzer/backend/internal/models"
)

// BedrockAIService uses AWS Bedrock with Claude for AI operations.
type BedrockAIService struct {
	client  *bedrockruntime.Client
	modelID string
}

// claudeRequest represents the request body for Claude on Bedrock.
type claudeRequest struct {
	AnthropicVersion string           `json:"anthropic_version"`
	MaxTokens        int              `json:"max_tokens"`
	Messages         []claudeMessage  `json:"messages"`
	System           string           `json:"system,omitempty"`
	Temperature      float64          `json:"temperature,omitempty"`
}

// claudeMessage represents a single message in the Claude conversation.
type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// claudeResponse represents the response from Claude on Bedrock.
type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

// NewBedrockAIService creates a new Bedrock AI service with the configured AWS credentials.
func NewBedrockAIService(cfg *config.Config) (*BedrockAIService, error) {
	var awsOpts []func(*awsconfig.LoadOptions) error

	awsOpts = append(awsOpts, awsconfig.WithRegion(cfg.AWSRegion))

	// Use explicit credentials if provided
	if cfg.AWSAccessKey != "" && cfg.AWSSecretKey != "" {
		awsOpts = append(awsOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AWSAccessKey, cfg.AWSSecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(awsCfg)

	return &BedrockAIService{
		client:  client,
		modelID: cfg.BedrockModel,
	}, nil
}

// GenerateSummary sends the paper text to Claude via Bedrock for comprehensive analysis.
func (b *BedrockAIService) GenerateSummary(text string) (*AnalysisResult, error) {
	// Truncate text if it's too long for the model context window
	if len(text) > 80000 {
		text = text[:80000]
	}

	systemPrompt := `You are an expert academic research paper analyzer. Analyze the given paper text and provide a structured analysis. 
Respond ONLY with valid JSON in this exact format (no markdown, no code blocks, just raw JSON):
{
  "summary": "A comprehensive 3-5 sentence summary of the paper",
  "key_findings": "Bullet points starting with • for each key finding, separated by newlines",
  "methodology": "Description of the research methodology used",
  "limitations": "Bullet points starting with • for each limitation, separated by newlines",
  "future_scope": "Bullet points starting with • for each future research direction, separated by newlines",
  "keywords": "comma-separated list of 5-8 relevant keywords"
}`

	userMessage := fmt.Sprintf("Please analyze the following research paper:\n\n%s", text)

	response, err := b.invokeModel(systemPrompt, userMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Parse the JSON response
	var result struct {
		Summary     string `json:"summary"`
		KeyFindings string `json:"key_findings"`
		Methodology string `json:"methodology"`
		Limitations string `json:"limitations"`
		FutureScope string `json:"future_scope"`
		Keywords    string `json:"keywords"`
	}

	// Try to extract JSON from the response (handle potential markdown wrapping)
	jsonStr := extractJSON(response)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		// If JSON parsing fails, return the raw response as summary
		return &AnalysisResult{
			Summary:     response,
			KeyFindings: "• Analysis could not be fully structured. Please review the summary.",
			Methodology: "See summary for details.",
			Limitations: "• Automated analysis may not capture all nuances.",
			FutureScope: "• Further manual review recommended.",
			Keywords:    "research, analysis",
		}, nil
	}

	return &AnalysisResult{
		Summary:     result.Summary,
		KeyFindings: result.KeyFindings,
		Methodology: result.Methodology,
		Limitations: result.Limitations,
		FutureScope: result.FutureScope,
		Keywords:    result.Keywords,
	}, nil
}

// GenerateQuiz uses Claude to create multiple-choice questions about the paper.
func (b *BedrockAIService) GenerateQuiz(text string, numQuestions int) ([]models.QuizQuestionAI, error) {
	if len(text) > 60000 {
		text = text[:60000]
	}

	systemPrompt := fmt.Sprintf(`You are an expert quiz creator for academic papers. Create exactly %d multiple-choice questions about the given paper.
Respond ONLY with valid JSON array in this exact format (no markdown, no code blocks, just raw JSON):
[
  {
    "question": "The question text",
    "option_a": "First option",
    "option_b": "Second option",
    "option_c": "Third option",
    "option_d": "Fourth option",
    "correct_answer": "A",
    "explanation": "Why this answer is correct"
  }
]
The correct_answer must be exactly one of: "A", "B", "C", or "D".`, numQuestions)

	userMessage := fmt.Sprintf("Create quiz questions based on this paper:\n\n%s", text)

	response, err := b.invokeModel(systemPrompt, userMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate quiz: %w", err)
	}

	var questions []models.QuizQuestionAI
	jsonStr := extractJSON(response)
	if err := json.Unmarshal([]byte(jsonStr), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse quiz response: %w", err)
	}

	return questions, nil
}

// ChatWithContext answers a question using paper context and chat history.
func (b *BedrockAIService) ChatWithContext(question string, paperContext string, chatHistory []models.ChatMessage) (string, error) {
	systemPrompt := `You are a helpful research assistant. Answer the user's question based on the provided paper context. 
Be concise but thorough. If the context doesn't contain enough information to answer the question, say so honestly.
Do not make up information that isn't supported by the paper context.`

	// Build the user message with context
	userMessage := fmt.Sprintf("Paper context:\n%s\n\nQuestion: %s", paperContext, question)

	// Convert chat history to Claude format
	var history []claudeMessage
	for _, msg := range chatHistory {
		history = append(history, claudeMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	response, err := b.invokeModel(systemPrompt, userMessage, history)
	if err != nil {
		return "", fmt.Errorf("failed to chat: %w", err)
	}

	return response, nil
}

// invokeModel sends a request to Claude via Bedrock and returns the text response.
func (b *BedrockAIService) invokeModel(systemPrompt string, userMessage string, history []claudeMessage) (string, error) {
	// Build messages array with history + current message
	messages := make([]claudeMessage, 0)
	if history != nil {
		messages = append(messages, history...)
	}
	messages = append(messages, claudeMessage{
		Role:    "user",
		Content: userMessage,
	})

	req := claudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        4096,
		Messages:         messages,
		System:           systemPrompt,
		Temperature:      0.3,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := b.client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.modelID),
		ContentType: aws.String("application/json"),
		Body:        reqBody,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke Bedrock model: %w", err)
	}

	var response claudeResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		return "", fmt.Errorf("failed to parse Bedrock response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from Bedrock")
	}

	return response.Content[0].Text, nil
}

// extractJSON attempts to extract a JSON string from a response that might
// contain markdown code blocks or other wrapper text.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)

	// Check if the response is wrapped in markdown code blocks
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
	}

	return strings.TrimSpace(s)
}
