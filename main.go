package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initalize
	e := echo.New()
	port := os.Getenv("PORT")
	os.Mkdir("downloads", os.ModePerm)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", IndexPage)
	e.GET("/upload", UploadPage)
	e.POST("/upload", HandleUpload)
	e.GET("/download/:id", DownloadFile)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

// Handler
func IndexPage(c echo.Context) error {
	return c.HTML(http.StatusOK, "pages/index.html")
}

func UploadPage(c echo.Context) error {
	return c.HTML(http.StatusOK, "pages/upload.html")
}

var FileIds map[string]string

func HandleUpload(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	downloadId := RandomString(6)
	downloadLink := fmt.Sprintf("<a href=\"/download/%s\">%s</a>", downloadId, file.Filename)
	FileIds[downloadId] = "downloads/" + file.Filename

	return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Your file uploaded!</h2><p>File name: %s<br>Download Link: %s</p>", file.Filename, downloadLink))
}

func DownloadFile(c echo.Context) error {
	id := c.Param("id")
	for k, v := range FileIds {
		if k == id {
			return c.File(v)
		}
	}
	return c.String(http.StatusNotFound, "File not found")
}

func RandomString(count int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, count)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
