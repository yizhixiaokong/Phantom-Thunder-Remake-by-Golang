package main

import (
	"korok.io/korok"
	"korok.io/korok/asset"
	"korok.io/korok/engi"
	"korok.io/korok/game"
	"korok.io/korok/gfx"
	"korok.io/korok/gui"
	"korok.io/korok/math/f32"
)

type GameScene struct {
	ready struct {
		gfx.Tex2D
		gui.Rect
	}
	tap struct {
		gfx.Tex2D
		gui.Rect
	}
	bird, bg, ground engi.Entity
}

/*
	在 Korok 中是不主动删除前一个场景初始化的 Entity 的，
	所以即使场景转换了，
	在前一个场景中初始化的 Entity 还是会被绘制的
*/
/*
	场景切换后，原本的菜单和标题小时
	因为GUI 系统的工作原理，
	当前一帧的Update()方法不再被调用的时候，
	它的UI也就不再绘制了，
	而新场景的Update()方法开始被调用，
	开始绘制新场景的UI

*/
func (sn *GameScene) OnEnter(g *game.Game) {
	at, _ := asset.Texture.Atlas(`F:\Phantom Thunder (Remake by Golang)\Phantom-Thunder-Remake-by-Golang\demo\flappybird\assets\images\bird.png`)

	// ready and tap image
	sn.ready.Tex2D, _ = at.GetByName("getready.png")
	sn.ready.Rect = gui.Rect{
		X: (320 - 233) / 2,
		Y: 70,
		W: 233,
		H: 70,
	}
	sn.tap.Tex2D, _ = at.GetByName("tap.png")
	sn.tap.Rect = gui.Rect{
		X: (320 - 143) / 2,
		Y: 200,
		W: 143, // 286
		H: 123, // 246
	}
	// 重新调整鸟的位置
	korok.Transform.Comp(sn.bird).SetPosition(f32.Vec2{80, 240})
}
func (sn *GameScene) Update(dt float32) {
	// show ready
	gui.Image(1, sn.ready.Rect, sn.ready.Tex2D, nil)

	// show tap hint
	gui.Image(2, sn.tap.Rect, sn.tap.Tex2D, nil)
}
func (sn *GameScene) OnExit() {
	/*
		如果不希望保留前一场景的Entity可以在OnExit回调中删除，
		删除一个 Entity 只要删除它拥有的组件即可
	*/
}

// 为当前帧声明borrow方法，用于“借用”前一帧Entity
/*
	在前一帧(上一个场景)调用这里声明的borrow方法，
	将之前的Entity传递给当前帧(当前场景)
*/
func (sn *GameScene) borrow(bird, bg, ground engi.Entity) {
	sn.bird, sn.bg, sn.ground = bird, bg, ground
}
