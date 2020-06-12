package main

import (
	"fmt"

	"korok.io/korok"
	"korok.io/korok/anim/frame"
	"korok.io/korok/asset"
	"korok.io/korok/engi"
	"korok.io/korok/game"
	"korok.io/korok/gfx"
	"korok.io/korok/math/f32"
)

type StartScene struct {
	bird, bg, ground engi.Entity
}

func (sn *StartScene) Load() {
	//asset.Texture.LoadAtlas加载纹理图集
	/*
		纹理图集通常是把多张图片打包在一张大的图片中，
		并通过一个描述文件来描述每个子图片的位置。
		使用纹理图集可以显著的减少 DrawCall 提高游戏的性能。
	*/
	asset.Texture.LoadAtlas(`F:\Phantom Thunder (Remake by Golang)\Phantom-Thunder-Remake-by-Golang\demo\flappybird\assets\images\bird.png`,
		`F:\Phantom Thunder (Remake by Golang)\Phantom-Thunder-Remake-by-Golang\demo\flappybird\assets\images\bird.json`)
}

func (sn *StartScene) OnEnter(g *game.Game) {
	//  get textue
	/*
		加载了之前的纹理图集，并从中找出了 "bird1.png" 这张图片。
	*/
	at, _ := asset.Texture.Atlas(`F:\Phantom Thunder (Remake by Golang)\Phantom-Thunder-Remake-by-Golang\demo\flappybird\assets\images\bird.png`)

	/*
		找出了background.png,game_name.png,ground.png
	*/
	bg, _ := at.GetByName("background.png")
	// tt, _ := at.GetByName("game_name.png")
	ground, _ := at.GetByName("ground.png")
	// setup bg
	{
		entity := korok.Entity.New()
		spr := korok.Sprite.NewCompX(entity, bg)
		spr.SetSize(320, 480)
		xf := korok.Transform.NewComp(entity)
		xf.SetPosition(f32.Vec2{160, 240})
		sn.bg = entity
	}

	// setup ground {840 281}
	{
		entity := korok.Entity.New()
		spr := korok.Sprite.NewCompX(entity, ground)
		spr.SetSize(420, 140)
		spr.SetGravity(0, 1)
		spr.SetZOrder(1)
		xf := korok.Transform.NewComp(entity)
		xf.SetPosition(f32.Vec2{0, 100})
		sn.ground = entity
	}

	// flying animation
	bird1, _ := at.GetByName("bird1.png")
	bird2, _ := at.GetByName("bird2.png")
	bird3, _ := at.GetByName("bird3.png")

	frames := []gfx.Tex2D{bird1, bird2, bird3}
	g.AnimationSystem.SpriteEngine.NewAnimation("flying", frames, true)

	// setup bird
	/*
		新建了一个 Entity，
		并给它添加了 SpriteComp 和 Transform 组件，
		Transform 组件赋予了这个 Entity 在游戏世界中的位置，
		SpriteComp 让这个 Entity 可以绘制出来。
	*/

	bird := korok.Entity.New()
	spr := korok.Sprite.NewCompX(bird, bird1)
	spr.SetSize(48, 32)
	spr.SetZOrder(2)
	xf := korok.Transform.NewComp(bird)
	xf.SetPosition(f32.Vec2{160, 240})

	// play animation
	anim := korok.Flipbook.NewComp(bird)
	anim.SetRate(.1)
	fmt.Println(anim.Loop())
	anim.SetLoop(true, frame.Restart)
	anim.Play("flying")

	sn.bird = bird

}
func (sn *StartScene) Update(dt float32) {
	// draw something
}
func (sn *StartScene) OnExit() {
}

func main() {
	options := korok.Options{
		Title:  "Flappy Bird",
		Width:  320,
		Height: 480,
	}
	korok.Run(&options, &StartScene{})
}
