package main

import (
	"github.com/RohithBN/proxy"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	gateway := proxy.NewGateway()

	r.POST("/convert/image-pdf", gateway.Proxy("go", "ImageToPdf"))

	r.Run(":8080")

}
