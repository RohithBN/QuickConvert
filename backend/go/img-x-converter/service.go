package imagepdf

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
)

type Service struct {
	TempDir    string
	HTTPClient *http.Client
	PDFConfig  PDFConfig
}

type PDFConfig struct {
	PageSize    string
	Orientation string
	Unit        string
	ImageMargin float64
	ImageWidth  float64
}

func NewService() *Service {
	return &Service{
		TempDir: "./temp",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		PDFConfig: PDFConfig{
			PageSize:    "A4",
			Orientation: "P",
			Unit:        "mm",
			ImageMargin: 10,
			ImageWidth:  190,
		},
	}
}

func (s *Service) ConvertToPDF(imageURLs []string) (string, error) {

	if len(imageURLs) == 0 {
		return "", fmt.Errorf("no images provided")
	}

	pdf := fpdf.New(s.PDFConfig.Orientation, s.PDFConfig.Unit, s.PDFConfig.PageSize, "")

	if err:= os.MkdirAll(s.TempDir,0755); err!= nil{
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	uniqueKey := generateUniqueKey()

	for _, imageURL := range imageURLs {
		pdf.AddPage()
		pdf.ImageOptions(imageURL, s.PDFConfig.ImageMargin, s.PDFConfig.ImageMargin, s.PDFConfig.ImageWidth, 0, false, fpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	}

	err := pdf.OutputFileAndClose(s.TempDir + "/" + uniqueKey + ".pdf")
	if err != nil {
		return "", err
	}

	// once PDF is created , delete the input images

	for _, imgUrl := range imageURLs{
		if err:= os.Remove(imgUrl); err!= nil{
			fmt.Printf("failed to delete image %s: %v\n", imgUrl, err)
		}
	}


	// return the path or URL to the generated PDF
	return uniqueKey + ".pdf", nil
}

func generateUniqueKey() string {
	return uuid.New().String()
}
