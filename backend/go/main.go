package main

import (
	imagepdf "github.com/RohithBn/img-x-converter"
	pdfops "github.com/RohithBn/pdfOps"
	"github.com/RohithBn/router"
)


func main(){


	imageToXConvService:= imagepdf.NewService()
	imageToXConvServiceHandler:= imagepdf.NewHandler(imageToXConvService)

	pdfService:= pdfops.NewService()
	pdfServiceHandler:= pdfops.NewHandler(pdfService)


    r:= router.SetupRoutes(&router.App{ImgHandler: imageToXConvServiceHandler, PDFHandler: pdfServiceHandler})

	r.Run(":8081")

}