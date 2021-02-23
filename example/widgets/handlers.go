package widgets

import (
	"github.com/gin-gonic/gin"
)

type WidgetQuery struct {
	Limit int `json:"limit" form:"limit"`
}

func ListHandler(c *gin.Context) {
	var query WidgetQuery
	if err := c.BindQuery(&query); err != nil {
		return
	}
	c.JSON(200, query)
}

type WidgetCreate struct {
	Name string `json:"name"`
}

func CreateHandler(c *gin.Context) {
	var widget WidgetCreate
	if err := c.BindJSON(&widget); err != nil {
		return
	}
	c.JSON(200, widget)
}

type WidgetUpdate struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func UpdateHandler(c *gin.Context) {
	var widget WidgetCreate
	if err := c.BindJSON(&widget); err != nil {
		return
	}
	c.JSON(200, widget)
}

func DeleteHandler(c *gin.Context) {
	c.JSON(200, c.Param("id"))
}
