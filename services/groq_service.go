package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"resume-analyzer-api/models"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// GroqService handles AI analysis using Groq API (OpenAI-compatible)
type GroqService struct {
	client *openai.Client
}

// NewGroqService creates a new Groq service instance
func NewGroqService(apiKey string) *GroqService {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"

	client := openai.NewClientWithConfig(config)

	return &GroqService{
		client: client,
	}
}

// AnalyzeCV analyzes a CV against a job description using Groq
func (s *GroqService) AnalyzeCV(cvText, jobTitle, jobDescription string) (*models.AnalysisResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Build the prompt
	prompt := buildAnalysisPrompt(cvText, jobTitle, jobDescription)

	// Create chat completion request
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   4096,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Choices) == 0 {
		return nil, errors.New("empty response from Groq")
	}

	responseText := resp.Choices[0].Message.Content

	// Clean up the response (remove markdown code blocks if present)
	responseText = cleanJSONResponse(responseText)

	// Parse the JSON response
	var analysis models.AnalysisResponse
	if err := json.Unmarshal([]byte(responseText), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &analysis, nil
}

// buildAnalysisPrompt creates the structured prompt for CV analysis
func buildAnalysisPrompt(cvText, jobTitle, jobDescription string) string {
	return fmt.Sprintf(`You are an expert resume analyzer and career coach. Analyze the following CV/resume against the job requirements and provide a detailed analysis.

IMPORTANT: You must respond with ONLY valid JSON. Do not include any markdown, code blocks, or explanatory text. Just the raw JSON object.

CV/Resume Content:
---
%s
---

Job Title: %s

Job Description/Requirements:
---
%s
---

Analyze this CV against the job requirements and provide a comprehensive analysis in the following JSON format. Return ONLY the JSON, no other text:

{
  "match_score": <number 0-100>,
  "match_summary": "<brief summary of how well the CV matches the position>",
  "skill_analysis": {
    "matched_skills": ["<skill1>", "<skill2>"],
    "missing_skills": ["<skill1>", "<skill2>"],
    "recommendations": ["<specific recommendation1>", "<recommendation2>"]
  },
  "structure_feedback": {
    "overall_score": <number 0-100>,
    "sections": [
      {
        "name": "<section name>",
        "score": <number 0-100>,
        "feedback": "<specific feedback>",
        "suggestions": ["<suggestion1>", "<suggestion2>"]
      }
    ]
  },
  "improvement_suggestions": [
    {
      "priority": "<high|medium|low>",
      "section": "<section name>",
      "suggestion": "<specific actionable suggestion>"
    }
  ],
  "ats_tips": [
    "<tip1 for passing ATS systems>",
    "<tip2>"
  ]
}

Guidelines:
- match_score: Calculate based on skills match, experience relevance, and overall fit
- Identify at least 4-6 matched skills and 3-5 missing skills relevant to the job
- Provide 3-4 specific skill recommendations
- Analyze at least 4 CV sections (Work Experience, Skills, Education, Summary/Objective, etc.)
- Give 3-5 improvement suggestions with varying priorities
- Provide 3-4 ATS optimization tips specific to this CV and job

Return ONLY the JSON object, nothing else.`, cvText, jobTitle, jobDescription)
}

// cleanJSONResponse removes markdown code blocks and cleans up the response
func cleanJSONResponse(response string) string {
	// Remove markdown code blocks if present
	response = strings.TrimSpace(response)

	// Remove ```json prefix
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
	}

	// Remove ``` suffix
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}

	// Trim whitespace
	response = strings.TrimSpace(response)

	return response
}
