package extractor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/ledongthuc/pdf"
	rscpdf "rsc.io/pdf"
)

// ExtractPages reads a PDF file and returns extracted plain text per page.
func ExtractPages(filePath string) ([]string, error) {
	if err := validatePDFSignature(filePath); err != nil {
		return nil, err
	}

	pages, primaryErr := extractWithLedong(filePath)
	if primaryErr == nil && len(pages) > 0 {
		return pages, nil
	}

	pages, fallbackErr := extractWithRsc(filePath)
	if fallbackErr == nil && len(pages) > 0 {
		return pages, nil
	}

	pages, commandErr := extractWithPdfToText(filePath)
	if commandErr == nil && len(pages) > 0 {
		return pages, nil
	}

	if isUnexpectedDelimiterError(primaryErr) && isUnexpectedDelimiterError(fallbackErr) && isPdfToTextUnavailable(commandErr) {
		return nil, fmt.Errorf("pdf has unsupported or malformed content streams; install pdftotext (Poppler) and set PDFTOTEXT_PATH to enable fallback extraction")
	}

	if primaryErr != nil && fallbackErr != nil && commandErr != nil {
		return nil, fmt.Errorf("extract text failed (ledongthuc: %v; rscpdf: %v; pdftotext: %v)", primaryErr, fallbackErr, commandErr)
	}
	if primaryErr != nil {
		return nil, primaryErr
	}
	if fallbackErr != nil {
		return nil, fallbackErr
	}
	if commandErr != nil {
		return nil, commandErr
	}

	return nil, fmt.Errorf("no text pages extracted from pdf")
}

func isUnexpectedDelimiterError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unexpected delimiter '>'")
}

func isPdfToTextUnavailable(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "pdftotext not available")
}

func validatePDFSignature(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open pdf: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 2048)
	readBytes, readErr := file.Read(buffer)
	if readErr != nil && readErr != io.EOF {
		return fmt.Errorf("read pdf header: %w", readErr)
	}
	if readBytes == 0 {
		return fmt.Errorf("input does not appear to be a valid pdf: empty file")
	}

	trimmed := bytes.TrimSpace(buffer[:readBytes])
	if !bytes.HasPrefix(trimmed, []byte("%PDF-")) {
		preview := strings.TrimSpace(string(trimmed))
		if len(preview) > 120 {
			preview = preview[:120] + "..."
		}
		if preview != "" {
			return fmt.Errorf("input does not appear to be a valid pdf (header preview: %q)", preview)
		}
		return fmt.Errorf("input does not appear to be a valid pdf")
	}

	return nil
}

func extractWithLedong(filePath string) ([]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	pages := make([]string, 0, r.NumPage())
	var firstPageErr error

	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			if firstPageErr == nil {
				firstPageErr = fmt.Errorf("extract page %d: %w", i, err)
			}
			continue
		}

		trimmed := strings.TrimSpace(content)
		if trimmed != "" {
			pages = append(pages, trimmed)
		}
	}

	if len(pages) > 0 {
		return pages, nil
	}

	if firstPageErr != nil {
		return nil, firstPageErr
	}

	return nil, fmt.Errorf("no text pages extracted from pdf")
}

func extractWithRsc(filePath string) ([]string, error) {
	reader, err := rscpdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open pdf with rscpdf: %w", err)
	}

	numPages := reader.NumPage()

	pages := make([]string, 0, numPages)
	var firstPageErr error

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		content, err := safeRscPageContent(page)
		if err != nil {
			if firstPageErr == nil {
				firstPageErr = fmt.Errorf("extract text for page %d with rscpdf: %w", i, err)
			}
			continue
		}
		if len(content.Text) == 0 {
			continue
		}

		sort.Slice(content.Text, func(a, b int) bool {
			dy := content.Text[a].Y - content.Text[b].Y
			if dy > 0.5 || dy < -0.5 {
				return content.Text[a].Y > content.Text[b].Y
			}
			return content.Text[a].X < content.Text[b].X
		})

		var builder strings.Builder
		lastY := content.Text[0].Y
		for idx, token := range content.Text {
			if idx > 0 {
				if abs(lastY-token.Y) > 2.0 {
					builder.WriteString("\n")
				} else {
					builder.WriteString(" ")
				}
			}
			builder.WriteString(strings.TrimSpace(token.S))
			lastY = token.Y
		}

		text := strings.TrimSpace(builder.String())
		if text == "" {
			if firstPageErr == nil {
				firstPageErr = fmt.Errorf("extract text for page %d with rscpdf: empty content", i)
			}
			continue
		}

		pages = append(pages, text)
	}

	if len(pages) > 0 {
		return pages, nil
	}

	if firstPageErr != nil {
		return nil, firstPageErr
	}

	return nil, fmt.Errorf("no text pages extracted from pdf with rscpdf")
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func safeRscPageContent(page rscpdf.Page) (content rscpdf.Content, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic while reading page content: %v", recovered)
		}
	}()

	content = page.Content()
	return content, nil
}

