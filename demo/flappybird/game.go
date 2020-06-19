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

type BirdStateEnum int

/*
	鸟的状态
	Flying 鸟的正常飞行状态
	Dead 死亡状态
*/
const (
	Flying BirdStateEnum = iota
	Dead
)

// 定义重力加速度和点击冲量
const (
	Gravity    = 660 //重力加速度  //800
	TapImpulse = 260 //冲量    //300
)

//鸟旋转的常量
const (
	RotTrigger    = 150       //触发旋转需要的速度  //200
	MaxAngle      = 3.14 / 6  //最大仰角
	MinAngle      = -3.14 / 2 //最小俯角
	AngleVelocity = 3.14 * 4  //旋转速度   3.14 * 5
)

//地面滚动速度
const ScrollVelocity = 150 //200

type GameScene struct {
	state StateEnum //添加游戏状态属性，默认是Ready状态
	ready struct {
		gfx.Tex2D
		gui.Rect
	}
	gameover struct {
		gfx.Tex2D
		gui.Rect
	}
	score struct {
		gfx.Tex2D
		gui.Rect
	}
	restart struct {
		gfx.Tex2D
		gui.Rect
	}
	tap struct {
		gfx.Tex2D
		gui.Rect
	}
	bird struct { //鸟
		state BirdStateEnum //鸟状态
		engi.Entity
		f32.Vec2         //位置
		vy       float32 //y方向速度
		w, h     float32
		rotate   float32 //旋转角度
	}
	ground struct { //地面
		engi.Entity
		f32.Vec2
		vx float32
	}
	bg engi.Entity
	PipeSystem
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

	//gameover score restart 的初始化
	sn.gameover.Tex2D, _ = at.GetByName("gameover.png")
	sn.gameover.Rect = gui.Rect{
		X: (320 - 233) / 2,
		Y: 70,
		W: 233,
		H: 70,
	}
	sn.score.Tex2D, _ = at.GetByName("result_board.png")
	sn.score.Rect = gui.Rect{
		X: (320 - 240) / 2,
		Y: 200,
		W: 240,
		H: 120,
	}
	sn.restart.Tex2D, _ = at.GetByName("start.png")
	sn.restart.Rect = gui.Rect{
		X: (320 - 120) / 2,
		Y: 360,
		W: 120,
		H: 60,
	}

	// 重新调整鸟的位置
	sn.bird.Vec2 = f32.Vec2{100, 240}
	korok.Transform.Comp(sn.bird.Entity).SetPosition(sn.bird.Vec2)

	//地面参数初始化
	sn.ground.Vec2 = f32.Vec2{0, 100}
	sn.ground.vx = ScrollVelocity

	//初始化管道系统
	top, _ := at.GetByName("top_pipe.png")
	bottom, _ := at.GetByName("bottom_pipe.png")

	ps := &sn.PipeSystem
	ps.initialize(top, bottom, 6)
	ps.setDelay(0)        // 0 seconds
	ps.setRate(1.2)       // 生成管道间隔的秒数   //0.9
	ps.setGap(120)        // 上下管道间隙
	ps.setLimit(300, 150) //管道生成高度限制
	ps.StartScroll()

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

	// 旋转的仿真代码
	if sn.bird.vy > RotTrigger && sn.bird.rotate < MaxAngle {
		sn.bird.rotate += AngleVelocity * dt
	} else if sn.bird.vy < -RotTrigger && sn.bird.rotate > MinAngle {
		sn.bird.rotate += -AngleVelocity * dt
	}

	// update bird position
	b := korok.Transform.Comp(sn.bird.Entity)
	b.SetPosition(sn.bird.Vec2)
	b.SetRotation(sn.bird.rotate)

	/*
		不断的把“地面”向左移动，
		这样就会产生地面在“移动”的效果，
		当地面的右端被完全进入屏幕的时候，
		我们再重新开始这个动画，
		这样动画可以衔接起来产生一个地面在移动的效果。
	*/
	//地面滚动
	x := sn.ground.Vec2[0]
	if x < -100 {
		x = x + 90 // magic number (bridge start and end of the image)
	}
	x -= sn.ground.vx * dt
	sn.ground.Vec2[0] = x

	// update ground shift
	g := korok.Transform.Comp(sn.ground.Entity)
	g.SetPosition(sn.ground.Vec2)

	//管道系统更新
	sn.PipeSystem.Update(dt)

	// 天空与地面的碰撞检测
	if y := sn.bird.Vec2[1]; y > 480 {
		sn.bird.Vec2[1] = 480
	} else if y < 100 {
		y = 100
		sn.state = Over

		if sn.bird.state != Dead {
			sn.bird.state = Dead
			korok.Flipbook.Comp(sn.bird.Entity).Stop()
		}
	}
	//管道的碰撞检测
	ps := &sn.PipeSystem
	if c, _ := ps.CheckCollision(sn.bird.Vec2, f32.Vec2{sn.bird.w, sn.bird.h}); c {
		if sn.bird.state != Dead {
			ps.StopScroll()
			sn.bird.state = Dead
			korok.Flipbook.Comp(sn.bird.Entity).Stop() // 小鸟动画停止
			sn.state = Over
		}
	}

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
	sn.bird.Entity, sn.bg, sn.ground.Entity = bird, bg, ground
}

func (sn *GameScene) showReady(dt float32) {
	// 重新调整鸟的位置
	sn.bird.Vec2 = f32.Vec2{100, 240}
	korok.Transform.Comp(sn.bird.Entity).SetPosition(sn.bird.Vec2)
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
	// show game over
	gui.Image(1, sn.gameover.Rect, sn.gameover.Tex2D, nil)

	// show score
	gui.Image(2, sn.score.Rect, sn.score.Tex2D, nil)

	// show restart button
	e := gui.ImageButton(3, sn.restart.Rect, sn.restart.Tex2D, sn.restart.Tex2D, nil)
	if e.JustPressed() {
		// do something...
		sn.reStart()
	}
}

//游戏restart事件
func (sn *GameScene) reStart() {
	sn.state = Ready

	// bird
	sn.bird.state = Flying
	sn.bird.Vec2 = f32.Vec2{80, 240}
	sn.bird.vy = 0
	sn.bird.rotate = 0
	korok.Transform.Comp(sn.bird.Entity).SetRotation(0)
	korok.Flipbook.Comp(sn.bird.Entity).Play("flying")
	// pipes
	sn.PipeSystem.Reset()
	sn.PipeSystem.StartScroll()
}
