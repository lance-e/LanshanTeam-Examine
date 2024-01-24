package handle

import "github.com/gin-gonic/gin"

func HomePage(c *gin.Context) {
	var html = `<html><body>GitHub Login success </a></body></html>`
	c.Writer.Write([]byte(html))
}
