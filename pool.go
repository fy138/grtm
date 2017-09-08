package grtm

type Pool struct {
	PoolNum  int
	Queue    chan Command
	LimtChan chan bool
}
type Command struct {
	Func interface{}
	Args []interface{}
}

func NewPool(Num int) *Pool {
	pool := &Pool{}
	pool.PoolNum = Num
	pool.LimtChan = make(chan bool, Num)
	pool.Queue = make(chan Command)

	go func() {
		for {

			select {
			case cmd := <-pool.Queue:
				go func(cmd Command) {
					if len(cmd.Args) > 1 {
						cmd.Func.(func(...interface{}))(cmd.Args)
					} else if len(cmd.Args) == 1 {
						cmd.Func.(func(interface{}))(cmd.Args[0])
					} else {
						cmd.Func.(func())()
					}
					defer func() {
						<-pool.LimtChan
					}()
				}(cmd)
			}
		}
	}()

	return pool
}

func (p *Pool) AddTask(fc interface{}, args ...interface{}) {
	var cmd Command
	cmd.Func = fc
	cmd.Args = args
	p.Queue <- cmd
}
