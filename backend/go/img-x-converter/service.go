package imagepdf

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
)

type Service struct {
	TempPdfDir  string
	HTTPClient  *http.Client
	PDFConfig   PDFConfig
	TempJpegDir string
	TempPNGDir string
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
		TempPdfDir: "./temp",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		TempJpegDir: "./jpeg-images",
		PDFConfig: PDFConfig{
			PageSize:    "A4",
			Orientation: "P",
			Unit:        "mm",
			ImageMargin: 10,
			ImageWidth:  190,
		},
		TempPNGDir: "./png-images",
	}
}

func (s *Service) ConvertToPDF(imageURLs []string) (string, error) {

	if len(imageURLs) == 0 {
		return "", fmt.Errorf("no images provided")
	}

	pdf := fpdf.New(s.PDFConfig.Orientation, s.PDFConfig.Unit, s.PDFConfig.PageSize, "")

	if err := os.MkdirAll(s.TempPdfDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	uniqueKey := generateUniqueKey()

	for _, imageURL := range imageURLs {
		pdf.AddPage()
		pdf.ImageOptions(imageURL, s.PDFConfig.ImageMargin, s.PDFConfig.ImageMargin, s.PDFConfig.ImageWidth, 0, false, fpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	}

	err := pdf.OutputFileAndClose(s.TempPdfDir + "/" + uniqueKey + ".pdf")
	if err != nil {
		return "", err
	}

	// once PDF is created , delete the input images

	for _, imgUrl := range imageURLs {
		if err := os.Remove(imgUrl); err != nil {
			fmt.Printf("failed to delete image %s: %v\n", imgUrl, err)
		}
	}

	// return the path or URL to the generated PDF
	return uniqueKey + ".pdf", nil
}

func (s *Service) ConvertPNGToJPEG(imageUrl string) (string, error) {

	// first check if imageUrl is empty
	if len(imageUrl) == 0 {
		return "", fmt.Errorf("no image provided")
	}
	// check if jpeg dir exists , if not create it
	if err := os.MkdirAll("./jpeg-images", 0755); err != nil {
		return "", fmt.Errorf("failed to create jpeg directory: %w", err)
	}

	// generate unique key for jpeg image
	uniqueKey := generateUniqueKey()
	jpegImagePath := s.TempJpegDir + "/" + uniqueKey + ".jpeg"

	// download the png from imageUrl ie in the folder ./input
	img, err := os.Open(imageUrl)
	if err != nil {
		return "", fmt.Errorf("failed to open image %s: %w", imageUrl, err)
	}
	defer img.Close()

	// decode the png image
	pngImage, _, err := image.Decode(img)
	if err != nil {
		return "", fmt.Errorf("failed to decode png image %s: %w", imageUrl, err)
	}

	bounds := pngImage.Bounds()
	jpegImage := image.NewRGBA(bounds)

	// Fill entire image with white using
	white := &image.Uniform{color.White}

	//  fill the jpeg image with white background
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			jpegImage.Set(x, y, white)
		}
	}

	// now draw the png image over the white background
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			jpegImage.Set(x, y, pngImage.At(x, y))
		}
	}

	// Alternate method using draw package
	// draw.Draw(jpegImage,bounds,white,image.Point{},draw.Src)
	// Draw PNG over white background, compositing with alpha
	// draw.Draw(jpegImage,bounds,pngImage,bounds.Min,draw.Over)

	// save the jpeg file to jpegImagePath
	jpegFile, err := os.Create(jpegImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create jpeg file %s: %w", jpegImagePath, err)
	}
	defer jpegFile.Close()

	// encode the jpeg image
	err = jpeg.Encode(jpegFile, jpegImage, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encode jpeg image %s: %w", jpegImagePath, err)
	}

	return jpegImagePath, nil
}


func (s *Service) ConvertJPEGToPNG(imageUrl string) (string,error){
	// check if file path exists
	fmt.Println("Called j->png")
	if len(imageUrl) ==0 {
		return "", fmt.Errorf("no image provided")
	}

	// create a output dir if not exists
	if err:=os.MkdirAll(s.TempPNGDir,0755);err!=nil{
		return "", fmt.Errorf("failed to create png directory: %w", err)
	}

	//download the jpeg image from path
	img , err:= os.Open(imageUrl)
	if err!= nil {
		return "", fmt.Errorf("failed to open image %s: %w", imageUrl, err)
	}
	defer img.Close()

	// decode the jpeg image
	jpegImg,err:= jpeg.Decode(img)
	if err!= nil {
		return "", fmt.Errorf("failed to decode jpeg image %s: %w", imageUrl, err)
	}

	// generate unique key for png image
	uniqueKey:= generateUniqueKey()
	pngImagePath:= s.TempPNGDir + "/" + uniqueKey + ".png"

	// create the png file
	pngFile,err:= os.Create(pngImagePath)
	if err!= nil {
		return "", fmt.Errorf("failed to create png file %s: %w", pngImagePath, err)
	}
	defer pngFile.Close()

	// encode the png image
	err = png.Encode(pngFile,jpegImg)
	if err!= nil {
		return "", fmt.Errorf("failed to encode png image %s: %w", pngImagePath, err)
	}
	return pngImagePath, nil
}

func generateUniqueKey() string {
	return uuid.New().String()
}
