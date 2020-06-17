package main

import (
	"korok.io/korok"
	"korok.io/korok/asset"
	"korok.io/korok/engi"
	"korok.io/korok/game"
	"korok.io/korok/gfx"
	"korok.io/korok/gui"
	"korok.io/korok/hid/input"
	"korok.io/korok/math/f32"
)

type StateEnum int

/*
	三种状态
	Ready 准备状态，显示 Ready 菜单
	Running 运行状态，鸟受候物理作用并隐藏菜单
	Over 挂掉状态，显示失败菜单
*/
const (
	Ready StateEnum = iota
	Running
	Over
)

// 定义重力加速度和点击冲量
const (
	Gravity    = 600
	TapImpulse = 280
)

type GameScene struct {
	state StateEnum //添加游戏状态属性，默认是Ready状态
	ready struct {
		gfx.Tex2D
		gui.Rect
	}
	tap struct {
		gfx.Tex2D
		gui.Rect
	}
	bird struct {
		engi.Entity
		f32.Vec2         //位置
		vy       float32 //y方向速度
		w, h     float32
	}
	bg, ground engi.Entity
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
	sn.bird.Vec2 = f32.Vec2{80, 240}
	korok.Transform.Comp(sn.bird.Entity).SetPosition(sn.bird.Vec2)
}
func (sn *GameScene) Update(dt float32) {
	/*
		根据游戏状态属性的不同来管理生命周期
	*/
	if st := sn.state; st == Ready {
		sn.showReady(dt)
		return
	} else if st == Over {
		sn.showOver(dt)
		return
	}

	//Running状态

	// 检测屏幕点击，每次点击给鸟施加一次冲量
	if input.PointerButton(0).JustPressed() {
		sn.bird.vy = TapImpulse //冲量等于质量*速度 , 这里由于没有设定质量，所以可以直接将冲量换算成速度
	}
	// 模拟物理加速
	sn.bird.vy -= Gravity * dt         //模拟重力加速度对速度的影响
	sn.bird.Vec2[1] += sn.bird.vy * dt //根据速度改变鸟的位置

	// update bird position
	b := korok.Transform.Comp(sn.bird.Entity)
	b.SetPosition(sn.bird.Vec2)

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
	sn.bird.Entity, sn.bg, sn.ground = bird, bg, ground
}

func (sn *GameScene) showReady(dt float32) {
	// show ready
	gui.Image(1, sn.ready.Rect, sn.ready.Tex2D, nil)

	// show tap hint
	gui.Image(2, sn.tap.Rect, sn.tap.Tex2D, nil)

	/*
		通过用户输入系统来获取点击事件，
		检测屏幕点击事件
		改变游戏状态属性
	*/
	if input.PointerButton(0).JustPressed() {
		sn.state = Running
	}
}
func (sn *GameScene) showOver(dt float32) {
	//gameover
}
