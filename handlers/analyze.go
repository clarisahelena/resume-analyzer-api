package handlers

import (
	"context"
	"errors"
	"net/http"
	"resume-analyzer-api/models"
	"resume-analyzer-api/services"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler handles CV analysis requests
type AnalyzeHandler struct {
	groqService *services.GroqService
	pdfService  *services.PDFService
}

// NewAnalyzeHandler creates a new analyze handler instance
func NewAnalyzeHandler(groqService *services.GroqService, pdfService *services.PDFService) *AnalyzeHandler {
	return &AnalyzeHandler{
		groqService: groqService,
		pdfService:  pdfService,
	}
}

// Analyze handles POST /api/analyze requests
func (h *AnalyzeHandler) Analyze(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max memory
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse multipart form",
		})
		return
	}

	// Get form values
	jobTitle := c.PostForm("job_title")
	jobDescription := c.PostForm("job_description")

	// Validate required fields
	if jobTitle == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_field",
			Message: "job_title is required",
		})
		return
	}

	if jobDescription == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_field",
			Message: "job_description is required",
		})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("cv_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_file",
			Message: "CV file is required",
		})
		return
	}
	defer file.Close()

	// Validate PDF
	if err := h.pdfService.ValidatePDF(file, header); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_file",
			Message: err.Error(),
		})
		return
	}

	// Extract text from PDF
	cvText, err := h.pdfService.ExtractText(file, header)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Error:   "extraction_failed",
			Message: err.Error(),
		})
		return
	}

	// Analyze CV with Groq
	analysis, err := h.groqService.AnalyzeCV(cvText, jobTitle, jobDescription)
	if err != nil {
		// Check for specific error types
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, models.ErrorResponse{
				Error:   "timeout",
				Message: "Analysis request timed out",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "analysis_failed",
			Message: "Failed to analyze CV: " + err.Error(),
		})
		return
	}

	// Return successful analysis
	c.JSON(http.StatusOK, analysis)
}
