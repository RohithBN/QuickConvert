package main

import (
	imagepdf "github.com/RohithBn/img-x-converter"
	"github.com/RohithBn/router"
)


func main(){


	imageToXConvService:= imagepdf.NewService()
	imageToXConvServiceHandler:= imagepdf.NewHandler(imageToXConvService)


    r:= router.SetupRoutes(&router.App{ImgHandler: imageToXConvServiceHandler})

	r.Run(":8080")

}