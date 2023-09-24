package main

import (
	"2dal/core"
	"2dal/shortener"
	"2dal/shortener/db"
	"2dal/shortener/middleware"
	"2dal/shortener/models"
	"context"
	"embed"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/catinello/base62"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/kamva/mgm/v3"
	"golang.org/x/exp/slices"

	"github.com/getsentry/sentry-go"
)

type contextKey int

const SomeContextKey = contextKey(1)

//go:embed static
var staticFiles embed.FS

// gE#q?6a6-xCpTU
func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://32109d4bacdc44b2bdf09093dc28ff2a@o1306780.ingest.sentry.io/4504068517199872",
		EnableTracing: true,
		// We recommend adjusting these values in production:
		TracesSampleRate: 1.0,
		// The sampling rate for profiling is relative to TracesSampleRate:
		ProfilesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	ctx := context.WithValue(context.Background(), SomeContextKey, "some details about your panic")

	func() {
		defer sentry.RecoverWithContext(ctx)
		// do all of the scary things here
	}()

	err = db.GetConnection()

	if err != nil {
		sentry.CaptureException(err)
	}

	router := gin.Default()

	router.Use(middleware.Errors())
	router.Use(middleware.Options)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/api/create", gin.Bind(shortener.CreateLinkStruct{}), shortener.CreateLink)

	//router.GET("/api/qr/png/:alias", shortener.CreatePngQR)

	router.GET("/:alias", func(ctx *gin.Context) {
		ctx.ShouldBindUri(&shortener.OpenLinkStruct{})
	}, shortener.OpenLink)

	// SVG builder
	router.GET("/:alias/qr", shortener.CreateSvgQR)

	fs := core.EmbedFolder(staticFiles, "static", true)

	router.Use(static.Serve("/", fs))

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("url", core.UrlValidator())
		v.RegisterValidation("isUnique", db.IsUniqueValidator())
	}

	router.Run(":80")

	//mi := i.Uint64()
	//fmt.Println(len(uniques))
	//fmt.Println(len(uniques))
	//fmt.Println(len(uniques))

	func() {
		defer func() {
			err := recover()

			if err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 5)
			}
		}()

		// do all of the scary things here
	}()
}

func someTest() {
	a := big.NewInt(2103)

	m := big.NewInt(100000)

	i := big.Int{}
	i.ModInverse(a, m)

	list := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0}
	uniques := []string{}
	start_time := time.Now().UnixMilli()

	for i := 0; i < 90000; i++ {
		char := base62.Encode(i * 1881 % 100000)

		if slices.Contains(uniques, char) {
			list[7]++
		}

		uniques = append(uniques, char)
		list[len(char)]++

		if i%1000 == 0 {
			fmt.Println("took ", time.Now().UnixMilli()-start_time, " ms for, ", len(uniques), " items")
			start_time = time.Now().UnixMilli()
		}
	}

	fmt.Println(list)
}

func linkTest() {
	user := models.User{Username: "ndasim", Subscription: "sd"}

	link := models.Link{}
	link.Create("https://2d.al/", "", &user)

	link2 := &models.Link{}
	_ = mgm.Coll(link2).First(bson.M{}, link2)
}
