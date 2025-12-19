package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Gateway struct {
	GoServiceUrl           string
	PythonServiceUrl       string
	GoServiceEndpoints     map[string]string
	PythonServiceEndpoints map[string]string
}

func NewGateway() *Gateway {
	return &Gateway{
		GoServiceUrl:     "http://localhost:8081",
		PythonServiceUrl: "http://localhost:8082",
		GoServiceEndpoints: map[string]string{
			"ImageToPdf": "/convert/image-pdf",
			"MergePdfs":  "/merge",
		},
		PythonServiceEndpoints: map[string]string{
			"PdfToImage": "/pdf-to-jpg",
		},
	}
}

func (g *Gateway) Proxy(serviceType, endpointType string) gin.HandlerFunc {
	if serviceType == "" || endpointType == "" {
		return func(c *gin.Context) {
			c.JSON(500, gin.H{"error": fmt.Errorf("no service type"), "message": "Internal Server Error"})
		}
	} else {
		return func(c *gin.Context) {

			baseUrl,targetPath, err := g.GetServiceEndpoint(serviceType, endpointType)
			if err != nil {
				c.JSON(500, gin.H{"error": fmt.Errorf("no service type"), "message": "Internal Server Error"})
			}

			targetUrl, err := url.Parse(baseUrl)

			if err != nil {
				c.JSON(500, gin.H{"error": fmt.Errorf("url parsing error"), "message": "Internal Server Error"})
			}

			proxy := httputil.NewSingleHostReverseProxy(targetUrl)
			c.Request.URL.Path=targetPath

			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}

}

func (g *Gateway) GetServiceEndpoint(serviceType, endpointType string) (string,string, error) {

	var baseUrl string
	var targetPath string

	if serviceType == "go" {
		baseUrl= g.GoServiceUrl
		targetPath= g.GoServiceEndpoints[endpointType]
		

	} else if serviceType == "python" {
		baseUrl= g.PythonServiceUrl
		targetPath= g.PythonServiceEndpoints[endpointType]
	}

	return baseUrl,targetPath, nil

}
