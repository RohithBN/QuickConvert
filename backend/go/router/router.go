package router

import (
	imagepdf "github.com/RohithBn/img-x-converter"
	"github.com/RohithBn/lib"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	ImgHandler *imagepdf.Handler
}

func SetupRoutes(app *App) *gin.Engine {

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/convert/image-pdf", lib.GinErrorWrapper(app.ImgHandler.ConvertToPDFHandler))
	r.POST("/convert/png-jpeg",lib.GinErrorWrapper(app.ImgHandler.ConvertPNGToJPEGHandler))
	r.POST("/convert/jpeg-png",lib.GinErrorWrapper(app.ImgHandler.ConvertJPEGToPNGHandler))
	return r

}