func extractWithPdfToText(filePath string) ([]string, error) {
	binaryPath, err := resolvePdfToTextBinary()
	if err != nil {
		return nil, fmt.Errorf("pdftotext not available")
	}

	cmd := exec.Command(binaryPath, "-layout", "-enc", "UTF-8", filePath, "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("run pdftotext: %w (%s)", err, abbreviateText(strings.TrimSpace(string(output)), 300))
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return nil, fmt.Errorf("pdftotext returned empty content")
	}

	parts := strings.Split(raw, "\f")
	pages := make([]string, 0, len(parts))
	for _, part := range parts {
		text := strings.TrimSpace(part)
		if text != "" {
			pages = append(pages, text)
		}
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no text pages extracted from pdf via pdftotext")
	}

	return pages, nil
}

// PdfToTextStatus reports whether pdftotext can be resolved at runtime.
func PdfToTextStatus() (bool, string, error) {
	binaryPath, err := resolvePdfToTextBinary()
	if err != nil {
		return false, "", err
	}
	return true, binaryPath, nil
}

func abbreviateText(value string, maxLen int) string {
	if maxLen <= 0 || len(value) <= maxLen {
		return value
	}
	return value[:maxLen] + "..."
}

var (
	pdfToTextResolveMu  sync.Mutex
	pdfToTextBinaryPath string
)

func resolvePdfToTextBinary() (string, error) {
	pdfToTextResolveMu.Lock()
	defer pdfToTextResolveMu.Unlock()

	if strings.TrimSpace(pdfToTextBinaryPath) != "" {
		if _, err := os.Stat(pdfToTextBinaryPath); err == nil {
			return pdfToTextBinaryPath, nil
		}
		pdfToTextBinaryPath = ""
	}

	envPath := strings.TrimSpace(os.Getenv("PDFTOTEXT_PATH"))
	if envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			pdfToTextBinaryPath = envPath
			return pdfToTextBinaryPath, nil
		}
	}

	if path, err := exec.LookPath("pdftotext"); err == nil {
		pdfToTextBinaryPath = path
		return pdfToTextBinaryPath, nil
	}

	candidates := buildPdfToTextCandidates()
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, statErr := os.Stat(candidate); statErr == nil {
			pdfToTextBinaryPath = candidate
			return pdfToTextBinaryPath, nil
		}
	}

	return "", fmt.Errorf("pdftotext binary not found")
}

func buildPdfToTextCandidates() []string {
	paths := make([]string, 0, 16)

	if runtime.GOOS == "linux" {
		paths = append(paths,
			"/usr/bin/pdftotext",
			"/usr/local/bin/pdftotext",
			"/bin/pdftotext",
		)
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		wingetPattern := filepath.Join(localAppData, "Microsoft", "WinGet", "Packages", "oschwartz10612.Poppler*", "poppler-*", "Library", "bin", "pdftotext.exe")
		if matches, err := filepath.Glob(wingetPattern); err == nil {
			paths = append(paths, matches...)
		}
	}

	programFiles := os.Getenv("ProgramFiles")
	if programFiles != "" {
		paths = append(paths,
			filepath.Join(programFiles, "poppler", "Library", "bin", "pdftotext.exe"),
			filepath.Join(programFiles, "Poppler", "Library", "bin", "pdftotext.exe"),
		)
	}

	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	if programFilesX86 != "" {
		paths = append(paths,
			filepath.Join(programFilesX86, "poppler", "Library", "bin", "pdftotext.exe"),
			filepath.Join(programFilesX86, "Poppler", "Library", "bin", "pdftotext.exe"),
		)
	}

	return paths
}
