package main

import (
	"fmt"
	"github.com/murlokswarm/app"
	"strconv"
	"time"
)

const (
	// Activity type
	TASK       = ActivityType("Task")
	PAUSE      = ActivityType("Pause")
	LONG_PAUSE = ActivityType("Long pause")

	// Activity statuses
	NEW         = ActivityStatus("New")
	IN_PROGRESS = ActivityStatus("In Progress")
	DONE        = ActivityStatus("Done")

	// Durations in minutes
	TASK_LENGTH       = 25
	PAUSE_LENGTH      = 5
	LONG_PAUSE_LENGTH = 15
)

type ActivityType string
type ActivityStatus string

type Activity struct {
	Type   ActivityType
	Desc   string
	Status ActivityStatus
}

type TaskBox struct {
	Activity
	ActivityTimer
	ShowInput bool
}

type ActivityTimer struct {
	Minutes string
	Seconds string
}

var (
	task    = &Activity{}
	taskBox = &TaskBox{ShowInput: true, Activity: Activity{Type: TASK}}

	// type -> duration
	durations = map[ActivityType]int{TASK: TASK_LENGTH, PAUSE: PAUSE_LENGTH, LONG_PAUSE: LONG_PAUSE_LENGTH}

	// type -> progress icon
	progressIcons = map[ActivityType]string{TASK: "resources/pomodoro.png", PAUSE: "resources/pause.png", LONG_PAUSE: "resources/pause.png"}

	// type -> done icon
	doneIcons = map[ActivityType]string{TASK: "resources/pomodoro_done.png", PAUSE: "resources/pomodoro.png", LONG_PAUSE: "resources/pomodoro.png"}

	stop = make(chan bool, 1)
)

func (activity *Activity) Render() string {
	return `
    <div class="Activity">
        <h4>
            Enter your next task:
        </h4>
        <input type="text" placeholder="Next thing to do..." onchange="OnInputChange" />
        <p id="pauseButtons">
          <button onclick="OnPause">Short pause</button>
          <button onclick="OnLongPause">Long pause</button>
          <button onclick="OnStopPause">Stop</button>
        </p>
    </div>
    `
}

func (*Activity) OnPause() {
	taskBox.Status = NEW
	startActivity(PAUSE)
}

func (*Activity) OnLongPause() {
	taskBox.Status = NEW
	startActivity(LONG_PAUSE)
}

func (*Activity) OnStopPause() {
	stopActivity()
}

func (taskBox *TaskBox) Render() string {
	return `
    <div class="TaskBox">
      {{if eq .ShowInput true}}<Activity></Activity>
      {{if eq .Status "In Progress"}}<div class="PauseTimer"><span>{{html .ActivityTimer.Minutes}}:{{html .ActivityTimer.Seconds}}</span></div>{{end}}
      {{else}}
      <span oncontextmenu="OnContextMenu"><span id="taskDesc">{{html .Activity.Desc}}</span>:<span id="status">{{html .Status}}</span></span>
      <button onclick="OnStart">Start</button>
      <button onclick="OnStop">Stop</button>
      <button onclick="OnRemove">Remove</button>
      {{if eq .Status "In Progress"}}<div class="ActivityTimer"><span>{{html .ActivityTimer.Minutes}}:{{html .ActivityTimer.Seconds}}</span></div>{{end}}
      {{end}}
    </div>
`
}

func (*TaskBox) OnContextMenu() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&TaskMenu{})
}

func (*TaskBox) OnStart() {
	startActivity(TASK)
}

func (*TaskBox) OnStop() {
	stopActivity()
}

func (*TaskBox) OnRemove() {
	removeTask()
}

func (*Activity) OnInputChange(arg app.ChangeArg) {
	taskBox.ShowInput = false
	taskBox.Desc = arg.Value
	taskBox.Status = NEW

	app.Render(taskBox)
}

func startActivity(activityType ActivityType) {
	if taskBox.Status == NEW {
		taskBox.Type = activityType

		duration := durations[activityType]

		taskBox.Status = IN_PROGRESS
		dock.SetIcon(progressIcons[activityType])

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
						taskBox.ActivityTimer.Minutes = strconv.Itoa(duration)
						taskBox.ActivityTimer.Seconds = "00"
						dock.SetIcon("resources/pomodoro.png")
						break ticker
					}

				default:
					if elapsed > duration*60 {
						taskBox.Status = DONE
						taskBox.ActivityTimer.Minutes = "00"
						taskBox.ActivityTimer.Seconds = "00"
						dock.SetIcon(doneIcons[activityType])
						break ticker
					}
				}

				mins := (duration*60 - elapsed) / 60

				var secs int
				remainder := elapsed % 60
				if remainder == 0 {
					secs = 0
				} else {
					secs = 60 - elapsed%60
				}

				taskBox.ActivityTimer.Minutes = fmt.Sprintf("%02d", mins)
				taskBox.ActivityTimer.Seconds = fmt.Sprintf("%02d", secs)
				app.Render(taskBox)
			}
		}()

	}
}

func stopActivity() {
	if taskBox.Activity.Status == IN_PROGRESS {
		stop <- true
	}
}

func removeTask() {
	if taskBox.Activity.Status == IN_PROGRESS {
		stop <- true
	}
	task.Status = NEW
	task.Desc = ""
	taskBox.ShowInput = true

	app.Render(taskBox)
}
