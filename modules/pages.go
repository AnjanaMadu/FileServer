package modules

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func IndexPage(c echo.Context) error {
	return c.File("html/index.html")
}

func UploadPage(c echo.Context) error {
	return c.File("html/upload.html")
}

func GetFiles(c echo.Context) error {
	flist := make([]string, 0)
	files, err := ioutil.ReadDir("downloads")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		ahref := fmt.Sprintf("<a href=\"/dl/name/%s\">%s</a>", f.Name(), f.Name())
		flist = append(flist, ahref)
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Files</h2><p>%s</p>", strings.Join(flist, "<br>")))
}
