package main

import (
	"AsciiArtWebExport/pkg/AsciiArt"
	"net/http"
	"os"
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
