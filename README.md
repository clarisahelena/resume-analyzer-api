# Resume Analyzer API

A Go-based REST API that analyzes resumes/CVs against job descriptions using Groq AI (llama-3.3-70b-versatile). Extracts text from PDF resumes and provides detailed analysis including match scores, skill gap analysis, and improvement suggestions.

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Gin (HTTP routing)
- **AI**: Groq API (llama-3.3-70b-versatile model via OpenAI-compatible SDK)
- **PDF Parsing**: ledongthc/pdf library

## Features

- **Match Score**: Instant percentage score showing CV-to-job match
- **Skill Analysis**: Identifies matched and missing skills with recommendations
- **Structure Feedback**: Detailed feedback on CV sections
- **Improvement Suggestions**: Prioritized actionable suggestions
- **ATS Tips**: Optimization tips for Applicant Tracking Systems

## Project Structure

```
resume-analyzer-api/
├── main.go              # Entry point & server setup
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── Dockerfile           # Docker build configuration
├── .env.example         # Environment variables template
├── handlers/
│   ├── analyze.go       # CV analysis endpoint handler
│   └── health.go        # Health check endpoint
├── services/
│   ├── groq_service.go  # Groq AI integration
│   └── pdf_service.go   # PDF text extraction
└── models/
    └ analysis.go        # Request/response data structures
```

## Prerequisites

- Go 1.24 or later
- Groq API key ([Get one here](https://console.groq.com/keys))
- Docker (optional, for containerized deployment)

## Environment Variables

Create a `.env` file from the example:

```bash
cp .env.example .env
```

Required variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `GROQ_API_KEY` | Your Groq API key (required) | - |
| `PORT` | Server port | `8080` |
| `MAX_FILE_SIZE` | Max PDF upload size in bytes | `5242880` (5MB) |
| `CORS_ORIGIN` | Frontend URL for CORS (production) | `localhost:3000` only |

## Local Development

### 1. Install dependencies

```bash
go mod download
```

### 2. Set environment variables

Edit `.env` and add your Groq API key:

```env
GROQ_API_KEY=your_groq_api_key_here
```

### 3. Run the server

```bash
go run main.go
```

The API will start on `http://localhost:8080`.

## Docker Deployment

### Build the image

```bash
docker build -t resume-analyzer-api .
```

### Run the container

```bash
docker run -p 8080:8080 -e GROQ_API_KEY=your_key resume-analyzer-api
```

The API will be available at `http://localhost:8080`.

## API Endpoints

### GET /api/health

Health check endpoint.

**Response**:
```json
{
  "status": "healthy",
  "message": "Resume Analyzer API is running"
}
```

### POST /api/analyze

Analyzes a CV against a job description.

**Request** (multipart/form-data):
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `cv_file` | File | Yes | PDF file (max 5MB) |
| `job_title` | String | Yes | Job title |
| `job_description` | String | Yes | Full job description |

**Response**:
```json
{
  "match_score": 85,
  "match_summary": "Your CV is a strong match...",
  "skill_analysis": {
    "matched_skills": ["Go", "PostgreSQL", "REST API"],
    "missing_skills": ["Kubernetes", "gRPC"],
    "recommendations": ["Add Redis experience..."]
  },
  "structure_feedback": {
    "overall_score": 78,
    "sections": [
      {
        "name": "Work Experience",
        "score": 90,
        "feedback": "Well structured...",
        "suggestions": []
      }
    ]
  },
  "improvement_suggestions": [
    {
      "priority": "high",
      "section": "Skills",
      "suggestion": "Add missing technologies..."
    }
  ],
  "ats_tips": [
    "Use keywords from the job description...",
    "Avoid tables and columns..."
  ]
}
```

**Error Responses**:

| Status | Error Code | Description |
|--------|------------|-------------|
| 400 | `missing_field` | Required field missing |
| 400 | `missing_file` | No CV file uploaded |
| 400 | `invalid_file` | Invalid PDF format or size exceeds limit |
| 422 | `extraction_failed` | PDF text extraction failed |
| 500 | `analysis_failed` | AI analysis failed |
| 504 | `timeout` | Analysis request timed out |

## CORS Configuration

The API allows requests from `http://localhost:3000` by default for local development.

For production, set the `CORS_ORIGIN` environment variable to your frontend URL:

```env
CORS_ORIGIN=https://your-frontend.vercel.app
```

Multiple origins are supported — localhost is always included for development.

## File Requirements

- **Format**: PDF only
- **Size**: Maximum 5MB
- **Content**: Text-based PDFs work best (scanned documents may not parse correctly)

## Troubleshooting

### Server won't start

- Ensure `GROQ_API_KEY` is set in `.env`
- Verify Go 1.24+ is installed: `go version`
- Check port 8080 is not already in use

### PDF parsing fails

- Ensure the PDF is text-based, not a scanned image
- Check file size is under 5MB
- Verify the file is a valid PDF (not corrupted)

### Groq API errors

- Verify your API key is valid and has quota
- Check [Groq API status](https://status.groq.com/) for issues
- Ensure you're using the correct model (`llama-3.3-70b-versatile`)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the LICENSE file for details.

## Acknowledgments

- [Groq](https://groq.com/) for fast AI inference
- [ledongthc/pdf](https://github.com/ledongthc/pdf) for PDF parsing
- [Gin](https://gin-gonic.com/) web framework