package grtm

type Pool struct {
	PoolNum  int
	Queue    chan func()
	LimtChan chan bool
}

func NewPool(Num int) *Pool {
	pool := &Pool{}
	pool.PoolNum = Num
	pool.LimtChan = make(chan bool, Num)
	pool.Queue = make(chan func())

	go func() {
		for {

			select {
			case task := <-pool.Queue:
				go func() {
					task()
					defer func() {
						<-pool.LimtChan
					}()
				}()
			}
		}
	}()

	return pool
}

func (p *Pool) AddTask(task func()) {
	p.Queue <- task
}
