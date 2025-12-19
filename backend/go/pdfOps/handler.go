package pdfops

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) MergePDFsHandler(c *gin.Context) error {
	form, _ := c.MultipartForm()
	files, ok := form.File["files"]
	if !ok || len(files) == 0 {
		return fmt.Errorf("No files provided")
	}
    
	var pdfPaths []string

	// create a temporary directory to store uploaded PDFs
	if err := os.MkdirAll("./input/pdfs", 0755); err != nil {
		return err
	}

	for _, file := range files {
		pdfPath := "./input/pdfs/" + file.Filename
		if err := c.SaveUploadedFile(file, pdfPath); err != nil {
			return err
		}
		pdfPaths = append(pdfPaths, pdfPath)
	}

	mergedPdfPath, err := h.Service.MergePDFs(pdfPaths)
	if err != nil {
		return err
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=merged.pdf")
	c.File(mergedPdfPath)
	return nil
}
