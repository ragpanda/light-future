package light_future

type Pool interface {
	Exec(func())
}

type GoroutineInfanitePool struct{}

func (*GoroutineInfanitePool) Exec(exec func()) {
	go exec()
}
