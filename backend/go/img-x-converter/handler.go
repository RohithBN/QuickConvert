package imagepdf

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

func (h *Handler) ConvertToPDFHandler(c *gin.Context) error {
	form, _ := c.MultipartForm()
	files := form.File["files"]
	var imageUrls []string

	// create a temporary directory to store uploaded images

	if err := os.MkdirAll("./input", 0755); err != nil {
		return fmt.Errorf("failed to create input directory: %w", err)
	}

	for _, file := range files {
        imgPath := "./input/" + file.Filename
        if err := c.SaveUploadedFile(file, imgPath); err != nil {
            return fmt.Errorf("failed to save file %s: %w", file.Filename, err)
        }
        fmt.Println("Uploaded file:", imgPath)
        imageUrls = append(imageUrls, imgPath)
    }

	outputPdfPath, err := h.Service.ConvertToPDF(imageUrls)
	if err != nil {
		return err
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=converted.pdf")
	c.File("./temp/" + outputPdfPath)
	return nil

}
