package main

import "github.com/murlokswarm/app"

const (
	NEW         = "New"
	IN_PROGRESS = "In Progress"
	DONE        = "Done"
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
      <span oncontextmenu="OnContextMenu"><span id="taskDesc">{{html .Task.Desc}}</span>:<span id="status">{{html .Status}}</span></span>{{end}}
      {{if eq .Status "In Progress"}}<div class="TaskTimer"><span>{{html .TaskTimer.Minutes}}:{{html .TaskTimer.Seconds}}</span></div>{{end}}
    </div>
`
}

func (taskBox *TaskBox) OnContextMenu() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&TaskMenu{})
}

func (task *Task) OnInputChange(arg app.ChangeArg) {
	task.Desc = arg.Value
	task.Status = NEW
	taskBox.ShowInput = false
	taskBox.Desc = task.Desc
	taskBox.Status = NEW

	app.Render(taskBox)
}
