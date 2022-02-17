package handlers

import (
	"GolangStore/models"

	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) {
	prod := models.ProductFinder()
	// Call the render function with the name of the template to render
	render(c, gin.H{"title": "Home Page", "payload": prod}, "products.html")
	//ctx.HTML(http.StatusOK, "products", models.ProductFinder())
}
