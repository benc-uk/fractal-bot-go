package main

import (
	"encoding/json"
	"fractal-bot-go/pkg/functions"
	"fractal-bot-go/pkg/twitter"
	"image"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/benc-uk/gofract/pkg/colors"
	"github.com/benc-uk/gofract/pkg/fractals"
)

const outBindingName = "res"
const ratioWidth = 4.0
const ratioHeight = 3.0
const ratioHW = ratioHeight / ratioWidth

var defaultGradient = colors.GradientTable{}
var mandelPoints = []fractals.ComplexPair{}
var urlBase = "http://localhost:7071"

func init() {
	defaultGradient.AddToTable("#000762", 0.0)
	defaultGradient.AddToTable("#0B48C3", 0.2)
	defaultGradient.AddToTable("#ffffff", 0.4)
	defaultGradient.AddToTable("#E3A000", 0.5)
	defaultGradient.AddToTable("#000762", 0.9)
	if ub := os.Getenv("URL_BASE"); ub != "" {
		urlBase = ub
	}

	file := "points.json"
	jsonFile, _ := os.Open(file)
	defer jsonFile.Close()
	bytes, _ := io.ReadAll(jsonFile)

	_ = json.Unmarshal(bytes, &mandelPoints)

	// use time to seed
	rand.Seed(time.Now().UnixNano())
}

// ================================================================
// Tweets a fractal on a schedule
// ================================================================
func tweetHandler(w http.ResponseWriter, r *http.Request) {
	width := 1200
	var resp functions.InvokeResponse
	success := true

	f := randomFractal(width)
	log.Printf("### 🎨 Rendering fractal: %s", f.FractType)

	img := image.NewRGBA(image.Rect(0, 0, width, int(float64(width)*ratioHW)))

	// Randomise the gradient colors
	gt := colors.GradientTable{}
	gt.Randomise()

	// Randomise the gradient mode
	gt.Mode = rand.Intn(colors.MaxColorModes)

	f.Render(img, gt)

	url := urlBase + "/api/fractal?" + getQueryString(f, gt)
	log.Printf("### 🎺 URL: %s", url)

	id, err := twitter.UploadMediaImage(img)
	if err != nil {
		log.Printf("!!! Error uploading image: %v", err)
		success = false
	}
	log.Printf("### 👍 Uploaded media to twitter: %s", *id)

	err = twitter.SendTweet(url, id)
	if err != nil {
		log.Printf("!!! Error sending tweet: %v", err)
		success = false
	} else {
		log.Printf("### 🐦 Tweet successful")
	}

	w.Header().Set("Content-Type", "application/json")
	if success {
		resp = functions.NewInvokeResponse(outBindingName, "Tweeted fractal", "Fractal generated and tweeted. URL is "+url)
		w.WriteHeader(http.StatusOK)
	} else {
		resp = functions.NewInvokeResponse(outBindingName, "Error", "Error generating fractal")
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// ================================================================
// Returns a custom fractal PNG image over HTTP
// ================================================================
func fractalHandler(w http.ResponseWriter, r *http.Request) {
	width := 1000
	if p := r.URL.Query().Get("w"); p != "" {
		width, _ = strconv.Atoi(p)
	}
	centerR := -0.6
	if p := r.URL.Query().Get("cr"); p != "" {
		centerR, _ = strconv.ParseFloat(p, 64)
	}
	centerI := 0.0
	if p := r.URL.Query().Get("ci"); p != "" {
		centerI, _ = strconv.ParseFloat(p, 64)
	}
	mag := 1.0
	if p := r.URL.Query().Get("m"); p != "" {
		mag, _ = strconv.ParseFloat(p, 64)
	}
	iters := 90.0
	if p := r.URL.Query().Get("i"); p != "" {
		iters, _ = strconv.ParseFloat(p, 64)
	}
	repeats := 2
	if p := r.URL.Query().Get("r"); p != "" {
		repeats, _ = strconv.Atoi(p)
	}
	fractType := "mandelbrot"
	if p := r.URL.Query().Get("t"); p != "" {
		if p == "j" {
			fractType = "julia"
		}
		if p == "m" {
			fractType = "mandelbrot"
		}
	}
	juliaR := 0.355
	if p := r.URL.Query().Get("jr"); p != "" {
		juliaR, _ = strconv.ParseFloat(p, 64)
	}
	juliaI := 0.355
	if p := r.URL.Query().Get("ji"); p != "" {
		juliaI, _ = strconv.ParseFloat(p, 64)
	}

	paletteString := ""
	var gradient colors.GradientTable
	gradient = defaultGradient
	if p := r.URL.Query().Get("p"); p != "" {
		paletteString = p
		gradient = colors.GradientTable{}
	}
	blendMode := 0
	if p := r.URL.Query().Get("pm"); p != "" {
		blendMode, _ = strconv.Atoi(p)
	}

	parts := strings.Split(paletteString, "_")
	if len(parts) == 0 {
		gradient = defaultGradient
	} else {
		for _, part := range parts {
			pair := strings.Split(part, ",")
			if len(pair) != 2 {
				gradient = defaultGradient
				break
			}

			pos, err := strconv.ParseFloat(pair[1], 64)
			if err != nil {
				gradient = defaultGradient
				break
			}
			gradient.AddToTable("#"+pair[0], pos)
		}
	}
	gradient.Mode = blendMode

	height := int(float64(width) * ratioHW)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{width, height},
	})

	f := fractals.Fractal{
		FractType:    fractType,
		Center:       fractals.ComplexPair{R: centerR, I: centerI},
		MagFactor:    mag,
		MaxIter:      iters,
		W:            ratioWidth,
		H:            ratioHeight,
		ImgWidth:     width,
		JuliaSeed:    fractals.ComplexPair{R: juliaR, I: juliaI},
		InnerColor:   "#000000",
		FullScreen:   false,
		ColorRepeats: repeats,
	}

	f.Render(img, gradient)
	_ = png.Encode(w, img)
}
