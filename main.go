package main

import (
	"github.com/murlokswarm/app"
	_ "github.com/murlokswarm/mac"
)

var (
	win app.Windower
)

func init() {
	app.RegisterComponent(task)
	app.RegisterComponent(taskBox)
}

func main() {
	app.OnLaunch = func() {
		win = app.NewWindow(app.Window{
			Title:          "Pomodoro2Go Timer",
			Width:          400,
			Height:         50,
			TitlebarHidden: false,
			Vibrancy:       app.VibeMediumLight,
		})

		win.Mount(taskBox)
	}

	app.Run()
}
