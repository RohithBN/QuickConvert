package lib

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type handlerFunc func(c *gin.Context) error

func GinErrorWrapper(handlerFunc handlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := handlerFunc(c); err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "message": "Internal Server Error"})
		}
	}
}

func ParseStringToInt(str string) (int, error) {
	var result int
	_, err := fmt.Sscanf(str, "%d", &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
