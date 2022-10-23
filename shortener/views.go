package shortener

import (
	"2dal/shortener/models"
	"net/http"

	"github.com/gin-gonic/gin"

	qrcode "github.com/skip2/go-qrcode"
)

////// CREATE LINK ///////

type CreateLinkStruct struct {
	Url   string `form:"url" binding:"required,url,startswith=http,contains=://,max=255"`
	Alias string `form:"alias" binding:"omitempty,max=15,isUnique=alias@link"`
}

func CreateLink(c *gin.Context) {
	// Validate form data
	data := c.MustGet(gin.BindKey).(*CreateLinkStruct)

	user := models.User{Username: "ndasim", Subscription: "sd"}

	link := models.Link{}
	link.Create(data.Url, data.Alias, &user)

	c.IndentedJSON(http.StatusOK, map[string]string{
		"origin_url":  link.Origin_url,
		"short_url":   link.Alias,
		"valid_until": link.To_ts,
	})
}

////// CREATE QR ///////

func CreateQR(c *gin.Context) {
	qrcode.WriteFile("https://example.org", qrcode.Medium, -50, "cached/qr.png")

	c.File("cached/qr.png")
}
