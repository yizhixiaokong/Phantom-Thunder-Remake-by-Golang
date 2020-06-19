package main

import (
	"korok.io/korok"
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
	"korok.io/korok/math"
	"korok.io/korok/math/f32"
)

//管道系统
type Pipe struct {
	/*
		Pipe 表示一个管道，
		它包含了上下两根，
		还有相应的位置，移动速度，高度等参数。
	*/
	top struct {
		engi.Entity
		f32.Vec2
	}
	bottom struct {
		engi.Entity
		f32.Vec2
	}
	high   float32
	active bool
	x, vx  float32
}

//AABB碰撞检测
type AABB struct {
	x, y          float32
	width, height float32
}

func OverlapAB(a, b *AABB) bool {
	if a.x < b.x+b.width && a.x+a.width > b.x && a.y < b.y+b.height && a.y+a.height > b.y {
		return true
	}
	return false
}

//初始化
/*
	初始化管道各种参数
*/
func (p *Pipe) initialize(texTop, texBottom gfx.Tex2D) {
	top := korok.Entity.New()
	spr := korok.Sprite.NewComp(top)
	spr.SetSprite(texTop)
	spr.SetSize(65, 400)
	spr.SetGravity(.5, 0) //设置锚点

	bottom := korok.Entity.New()
	spr = korok.Sprite.NewComp(bottom)
	spr.SetSize(65, 400)
	spr.SetSprite(texBottom)
	spr.SetGravity(.5, 1) //设置锚点

	// out of screen
	korok.Transform.NewComp(top).SetPosition(f32.Vec2{-100, 210})
	korok.Transform.NewComp(bottom).SetPosition(f32.Vec2{-100, 160})

	p.top.Entity = top
	p.bottom.Entity = bottom
	p.vx = ScrollVelocity
}

//管道根据传入参数重置位置属性
func (p *Pipe) reset(x, high, gap float32) {
	p.active = true
	p.x = x
	p.top.Vec2 = f32.Vec2{x, high + gap}
	p.bottom.Vec2 = f32.Vec2{x, high}
}

//管道更新
func (p *Pipe) update(dt float32) {
	p.x -= p.vx * dt
	if p.x < -50 {
		p.active = false
	}

	p.top.Vec2[0] = p.x
	p.bottom.Vec2[0] = p.x

	korok.Transform.Comp(p.top.Entity).SetPosition(p.top.Vec2)
	korok.Transform.Comp(p.bottom.Entity).SetPosition(p.bottom.Vec2)
}

//管道系统
type PipeSystem struct {
	gap, top, bottom float32 // gap, top, bottom limit
	respawn          float32 // respawn location
	scroll           bool

	delay struct { //延迟参数
		clock float32
		limit float32
	}
	generate struct { //生成参数
		clock float32
		limit float32
	}

	pipes []*Pipe //运行管道
	frees []*Pipe //待用管道

	_pool []Pipe
}

//管道系统初始化
func (ps *PipeSystem) initialize(texTop, texBottom gfx.Tex2D, size int) {
	ps._pool = make([]Pipe, size)
	ps.frees = make([]*Pipe, size) // add to freelist
	for i := range ps._pool {
		ps.frees[i] = &ps._pool[i]
		ps.frees[i].initialize(texTop, texBottom)
	}
	ps.respawn = 320 + 20
}

//设置管道出现的延迟时间
func (ps *PipeSystem) setDelay(d float32) {
	ps.delay.limit = d
}

//设置管道生成速率
func (ps *PipeSystem) setRate(r float32) {
	ps.generate.limit = r
}

//设置上下管道间的间隙
func (ps *PipeSystem) setGap(gap float32) {
	ps.gap = gap
}

//设置随机的高度范围从b到top
func (ps *PipeSystem) setLimit(top, b float32) {
	ps.top, ps.bottom = top, b
}

//管道系统更新
func (ps *PipeSystem) Update(dt float32) {
	if !ps.scroll {
		return
	}

	// delay some time
	//延迟多久更新一次
	if d := &ps.delay; d.clock < d.limit {
		d.clock += dt
		return
	}

	// generate new pipe
	//控制速率生成新的管道
	if g := &ps.generate; g.clock < g.limit {
		g.clock += dt
	} else {
		g.clock = 0
		ps.newPipe()
	}

	// update pipe
	//管道更新
	for _, p := range ps.pipes {
		p.update(dt)
	}

	// recycle
	//管道循环
	ps.recycle()
}

//管道系统停止滚动
func (ps *PipeSystem) StopScroll() {
	ps.scroll = false
}

//管道系统开始滚动
func (ps *PipeSystem) StartScroll() {
	ps.scroll = true
}

//管道系统重置
func (ps *PipeSystem) Reset() {
	for _, p := range ps.pipes {
		p.x = -100
		// out of screen
		korok.Transform.NewComp(p.top.Entity).SetPosition(f32.Vec2{-100, 210})
		korok.Transform.NewComp(p.bottom.Entity).SetPosition(f32.Vec2{-100, 160})
	}
}

//生成新管道
func (ps *PipeSystem) newPipe() {
	if sz := len(ps.frees); sz > 0 {
		p := ps.frees[sz-1]
		ps.frees = ps.frees[:sz-1]
		ps.pipes = append(ps.pipes, p)
		p.reset(ps.respawn, math.Random(ps.bottom, ps.top), ps.gap)
	}
}

// inactive pipes come first
//管道循环，滚动到最左边的管道移动到待用管道数组里
func (ps *PipeSystem) recycle() {
	pipes, inactive := ps.pipes, -1
	for i, p := range pipes {
		if p.active {
			break
		}
		inactive = i
	}
	if inactive >= 0 {
		ps.pipes = pipes[inactive+1:]
		ps.frees = append(ps.frees, pipes[:inactive+1]...)
	}
}

//管道系统碰撞检测
func (ps *PipeSystem) CheckCollision(p f32.Vec2, sz f32.Vec2) (bool, float32) {
	//p鸟位置 sz鸟大小
	tolerance := float32(8) //允许的误差(判定范围)
	sz[0], sz[1] = sz[0]-tolerance, sz[1]-tolerance
	bird := &AABB{p[0] - sz[0]/2, p[1] - sz[1]/2, sz[0], sz[1]}
	for _, p := range ps.pipes {
		top := &AABB{
			p.top.Vec2[0] - 32,
			p.top.Vec2[1],
			65,
			400,
		}
		if OverlapAB(bird, top) {
			return true, bird.x - top.x
		}

		bottom := &AABB{
			p.bottom.Vec2[0] - 32,
			p.bottom.Vec2[1] - 400,
			65,
			400,
		}
		if OverlapAB(bird, bottom) {
			return true, bird.x - top.x
		}
	}
	return false, 0 //没有碰撞
}
