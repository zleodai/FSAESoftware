package main

import (
	"fmt"
	"image/color"
	"time"

	"databaseAPI"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	throttleData     []float64 = make([]float64, 0, maxDataPoints)
	brakeData        []float64
	steeringData     []float64
	xAccelData       []float64
	yAccelData       []float64
	zAccelData       []float64
	tireTempData     []float64
	tirePressureData []float64
	pitchData        []float64
	yawData          []float64
	rollData         []float64

	throttle      float64
	brake         float64
	steeringAngle float64
	xAccel        float64
	yAccel        float64
	zAccel        float64
	tireTemp      float64
	tirePressure  float64
	pitch         float64
	yaw           float64
	roll          float64
)

const (
	maxDataPoints = 400
	graphHeight   = 100
	graphWidth    = 1000
)

func main() {
	connection := databaseAPI.NewConnection()
	databaseAPI.InsertIntoPool(connection, []databaseAPI.TelemetryPacket{databaseAPI.TempTelemtryPacket()})

	a := app.New()
	w := a.NewWindow("Vehicle Telemetry Monitor")
	w.Resize(fyne.NewSize(1080, 1920))

	throttleLabel := widget.NewLabel("Throttle: 0%")
	throttleGraph := container.NewWithoutLayout()
	brakeLabel := widget.NewLabel("Brake: 0%")
	brakeGraph := container.NewWithoutLayout()
	steeringLabel := widget.NewLabel("Steering: 0°")
	steeringGraph := container.NewWithoutLayout()

	xAccelLabel := widget.NewLabel("X Acceleration: 0 m/s²")
	xAccelGraph := container.NewWithoutLayout()
	yAccelLabel := widget.NewLabel("Y Acceleration: 0 m/s²")
	yAccelGraph := container.NewWithoutLayout()
	zAccelLabel := widget.NewLabel("Z Acceleration: 0 m/s²")
	zAccelGraph := container.NewWithoutLayout()

	tireTempLabel := widget.NewLabel("Tire Temperature: 0°F")
	tireTempGraph := container.NewWithoutLayout()
	tirePressureLabel := widget.NewLabel("Tire Pressure: 0 psi")
	tirePressureGraph := container.NewWithoutLayout()

	pitchLabel := widget.NewLabel("Pitch: 0°")
	pitchGraph := container.NewWithoutLayout()
	yawLabel := widget.NewLabel("Yaw: 0°")
	yawGraph := container.NewWithoutLayout()
	rollLabel := widget.NewLabel("Roll: 0°")
	rollGraph := container.NewWithoutLayout()

	updateTelemetryLabels := func() {
		throttleLabel.SetText("Throttle: " + fmt.Sprintf("%.2f", throttle) + "%")
		brakeLabel.SetText("Brake: " + fmt.Sprintf("%.2f", brake) + "%")
		steeringLabel.SetText("Steering: " + fmt.Sprintf("%.2f", steeringAngle) + "°")
		xAccelLabel.SetText("X Acceleration: " + fmt.Sprintf("%.2f", xAccel) + " m/s²")
		yAccelLabel.SetText("Y Acceleration: " + fmt.Sprintf("%.2f", yAccel) + " m/s²")
		zAccelLabel.SetText("Z Acceleration: " + fmt.Sprintf("%.2f", zAccel) + " m/s²")
		tireTempLabel.SetText("Tire Temperature: " + fmt.Sprintf("%.2f", tireTemp) + "°F")
		tirePressureLabel.SetText("Tire Pressure: " + fmt.Sprintf("%.2f", tirePressure) + " psi")
		pitchLabel.SetText("Pitch: " + fmt.Sprintf("%.2f", pitch) + "°")
		yawLabel.SetText("Yaw: " + fmt.Sprintf("%.2f", yaw) + "°")
		rollLabel.SetText("Roll: " + fmt.Sprintf("%.2f", roll) + "°")
	}

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyW:
			throttle += 5
			if throttle > 100 {
				throttle = 100
			}
		case fyne.KeyS:
			throttle -= 5
			if throttle < 0 {
				throttle = 0
			}
		case fyne.KeyE:
			brake += 5
			if brake > 100 {
				brake = 100
			}
		case fyne.KeyD:
			brake -= 5
			if brake < 0 {
				brake = 0
			}
		case fyne.KeyR:
			steeringAngle += 10
			if steeringAngle > 180 {
				steeringAngle = 180
			}
		case fyne.KeyF:
			steeringAngle -= 10
			if steeringAngle < -180 {
				steeringAngle = -180
			}
		}
		updateTelemetryLabels()
	})

	go func() {
		for range time.Tick(time.Millisecond * 20) {
			updateData(&throttleData, throttle)
			updateData(&brakeData, brake)
			updateData(&steeringData, steeringAngle)
			updateData(&xAccelData, xAccel)
			updateData(&yAccelData, yAccel)
			updateData(&zAccelData, zAccel)
			updateData(&tireTempData, tireTemp)
			updateData(&tirePressureData, tirePressure)
			updateData(&pitchData, pitch)
			updateData(&yawData, yaw)
			updateData(&rollData, roll)

			drawGraph(throttleGraph, throttleData, 100, 0)
			drawGraph(brakeGraph, brakeData, 100, 0)
			drawGraph(steeringGraph, steeringData, 180, -180)
			drawGraph(xAccelGraph, xAccelData, 30, -5)
			drawGraph(yAccelGraph, yAccelData, 15, -15)
			drawGraph(zAccelGraph, zAccelData, 10, -10)
			drawGraph(tireTempGraph, tireTempData, 200, 0)
			drawGraph(tirePressureGraph, tirePressureData, 50, 0)
			drawGraph(pitchGraph, pitchData, 100, -100)
			drawGraph(yawGraph, yawData, 100, -100)
			drawGraph(rollGraph, rollData, 100, -100)
		}
	}()

	inputGraphs := container.NewVBox(
		container.NewVBox(
			throttleLabel,
			addGraphWithBackground(throttleGraph, 100, 0),
		),
		container.NewVBox(
			brakeLabel,
			addGraphWithBackground(brakeGraph, 100, 0),
		),
		container.NewVBox(
			steeringLabel,
			addGraphWithBackground(steeringGraph, 180, -180),
		),
	)

	accelGraphs := container.NewVBox(
		container.NewVBox(
			xAccelLabel,
			addGraphWithBackground(xAccelGraph, 30, -5),
		),
		container.NewVBox(
			yAccelLabel,
			addGraphWithBackground(yAccelGraph, 15, -15),
		),
		container.NewVBox(
			zAccelLabel,
			addGraphWithBackground(zAccelGraph, 10, -10),
		),
	)

	tireGraphs := container.NewVBox(
		container.NewVBox(
			tireTempLabel,
			addGraphWithBackground(tireTempGraph, 200, 0),
		),
		container.NewVBox(
			tirePressureLabel,
			addGraphWithBackground(tirePressureGraph, 50, 0),
		),
	)

	gyroGraphs := container.NewVBox(
		container.NewVBox(
			pitchLabel,
			addGraphWithBackground(pitchGraph, 100, -100),
		),
		container.NewVBox(
			yawLabel,
			addGraphWithBackground(yawGraph, 100, -100),
		),
		container.NewVBox(
			rollLabel,
			addGraphWithBackground(rollGraph, 100, -100),
		),
	)

	content := container.NewVBox()
	buttons := container.NewHBox(
		widget.NewButton("Inputs", func() {
			content.Objects = []fyne.CanvasObject{inputGraphs}
			content.Refresh()
		}),
		widget.NewButton("Acceleration", func() {
			content.Objects = []fyne.CanvasObject{accelGraphs}
			content.Refresh()
		}),
		widget.NewButton("Tires", func() {
			content.Objects = []fyne.CanvasObject{tireGraphs}
			content.Refresh()
		}),
		widget.NewButton("Gyro", func() {
			content.Objects = []fyne.CanvasObject{gyroGraphs}
			content.Refresh()
		}),
	)

	content.Add(inputGraphs)
	mainContent := container.NewVBox(buttons, content)

	w.SetContent(mainContent)
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
	graphColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
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
		label := canvas.NewText(fmt.Sprintf("%.0f", minScale+(maxScale-minScale)*float64(i)/5), color.RGBA{R: 255, G: 0, B: 0, A: 255})
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
