package main

import (
	"korok.io/korok"
	"korok.io/korok/anim/frame"
	"korok.io/korok/asset"
	"korok.io/korok/engi"
	"korok.io/korok/game"
	"korok.io/korok/gfx"
	"korok.io/korok/gui"
	"korok.io/korok/math/f32"
)

type StartScene struct {
	title struct {
		gfx.Tex2D
		gui.Rect
	}
	start struct {
		btnNormal  gfx.Tex2D
		btnPressed gfx.Tex2D
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

//加载游戏场景的方法
func (sn *StartScene) LoadGame() {
	gsn := &GameScene{}
	//调用GameScene声明的borrow方法，将当前帧的Entity传递给下一帧
	gsn.borrow(sn.bird.Entity, sn.bg, sn.ground)

	// load game scene
	/*
		调用 Load 方法加载场景，
		然后调用 Push 方法把当前场景入栈，
		这个操作会导致前一个场景“退出”(并没有真的退出只是它的 Update 方法不会再被调用)。
		同时新的场景的生命周期方法开始被依次调用，
		比如 OnEnter 和 OnUpdate
	*/
	korok.SceneMan.Load(gsn)
	korok.SceneMan.Push(gsn)
}
func (sn *StartScene) OnEnter(g *game.Game) {
	//  get textue
	/*
		加载了之前的纹理图集，并从中找出了 "bird1.png" 这张图片。
	*/
	at, _ := asset.Texture.Atlas(`F:\Phantom Thunder (Remake by Golang)\Phantom-Thunder-Remake-by-Golang\demo\flappybird\assets\images\bird.png`)
	//因为运行路径的原因，这里直接用绝对路径避免出错 //就是长了点
	/*
		找出了background.png,game_name.png,ground.png
	*/
	bg, _ := at.GetByName("background.png")
	ground, _ := at.GetByName("ground.png")
	tt, _ := at.GetByName("game_name.png")
	btn, _ := at.GetByName("start.png")
	/*
		接下来是分别将背景、地面、标题、GUI等场景加载进游戏
	*/
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
	// setup gui
	/*
		Korok 中使用的 ImmediateMode GUI(IMGUI 即时图形GUI)
	*/
	// title
	sn.title.Tex2D = tt
	sn.title.Rect = gui.Rect{ //UI边界大小
		X: (320 - 233) / 2,
		Y: 80,
		W: 233,
		H: 70,
	}
	// start button
	sn.start.btnNormal = btn
	sn.start.btnPressed = btn
	sn.start.Rect = gui.Rect{ //UI边界大小
		X: (320 - 120) / 2,
		Y: 300,
		W: 120,
		H: 60,
	}
	// flying animation
	/*
		获取小鸟煽动翅膀的三张图片
		将三张图片组成一个动画，以名字flying为索引存储在动画池中
	*/
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
	/*
		为Entity添加Flipbook组件
		使其循环播放flying动画
		形成小鸟飞行的效果
	*/
	anim := korok.Flipbook.NewComp(bird)
	anim.SetRate(.1)
	// fmt.Println(anim.Loop())
	anim.SetLoop(true, frame.Restart) //设置循环播放开启，播放模式为restart
	//按理说这里不用自己设置循环的，不知道为啥引擎没设置
	anim.Play("flying")

	sn.bird.Entity = bird

}
func (sn *StartScene) Update(dt float32) {
	// draw something
	/*
		直接调用 gui.Image 就可以绘制一张图片，
		使用 gui.ImageButton 就可以绘制一个按钮，
		并且可以通过返回值来得到按钮的状态
	*/
	/*
		当你调用这个方法的时候这个 UI 就已经在绘制了；
		同样当你不再调用这个方法的时候，
		UI就不再绘制没有任何缓存的状态存在。
	*/
	/*
		因为 Update 方法不断被调用的缘故，
		实际上是在不断重新绘制UI，
		所以它能显示出来
	*/

	// draw title
	gui.Image(1, sn.title.Rect, sn.title.Tex2D, nil)
	// draw start button
	e := gui.ImageButton(2, sn.start.Rect, sn.start.btnNormal, sn.start.btnPressed, nil)
	if e.JustPressed() { //按键按下
		// do something
		//游戏开始
		sn.LoadGame()
	}
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
