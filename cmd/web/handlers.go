package main

import (
	"AsciiArtWebExport/pkg/AsciiArt"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/":
		app.notFound(w)
		return
	case r.Method != http.MethodGet:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	os.Create("output.txt")

	app.render(w, r, "home.page.html", nil)
}

func (app *application) asciiArtWeb(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/ascii-art-web":
		app.notFound(w)
		return
	case r.Method != http.MethodPost:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	banner := r.FormValue("banner")
	rawInput := r.FormValue("rawInput")

	AsciiOutput, err := AsciiArt.AsciiArt(rawInput, banner)
	if err != nil {
		app.serverError(w, err)
		return
	}

	td := &templateData{
		AsciiOutput: AsciiOutput,
	}

	f, err := os.Create("output.txt")
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(AsciiOutput))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "ascii-art.page.html", td)
}

func (app *application) download(w http.ResponseWriter, r *http.Request) {
	// if err := r.ParseForm(); err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	f, err := os.Open("output.txt")
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer f.Close()

	contentDisposition := fmt.Sprintf("attachment; filename=%s", f.Name())
	w.Header().Set("Content-Disposition", contentDisposition)

	contentType := "text/plain; charset=utf-8"
	w.Header().Set("Content-Type", contentType)

	finfo, err := f.Stat()
	if err != nil {
		app.serverError(w, err)
	}

	contentLength := strconv.Itoa(int(finfo.Size()))
	w.Header().Set("Content-Length", contentLength)

	if _, err := io.Copy(w, f); err != nil {
		app.serverError(w, err)
		return
	}
}
