package main

import (
	"github.com/murlokswarm/app"
	_ "github.com/murlokswarm/mac"
)

var (
	win  app.Windower
	dock app.Docker
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
			Height:         150,
			TitlebarHidden: false,
			Vibrancy:       app.VibeMediumLight,
			OnClose: func() bool {
				return true
			},
		})

		ok := false
		if dock, ok = app.Dock(); ok {
			dock.SetIcon("resources/pomodoro.png")
		}

		win.Mount(taskBox)
	}

	app.OnFinalize = func() {
		if taskBox.Activity.Status == IN_PROGRESS {
			stop <- true
		}
	}

	app.OnTerminate = func() bool {
		return true
	}

	app.Run()
}
