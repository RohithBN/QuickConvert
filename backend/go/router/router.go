package router

import (
	imagepdf "github.com/RohithBn/img-x-converter"
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

	r.POST("/convert/image-pdf", func(c *gin.Context) {
		if err := app.ImgHandler.ConvertToPDFHandler(c); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})

	return r

}
