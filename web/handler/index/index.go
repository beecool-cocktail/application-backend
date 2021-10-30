package index

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowIndexPage(c *gin.Context) {

	// Call the render function with the name of the template to render
	c.HTML(http.StatusOK, "index.html", gin.H{})
}