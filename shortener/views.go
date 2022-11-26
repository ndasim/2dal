package shortener

import (
	"2dal/shortener/models"
	"2dal/shortener/render"
	"image/png"
	"net/http"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	//qrcode "github.com/skip2/go-qrcode"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	pngqr "github.com/boombuler/barcode/qr"
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

	// Prevent BASIC users to create link with alias
	if user.Subscription == "BASIC" && data.Alias != "" {
		c.JSON(http.StatusBadRequest, gin.H{"auth": "subscription level is not sufficient"})
		return
	}

	// First look for the existing origin url on the db, resources are valuable
	existErr := link.FindOrigin(data.Url)
	if existErr != nil {
		err = link.Create(data.Url, data.Alias, &user)

		if err != nil {
			sentry.CaptureException(err)
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	} else {
		if data.Alias != "" && user.Username != link.User.Username {
			c.JSON(http.StatusUnauthorized, gin.H{"permission": "You can only override your short urls"})
			return
		} else {
			link.Update(data.Alias, &user)
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

	layout := "2006-01-02T15:04:05Z"
	str := link.To_ts
	to_ts, err := time.Parse(layout, str)
	today := time.Now()

	if err != nil {
		sentry.CaptureException(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		panic(err)
	}

	if today.After(to_ts) {
		c.HTML(http.StatusOK, "expired.html", gin.H{
			"url": link.Origin_url,
		})
	} else {
		c.HTML(http.StatusOK, "forwarder.html", gin.H{
			"url": link.Origin_url,
		})
	}
}

// //// CREATE QR ///////
type QRStruct struct {
	Alias string `form:"alias" uri:"alias" binding:"required,max=15"`
}

func CreateSvgQR(c *gin.Context) { // Validate form data
	data := OpenLinkStruct{}
	if err := c.BindUri(&data); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var w http.ResponseWriter = c.Writer
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)

	// Create the barcode
	qrCode, _ := qr.Encode("2d.al/"+data.Alias, qr.M, qr.Auto)

	// Write QR code to SVG
	qs := render.NewQrSVG(qrCode, 20)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	s.End()
	//qrcode.WriteFile("https://example.org", qrcode.Medium, -50, "cached/qr.png")

	//c.File("cached/qr.png")
}

func CreatePngQR(c *gin.Context) {
	data := OpenLinkStruct{}
	if err := c.BindUri(&data); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var w http.ResponseWriter = c.Writer

	qrCode, _ := pngqr.Encode(data.Alias, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	png.Encode(w, qrCode)
}
