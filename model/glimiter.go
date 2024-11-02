package model

// 简易限流
type GLimiter struct {
	size int
	ch   chan struct{}
}

func NewGlimiter(size int) *GLimiter {
	return &GLimiter{
		size: size,
		ch:   make(chan struct{}, size),
	}
}

func (g *GLimiter) Run(fc func()) {
	g.ch <- struct{}{}
	defer func() {
		<-g.ch
	}()
	go fc()
}
