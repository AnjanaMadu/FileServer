package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initalize
	e := echo.New()
	port := os.Getenv("PORT")
	os.Mkdir("downloads", os.ModePerm)

	// Routes
	e.GET("/", IndexPage)
	e.GET("/upload", UploadPage)
	e.POST("/upload", HandleUpload)
	e.GET("/download/:id", DownloadFile)
	e.GET("/files", GetFiles)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

// Handler
func IndexPage(c echo.Context) error {
	return c.File("pages/index.html")
}

func UploadPage(c echo.Context) error {
	return c.File("pages/upload.html")
}

var FileIds = make(map[string]string)

func HandleUpload(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	log.Println("Upload:", file.Filename)
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("downloads/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	downloadId := RandomString(6)
	log.Println("Download ID:", downloadId)
	downloadLink := fmt.Sprintf("<a href=\"/download/%s\">%s</a>", downloadId, file.Filename)
	FileIds[downloadId] = file.Filename

	return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Your file uploaded!</h2><p>File name: %s<br>Download Link: %s</p>", file.Filename, downloadLink))
}

func DownloadFile(c echo.Context) error {
	fileId := c.Param("id")

	for id, name := range FileIds {
		if fileId == id {

			filePath := fmt.Sprintf("downloads/%s", id)
			log.Println("Download:", filePath)
			return c.Attachment(filePath, name)

		}
	}
	return c.String(http.StatusNotFound, "File not found")
}

func GetFiles(c echo.Context) error {
	flist := make([]string, 0)
	files, err := ioutil.ReadDir("downloads")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		flist = append(flist, f.Name())
	}
	return c.String(http.StatusOK, fmt.Sprintf("<h2>Files</h2><p>%s</p>", strings.Join(flist, "<br>")))
}

func RandomString(count int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, count)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
