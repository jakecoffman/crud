package widgets

import (
	"github.com/gin-gonic/gin"
)

func ListHandler(c *gin.Context) {
	c.JSON(200, c.Request.URL.Query())
}

func CreateHandler(c *gin.Context) {
	var widget interface{}
	if err := c.BindJSON(&widget); err != nil {
		return
	}
	c.JSON(200, widget)
}

func GetHandler(c *gin.Context) {
	c.JSON(200, c.Params)
}

func UpdateHandler(c *gin.Context) {
	var widget interface{}
	if err := c.BindJSON(&widget); err != nil {
		return
	}
	c.JSON(200, widget)
}

func DeleteHandler(c *gin.Context) {
	c.JSON(200, c.Params)
}
