package mailbox

type Dispatcher interface {
	Schedule(fn func())
	Throughput() int
}

type goroutineDispatcher int

// 缺省的调度器，调度的时候直接开启一个goroutine
func (goroutineDispatcher) Schedule(fn func()) {
	go fn()
}

func (d goroutineDispatcher) Throughput() int {
	return int(d)
}

func NewDefaultDispatcher(throughput int) Dispatcher {
	return goroutineDispatcher(throughput)
}

type synchronizedDispatcher int

// 同步调度器：直接在当前环境下执行函数调用
func (synchronizedDispatcher) Schedule(fn func()) {
	fn()
}

func (d synchronizedDispatcher) Throughput() int {
	return int(d)
}

func NewSynchronizedDispatcher(throughput int) Dispatcher {
	return synchronizedDispatcher(throughput)
}
