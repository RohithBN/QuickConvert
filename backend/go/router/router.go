package router

import (
	imagepdf "github.com/RohithBn/img-x-converter"
	"github.com/gin-gonic/gin"
)

type App struct {
	ImgHandler *imagepdf.Handler
}

func SetupRoutes(app *App) *gin.Engine {

	r:= gin.Default()

	r.POST("/convert/image-pdf", func(c *gin.Context) {
		if err := app.ImgHandler.ConvertToPDFHandler(c); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})

	return r;


}