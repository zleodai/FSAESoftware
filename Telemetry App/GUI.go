package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("FSAE Telemetry")
	a.Settings().SetTheme(theme.DarkTheme())
	w.SetMaster()
	w.Resize(fyne.NewSize(1280, 720))

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}