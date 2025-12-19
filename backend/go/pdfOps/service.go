package pdfops

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Service struct {
	TempPdfDir string
}

func NewService() *Service {
	return &Service{
		TempPdfDir: "./temp/pdfs",
	}
}

func (s *Service) MergePDFs(pdfPaths []string) (string, error) {

	if pdfPaths == nil || len(pdfPaths) == 0 {
		return "", nil
	}

	//create a temp folder if not exists and ensure concurrency safety
	if err := os.MkdirAll(s.TempPdfDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	uniqueKey := generateUniqueKey()
	outputPdfPath := s.TempPdfDir + "/" + uniqueKey + ".pdf"

	config := model.NewDefaultConfiguration()

	err := api.MergeCreateFile(pdfPaths, outputPdfPath, false, config)
	if err != nil {
		return "", fmt.Errorf("failed to merge PDFs: %w", err)
	}

	return outputPdfPath, nil

}

func generateUniqueKey() string {
	return uuid.New().String()
}
