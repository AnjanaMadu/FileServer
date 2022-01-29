package main

import (
	"os"

	"FileServer/modules"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initalize
	e := echo.New()
	port := os.Getenv("PORT")
	os.Mkdir("downloads", os.ModePerm)

	// Routes
	e.GET("/", modules.IndexPage)
	e.GET("/upload", modules.UploadPage)
	e.POST("/api/upload", modules.HandleUpload)
	e.GET("/dl/id/:id", modules.DownloadFile)
	e.GET("/dl/name/:name", modules.DownloadFile)
	e.GET("/files", modules.GetFiles)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
