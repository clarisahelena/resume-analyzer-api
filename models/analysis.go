package models

// AnalysisRequest represents the incoming analysis request
type AnalysisRequest struct {
	JobTitle       string `form:"job_title"`
	JobDescription string `form:"job_description"`
}

// AnalysisResponse represents the complete AI analysis response
type AnalysisResponse struct {
	MatchScore             int                     `json:"match_score"`
	MatchSummary           string                  `json:"match_summary"`
	SkillAnalysis          SkillAnalysis           `json:"skill_analysis"`
	StructureFeedback      StructureFeedback       `json:"structure_feedback"`
	ImprovementSuggestions []ImprovementSuggestion `json:"improvement_suggestions"`
	ATSTips                []string                `json:"ats_tips"`
}

// SkillAnalysis contains matched and missing skills analysis
type SkillAnalysis struct {
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	Recommendations []string `json:"recommendations"`
}

// StructureFeedback contains feedback on CV structure
type StructureFeedback struct {
	OverallScore int       `json:"overall_score"`
	Sections     []Section `json:"sections"`
}

// Section represents feedback for a single CV section
type Section struct {
	Name        string   `json:"name"`
	Score       int      `json:"score"`
	Feedback    string   `json:"feedback"`
	Suggestions []string `json:"suggestions"`
}

// ImprovementSuggestion represents a prioritized improvement suggestion
type ImprovementSuggestion struct {
	Priority   string `json:"priority"`
	Section    string `json:"section"`
	Suggestion string `json:"suggestion"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
