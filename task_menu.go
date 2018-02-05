package main

import (
	"github.com/murlokswarm/app"
)

type TaskMenu struct {
}

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
	startActivity(TASK)
}

func (tm *TaskMenu) OnStop() {
	stopActivity()
}

func (tm *TaskMenu) OnNewTask() {
	removeTask()
}

func init() {
	app.RegisterComponent(&TaskMenu{})
}
