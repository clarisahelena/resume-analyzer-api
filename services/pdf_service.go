package services

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PDFService handles PDF text extraction
type PDFService struct{}

// NewPDFService creates a new PDF service instance
func NewPDFService() *PDFService {
	return &PDFService{}
}

// ExtractText extracts text content from a PDF file
func (s *PDFService) ExtractText(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024)
	if header.Size > maxSize {
		return "", errors.New("file size exceeds 5MB limit")
	}

	// Read file content into buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", errors.New("failed to read file content")
	}

	// Validate PDF magic bytes
	pdfMagic := []byte{0x25, 0x50, 0x44, 0x46} // %PDF
	if !bytes.HasPrefix(buf.Bytes(), pdfMagic) {
		return "", errors.New("invalid file format: not a valid PDF")
	}

	// Create pdf reader from bytes
	reader, err := pdf.NewReader(bytes.NewReader(buf.Bytes()), header.Size)
	if err != nil {
		return "", errors.New("failed to parse PDF file")
	}

	// Get number of pages
	numPages := reader.NumPage()

	// Extract text from all pages
	var textBuilder strings.Builder
	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		// Extract text from page
		text, err := page.GetPlainText(nil)
		if err != nil {
			// Skip pages that can't be parsed
			continue
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n")
	}

	extractedText := strings.TrimSpace(textBuilder.String())
	if extractedText == "" {
		return "", errors.New("no text could be extracted from PDF")
	}

	return extractedText, nil
}

// ValidatePDF validates that a file is a valid PDF
func (s *PDFService) ValidatePDF(file multipart.File, header *multipart.FileHeader) error {
	// Check file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		return errors.New("file must be a PDF")
	}

	// Check file size
	maxSize := int64(5 * 1024 * 1024)
	if header.Size > maxSize {
		return errors.New("file size exceeds 5MB limit")
	}

	// Read first 4 bytes to check magic number
	buf := make([]byte, 4)
	if _, err := file.Read(buf); err != nil {
		return errors.New("failed to read file")
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return errors.New("failed to reset file position")
	}

	// Check PDF magic bytes (%PDF)
	pdfMagic := []byte{0x25, 0x50, 0x44, 0x46}
	if !bytes.Equal(buf, pdfMagic) {
		return errors.New("invalid file format: not a valid PDF")
	}

	return nil
}
