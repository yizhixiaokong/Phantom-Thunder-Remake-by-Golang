package main

import (
	"korok.io/korok"
	"korok.io/korok/game"
	"korok.io/korok/gfx/dbg"
)

type MainScene struct {
}

func (m *MainScene) OnEnter(g *game.Game) {
}

func (m *MainScene) Update(dt float32) {
	dbg.DrawStr(180, 160, "Hello World")
}

func (*MainScene) OnExit() {
}

func main() {
	options := &korok.Options{
		Title:  "Hello, Korok Engine",
		Width:  480,
		Height: 320,
	}
	korok.Run(options, &MainScene{})
}
