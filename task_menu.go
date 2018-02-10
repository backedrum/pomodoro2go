package main

import (
	"github.com/murlokswarm/app"
)

type TaskMenu struct {
	Disabled DisabledActivities
}

func (tm *TaskMenu) Render() string {
	return `
   <menu label="TaskMenu">
      <menuitem label="Start" shortcut="meta+s" onclick="OnStart" {{if eq (.GetDisabledActivity "Task" "start") true}}disabled="true"{{end}}/>
      <menuitem label="Stop" shortcut="meta+z" onclick="OnStop" {{if eq (.GetDisabledActivity "Task" "stop") true}}disabled="true"{{end}}/>
      <menuitem label="New task" shortcut="meta+n" onclick="OnNewTask"/>
   </menu>
`
}

func (tm *TaskMenu) OnStart() {
	startActivity(TASK)
}

func (tm *TaskMenu) OnStop() {
	stopActivity()
}

func (tm *TaskMenu) OnNewTask() {
	removeTask()
}

func init() {
	app.RegisterComponent(taskMenu)
}
