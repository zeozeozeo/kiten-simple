package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/zeozeozeo/kiten"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

var (
	canvas  *kiten.Canvas
	canvas2 *kiten.Canvas = kiten.NewCanvas(50, 50, kiten.BlendAdd)

	globalTime float64 = 0
	frameTime  time.Duration
	width      int = 1280
	height     int = 720

	circleRadius int = 70
	circleX      int = circleRadius
	circleY      int = circleRadius
	circleSpeed  int = 10
	circleDy     int = 1
	circleDx     int = 1
)

func main() {
	log.SetFlags(0)

	driver.Main(func(s screen.Screen) {
		// Create a new window
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Demo",
			Width:  width,
			Height: height,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		// Get screen buffer and texture
		buf, err := s.NewBuffer(image.Pt(width, height))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if buf != nil {
				buf.Release()
				buf = nil
			}
		}()

		// Create a new canvas
		canvas = kiten.CanvasFromImageRGBA(buf.RGBA(), kiten.BlendAdd)

		// Start draw loop
		go drawLoop(buf, w)

		// Window loop
		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case key.Event:
				// Close when the escape button is pressed
				if e.Code == key.CodeEscape {
					return
				}
			}
		}
	})
}

func drawLoop(buf screen.Buffer, w screen.Window) {
	for {
		start := time.Now()
		dt := frameTime.Seconds()
		globalTime += dt
		// fmt.Printf("  %f\r", dt)

		// Draw frame
		draw(dt)

		if buf != nil && w != nil {
			w.Upload(image.Point{}, buf, canvas.Image.Bounds())
		}
		frameTime = time.Since(start)
	}
}

func draw(dt float64) {
	canvas.Fill(color.RGBA{0, 0, 0, 255})

	for x := 0; x < canvas2.Width; x++ {
		for y := 0; y < canvas2.Height; y++ {
			random := uint8(rand.Intn(255 + 1))
			canvas2.SetPixel(x, y, color.RGBA{random, random, random, 255})
		}
	}
	canvas.PutCanvas(800, 300, int(globalTime*100)%200, int(globalTime*170)%400, canvas2)

	// Rectangles
	for i := 0; i < 100-int(globalTime*100)%100; i += 10 {
		canvas.Rect(100+i, 100+i, 500+i, 500+i, color.RGBA{255, 0, 0, 255})
	}

	// Lines
	for i := 0; i < 10*100; i += 100 {
		canvas.Line(0, canvas.Height-i, i+int(globalTime*150)%900, 0, color.RGBA{0, 255, 0, 255})
	}

	canvas.CircleFilled(1280/4, 720/4, int(globalTime*50)%70, color.RGBA{128, 255, 255, 255})
	canvas.CircleOutline(1280-150, 720-150, 70-int(globalTime*50)%70, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255})
	canvas.RectFilled(1280-int(globalTime*100)%100-15, 15, 1280-15, 100, color.RGBA{255, 255, 255, 255})

	// Bouncing circle
	// Collide with walls
	if circleX+circleDx > canvas.Width-circleRadius || circleX+circleDx < circleRadius {
		circleDx = -circleDx
	}
	if circleY+circleDy > canvas.Height-circleRadius || circleY+circleDy < circleRadius {
		circleDy = -circleDy
	}

	circleX += circleDx * circleSpeed
	circleY += circleDy * circleSpeed

	// Draw circle
	canvas.Circle(
		circleX,
		circleY,
		circleRadius,
		color.RGBA{
			uint8(int(globalTime*150) % 255),
			255,
			255,
			255,
		})

	// Text
	text := "Hello World! 123456789@#$%^&*)"
	text = text[:int(globalTime*15)%len(text)+1]
	canvas.Text(text, 15, 300, inconsolata.Regular8x16, color.RGBA{255, 255, 255, 255})

	// FPS text
	canvas.Text(fmt.Sprint(math.Round(1/dt))+" FPS", 16, 16, inconsolata.Bold8x16, color.RGBA{255, 255, 255, 255})

	// Draw path
	path := []image.Point{}
	for i := 1; i < 200; i++ {
		path = append(path, image.Pt(i*2+700, 100+rand.Intn(50)))
	}

	canvas.DrawPath(path, color.RGBA{123, 123, 123, 255})
}
