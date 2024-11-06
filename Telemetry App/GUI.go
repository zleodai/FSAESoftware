package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	speedData        []float64
	steeringData     []float64
	accelerationData []float64
	speed            float64
	steeringAngle    float64
	acceleration     float64
)

const (
	maxDataPoints = 400
	graphHeight   = 100
	graphWidth    = 1000
)

func main() {
	a := app.New()
	w := a.NewWindow("Vehicle Telemetry Monitor")
	w.Resize(fyne.NewSize(1080, 1920))
	speedLabel := widget.NewLabel("Speed: 0")
	speedGraph := container.NewWithoutLayout()
	steeringLabel := widget.NewLabel("Steering Angle: 0")
	steeringGraph := container.NewWithoutLayout()
	accelLabel := widget.NewLabel("Acceleration: 0")
	accelGraph := container.NewWithoutLayout()
	updateTelemetryLabels := func() {
		speedLabel.SetText("Speed: " + fmt.Sprintf("%.2f", speed))
		steeringLabel.SetText("Steering Angle: " + fmt.Sprintf("%.2f", steeringAngle))
		accelLabel.SetText("Acceleration: " + fmt.Sprintf("%.2f", acceleration))
	}

	// FOR TESTING WITH DUMMY DATA USING KEYBOARD INPUTS
	//REMEMBER WHEN TAKING IN REAL INPUTS I WILL NEED TO
	//CLAMP DATA IF IT GOES OFF GRAPH OR JUST ADJUST GRAPH
	//BUT STILL CLAMP DATA JUST INCASE IDK
	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyW:
			speed += 5
			if speed > 100 {
				speed = 100
			}
		case fyne.KeyS:
			speed -= 5
			if speed < 0 {
				speed = 0
			}
		case fyne.KeyE:
			steeringAngle += 5
			if steeringAngle > 90 {
				steeringAngle = 90
			}
		case fyne.KeyD:
			steeringAngle -= 5
			if steeringAngle < -90 {
				steeringAngle = -90
			}
		case fyne.KeyR:
			acceleration += 1
			if acceleration > 50 {
				acceleration = 50
			}
		case fyne.KeyF:
			acceleration -= 1
			if acceleration < -10 {
				acceleration = -10
			}
		}
		updateTelemetryLabels()
	})

	go func() {
		for range time.Tick(time.Millisecond * 20) {
			updateData(&speedData, speed)
			updateData(&steeringData, steeringAngle)
			updateData(&accelerationData, acceleration)

			updateTelemetryLabels()
			drawGraph(speedGraph, speedData, 100, 0)
			drawGraph(steeringGraph, steeringData, 90, -90)
			drawGraph(accelGraph, accelerationData, 50, -10)
		}
	}()

	content := container.NewVBox(
		container.NewVBox(
			speedLabel,
			addGraphWithBackground(speedGraph, 100, 0),
			widget.NewSeparator(),
		),
		container.NewVBox(
			steeringLabel,
			addGraphWithBackground(steeringGraph, 90, -90),
			widget.NewSeparator(),
		),
		container.NewVBox(
			accelLabel,
			addGraphWithBackground(accelGraph, 10, -10),
		),
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func updateData(data *[]float64, value float64) {
	if len(*data) >= maxDataPoints {
		*data = (*data)[1:]
	}
	*data = append(*data, value)
}

func drawGraph(container *fyne.Container, data []float64, maxScale, minScale float64) {
	container.Objects = nil
	graphColor := color.RGBA{R: 255, G: 0, B: 0, A: 255} //graph line color set here
	for i := 1; i < len(data); i++ {
		x1 := float64(i-1) * graphWidth / maxDataPoints
		y1 := graphHeight - (data[i-1]-minScale)/(maxScale-minScale)*graphHeight
		x2 := float64(i) * graphWidth / maxDataPoints
		y2 := graphHeight - (data[i]-minScale)/(maxScale-minScale)*graphHeight
		line := canvas.NewLine(graphColor)
		line.StrokeWidth = 2
		line.Position1 = fyne.NewPos(float32(x1), float32(y1))
		line.Position2 = fyne.NewPos(float32(x2), float32(y2))
		container.Add(line)
	}
	container.Refresh()
}

func addGraphWithBackground(graph *fyne.Container, maxScale, minScale float64) *fyne.Container {
	bg := canvas.NewRectangle(color.RGBA{R: 220, G: 220, B: 220, A: 255})
	bg.Resize(fyne.NewSize(graphWidth, graphHeight))
	valueMarkers := container.NewWithoutLayout(bg)
	for i := 0; i <= 5; i++ {
		y := float32(graphHeight - (float64(i) * graphHeight / 5))
		line := canvas.NewLine(color.RGBA{R: 180, G: 180, B: 180, A: 255})
		line.Position1 = fyne.NewPos(0, y)
		line.Position2 = fyne.NewPos(float32(graphWidth), y)
		valueMarkers.Add(line)
		label := canvas.NewText(fmt.Sprintf("%.0f", minScale+(maxScale-minScale)*float64(i)/5), color.RGBA{R: 255, G: 0, B: 0, A: 255}) //marker text color set here
		label.TextSize = 10
		label.Move(fyne.NewPos(5, y-7))
		valueMarkers.Add(label)
	}

	paddedGraph := container.NewVBox(
		container.NewWithoutLayout(valueMarkers, graph),
		widget.NewLabel(" "),
		widget.NewLabel(" "),
		widget.NewLabel(" "),
	)
	return paddedGraph
}
