package main

import (
	"fmt"
	"github.com/murlokswarm/app"
	"strconv"
	"time"
)

const (
	// Task statuses
	NEW         = "New"
	IN_PROGRESS = "In Progress"
	DONE        = "Done"

	// Durations in minutes
	TASK_LENGTH       = 1
	PAUSE_LENGTH      = 5
	LONG_PAUSE_LENGTH = 15
)

type Task struct {
	Desc   string
	Status string
}

type TaskBox struct {
	Task
	TaskTimer
	ShowInput bool
}

type TaskTimer struct {
	Minutes string
	Seconds string
}

var (
	task    = &Task{}
	taskBox = &TaskBox{ShowInput: true}
	stop    = make(chan bool, 1)
)

func (task *Task) Render() string {
	return `
    <div class="Task">
        <h4>
            Enter your next task:
        </h4>
        <input type="text" placeholder="Next thing to do..." onchange="OnInputChange" />
    </div>
    `
}

func (taskBox *TaskBox) Render() string {
	return `
    <div class="TaskBox">
      {{if eq .ShowInput true}}<Task></Task>{{else}}
      <span oncontextmenu="OnContextMenu"><span id="taskDesc">{{html .Task.Desc}}</span>:<span id="status">{{html .Status}}</span></span>
      <button onclick="OnStart">Start</button>
      <button onclick="OnStop">Stop</button>
      <button onclick="OnRemove">Remove</button>
      {{end}}
      {{if eq .Status "In Progress"}}<div class="TaskTimer"><span>{{html .TaskTimer.Minutes}}:{{html .TaskTimer.Seconds}}</span></div>{{end}}
    </div>
`
}

func (taskBox *TaskBox) OnContextMenu() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&TaskMenu{})
}

func (taskBox *TaskBox) OnStart() {
	startTask()
}

func (taskBox *TaskBox) OnStop() {
	stopTask()
}

func (taskBox *TaskBox) OnRemove() {
	removeTask()
}

func (task *Task) OnInputChange(arg app.ChangeArg) {
	task.Desc = arg.Value
	task.Status = NEW
	taskBox.ShowInput = false
	taskBox.Desc = task.Desc
	taskBox.Status = NEW

	app.Render(taskBox)
}

func startTask() {
	if taskBox.Status == NEW {
		taskBox.Status = IN_PROGRESS
		dock.SetIcon("resources/pomodoro.png")

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
						taskBox.TaskTimer.Minutes = strconv.Itoa(TASK_LENGTH)
						taskBox.TaskTimer.Seconds = "00"
						dock.SetIcon("resources/pomodoro.png")
						break ticker
					}

				default:
					if elapsed > TASK_LENGTH*60 {
						taskBox.Status = DONE
						taskBox.TaskTimer.Minutes = "00"
						taskBox.TaskTimer.Seconds = "00"
						dock.SetIcon("resources/pomodoro_done.png")
						break ticker
					}
				}

				mins := (TASK_LENGTH*60 - elapsed) / 60

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

func stopTask() {
	if taskBox.Task.Status == IN_PROGRESS {
		stop <- true
	}
}

func removeTask() {
	if taskBox.Task.Status == IN_PROGRESS {
		stop <- true
	}
	task.Status = NEW
	task.Desc = ""
	taskBox.ShowInput = true

	app.Render(taskBox)
}
