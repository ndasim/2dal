package shortener

import (
	"2dal/shortener/models"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	qrcode "github.com/skip2/go-qrcode"
)

type RapidAPIHeaders struct {
	User         string `header:"X-RapidAPI-User" binding:"required"`
	Subscription string `header:"X-RapidAPI-Subscription" binding:"required"`
	Ip           string `header:"X-Forwarded-For" binding:"required"`
}

////// CREATE LINK ///////

type CreateLinkStruct struct {
	Url   string `form:"url" binding:"required,url,startswith=http,contains=://,max=255"`
	Alias string `form:"alias" binding:"omitempty,max=15,isUnique=alias@link"`
}

func CreateLink(c *gin.Context) {
	// Validate header
	header := RapidAPIHeaders{}
	err := c.ShouldBindHeader(&header)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Validate form data
	data := c.MustGet(gin.BindKey).(*CreateLinkStruct)

	user := models.User{Username: header.User, Subscription: header.Subscription, IpAddress: header.Ip}

	link := models.Link{}

	// First look for the existing origin url on the db, resources are valueble
	existErr := link.FindOrigin(data.Url)
	if existErr != nil {
		err = link.Create(data.Url, data.Alias, &user)

		if err != nil {
			sentry.CaptureException(err)
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	}

	c.IndentedJSON(http.StatusOK, map[string]string{
		"origin_url":  link.Origin_url,
		"short_url":   "2d.al/" + link.Alias,
		"valid_until": link.To_ts,
	})
}

////// OPEN LINK ///////

type OpenLinkStruct struct {
	Alias string `form:"alias" uri:"alias" binding:"required,max=15"`
}

func OpenLink(c *gin.Context) {
	// Validate form data
	data := OpenLinkStruct{}
	if err := c.BindUri(&data); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	link := models.Link{}

	/// Fetch result
	err := link.FindLink(data.Alias)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	/// If not found
	if link.Origin_url == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "forwarder.html", gin.H{
		"url": link.Origin_url,
	})
}

////// CREATE QR ///////

func CreateQR(c *gin.Context) {
	qrcode.WriteFile("https://example.org", qrcode.Medium, -50, "cached/qr.png")

	c.File("cached/qr.png")
}
