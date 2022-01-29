package modules

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var FileIds = make(map[string]string)

func DownloadFile(c echo.Context) error {
	fileId := c.Param("id")
	mode := c.QueryParam("mode")

	if mode == "name" {
		filePath := fmt.Sprintf("downloads/%s", fileId)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "File not found")
		}
		return c.Attachment(filePath, fileId)
	}

	for id, name := range FileIds {
		if fileId == id {
			filePath := fmt.Sprintf("downloads/%s", name)
			return c.Attachment(filePath, name)
		}
	}
	return c.String(http.StatusNotFound, "File not found")
}

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
	shortLink := fmt.Sprintf("<a href=\"/download/%s\">%s</a>", downloadId, file.Filename)
	longLink := fmt.Sprintf("<a href=\"/download/%s?mode=name\">%s</a>", file.Filename, file.Filename)
	FileIds[downloadId] = file.Filename

	return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Your file uploaded!</h2><p>File name: %s<br>Download Links:<br>	%s (short link)<br>		%s (long link)</p>", file.Filename, shortLink, longLink))
}
