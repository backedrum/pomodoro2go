package main

import (
	"fmt"
	"github.com/murlokswarm/app"
	"time"
)

type TaskMenu struct {
}

var stop = make(chan bool, 1)

func (tm *TaskMenu) Render() string {
	return `
   <menu label="TaskMenu">
      <menuitem label="Start" shortcut="meta+s" onclick="OnStart"/>
      <menuitem label="Stop" shortcut="meta+z" onclick="OnStop"/>
      <menuitem label="New task" shortcut="meta+n" onclick="OnNewTask"/>
   </menu>
`
}

func (tm *TaskMenu) OnStart() {
	if taskBox.Status == NEW {
		fmt.Println("Starting task..")
		taskBox.Status = IN_PROGRESS

		go func() {
			elapsed := 0
		ticker:
			for range time.Tick(time.Duration(1) * time.Second) {
				elapsed++

				select {
				case stopped := <-stop:
					if stopped {
						defer app.Render(taskBox)
						taskBox.Status = NEW
						taskBox.TaskTimer.Minutes = "25"
						taskBox.TaskTimer.Seconds = "00"
						break ticker
					}

				default:
					if elapsed > 25*60 {
						taskBox.Status = DONE
						taskBox.TaskTimer.Minutes = "00"
						taskBox.TaskTimer.Seconds = "00"
						break ticker
					}
				}

				mins := (25*60 - elapsed) / 60

				var secs int
				remainder := elapsed % 60
				if remainder == 0 {
					secs = 0
				} else {
					secs = 60 - elapsed%60
				}

				taskBox.TaskTimer.Minutes = fmt.Sprintf("%02d", mins)
				taskBox.TaskTimer.Seconds = fmt.Sprintf("%02d", secs)
				app.Render(taskBox)
			}
		}()

	}
}

func (tm *TaskMenu) OnStop() {
	if taskBox.Task.Status == IN_PROGRESS {
		fmt.Println("Stopping current task..")
		stop <- true
	}
}

func (tm *TaskMenu) OnNewTask() {
	fmt.Println("Switch back to task input..")
	stop <- true
	task.Status = NEW
	task.Desc = ""
	taskBox.ShowInput = true

	app.Render(taskBox)
}

func init() {
	app.RegisterComponent(&TaskMenu{})
}
