// Package pdf provides PDF text extraction and chunking services.
// It extracts raw text from uploaded PDF files and splits them into
// manageable chunks for AI processing and context retrieval.
package pdf

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	gopdf "github.com/ledongthuc/pdf"
)

// Service handles PDF text extraction and chunking.
type Service struct{}

// NewService creates a new PDF service instance.
func NewService() *Service {
	return &Service{}
}

// ExtractText reads a PDF file and extracts all text content.
// If the primary extraction method fails, it returns a fallback message.
// The filePath parameter should point to a locally accessible PDF file.
func (s *Service) ExtractText(filePath string) (string, error) {
	// Try primary extraction using ledongthuc/pdf
	text, err := s.extractWithLibrary(filePath)
	if err != nil {
		log.Printf("⚠️  Primary PDF extraction failed: %v. Using fallback.", err)
		return s.fallbackExtraction(filePath)
	}

	// Check if we actually extracted meaningful text
	cleanText := strings.TrimSpace(text)
	if len(cleanText) < 50 {
		log.Println("⚠️  Extracted text too short, using fallback")
		return s.fallbackExtraction(filePath)
	}

	return cleanText, nil
}

// extractWithLibrary performs text extraction using the ledongthuc/pdf library.
func (s *Service) extractWithLibrary(filePath string) (string, error) {
	// Open the PDF file
	f, r, err := gopdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	// Read text from all pages
	var buf bytes.Buffer
	totalPages := r.NumPage()

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			log.Printf("⚠️  Warning: failed to extract text from page %d: %v", pageNum, err)
			continue
		}

		buf.WriteString(text)
		buf.WriteString("\n\n") // Separate pages with double newline
	}

	return buf.String(), nil
}

// fallbackExtraction returns a message when PDF extraction fails.
// This could be enhanced to use OCR or other extraction methods.
func (s *Service) fallbackExtraction(filePath string) (string, error) {
	// Get file info for the response
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	fallbackText := fmt.Sprintf(
		"[PDF Text Extraction Notice]\n\n"+
			"The text content could not be fully extracted from this PDF file (%s, %.1f KB). "+
			"This may be because the PDF contains scanned images rather than selectable text, "+
			"or uses an encoding that is not supported by the current extraction library.\n\n"+
			"For best results, please ensure your PDF contains selectable text (not scanned images). "+
			"You can still use the chat feature to ask questions about this paper by providing context manually.",
		fileInfo.Name(),
		float64(fileInfo.Size())/1024.0,
	)

	return fallbackText, nil
}

// ExtractTextFromReader extracts text from a PDF provided as an io.Reader.
// It writes to a temporary file first since the PDF library requires file access.
func (s *Service) ExtractTextFromReader(reader io.Reader) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "paper-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write reader content to temp file
	if _, err := io.Copy(tmpFile, reader); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	// Close before reading to flush writes
	tmpFile.Close()

	return s.ExtractText(tmpFile.Name())
}

// SplitIntoChunks divides text into chunks of approximately chunkSize words
// with an overlap of overlapSize words between consecutive chunks.
// This overlap ensures context continuity at chunk boundaries.
func (s *Service) SplitIntoChunks(text string, chunkSize int, overlapSize int) []string {
	// Default values
	if chunkSize <= 0 {
		chunkSize = 500
	}
	if overlapSize <= 0 {
		overlapSize = 50
	}

	// Split text into words
	words := strings.Fields(text)

	// If the text is shorter than one chunk, return it as-is
	if len(words) <= chunkSize {
		return []string{strings.Join(words, " ")}
	}

	var chunks []string
	start := 0

	for start < len(words) {
		// Calculate end position for this chunk
		end := start + chunkSize
		if end > len(words) {
			end = len(words)
		}

		// Create the chunk
		chunk := strings.Join(words[start:end], " ")
		chunks = append(chunks, chunk)

		// Move start position forward, accounting for overlap
		start = end - overlapSize

		// Prevent infinite loop if overlap >= chunk size
		if start <= 0 || start >= len(words) {
			break
		}
		// If the remaining text would be too small, include it in the last chunk
		if len(words)-start < overlapSize {
			break
		}
	}

	return chunks
}
