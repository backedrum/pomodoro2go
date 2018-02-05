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
			OnClose: func() bool {
				return true
			},
		})

		if dock, ok := app.Dock(); ok {
			dock.SetIcon("resources/pomodoro.png")
		}

		win.Mount(taskBox)
	}

	app.OnFinalize = func() {
		if taskBox.Task.Status == IN_PROGRESS {
			stop <- true
		}
	}

	app.OnTerminate = func() bool {
		return true
	}

	app.Run()
}
