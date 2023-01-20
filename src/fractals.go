package main

import (
	"fmt"
	"math/rand"

	"github.com/benc-uk/gofract/pkg/colors"
	"github.com/benc-uk/gofract/pkg/fractals"
)

// ================================================================
// Returns a randomised fractal
// ================================================================
func randomFractal(width int) fractals.Fractal {
	fractType := "mandelbrot"
	if rand.Intn(2) == 1 {
		fractType = "julia"
	}

	f := fractals.Fractal{
		FractType:    fractType,
		W:            ratioWidth,
		H:            ratioHeight,
		ImgWidth:     width,
		InnerColor:   "#000000",
		FullScreen:   false,
		ColorRepeats: 5,
	}

	if fractType == "julia" {
		seeds := []fractals.ComplexPair{
			{R: +0.000, I: +0.800},
			{R: +0.370, I: +0.100},
			{R: -0.540, I: +0.540},
			{R: -0.400, I: -0.590},
			{R: +0.340, I: -0.050},
			{R: -0.790, I: +0.150},
			{R: -0.162, I: +1.040},
			{R: +0.280, I: +0.008},
			{R: -1.476, I: +0.000},
		}
		f.JuliaSeed = seeds[rand.Intn(len(seeds))]
		f.JuliaSeed.I += (rand.Float64() * 0.005)
		f.JuliaSeed.R += (rand.Float64() * 0.005)
		f.Center = fractals.ComplexPair{R: 0.0, I: 0.0}
		f.MagFactor = 0.04 + (rand.Float64() * 0.3)
		f.MaxIter = 100/f.MagFactor + 10
	} else {
		if rand.Intn(2) == 1 {
			f.MagFactor = 0.0001 + (rand.Float64() * 0.001)
			f.MaxIter = 1/f.MagFactor + 10
		} else {
			f.MagFactor = 0.01 + (rand.Float64() * 0.1)
			f.MaxIter = 100/f.MagFactor + 10
		}
		f.Center = mandelPoints[rand.Intn(len(mandelPoints))]
	}

	return f
}

// ================================================================
// Returns a query string from a Fractal struct
// ================================================================
func getQueryString(f fractals.Fractal, g colors.GradientTable) string {
	fType := "m"
	if f.FractType == "julia" {
		fType = fmt.Sprintf("j&jr=%f&ji=%f", f.JuliaSeed.R, f.JuliaSeed.I)
	}

	q := fmt.Sprintf("t=%s&w=%d&cr=%f&ci=%f&m=%f&i=%f", fType, f.ImgWidth, f.Center.R, f.Center.I, f.MagFactor, f.MaxIter)

	table := g.GetTable()
	colorString := ""
	for _, c := range table {
		hexColor := c.Col.Hex()
		delim := "_"
		if c == table[len(table)-1] {
			delim = ""
		}
		colorString += fmt.Sprintf("%s,%.1f%s", hexColor[1:], c.Pos, delim)
	}

	q += fmt.Sprintf("&p=%s&pm=%d&r=%d", colorString, g.Mode, f.ColorRepeats)

	return q
}
