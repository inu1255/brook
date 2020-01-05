package brook

import (
	"sync"
	"time"
)

type traffic struct {
	ThresHold int
	Speed     int
	ondata    func(int, int, int)
	ports     map[int]*portTraffic
}

type portTraffic struct {
	up   int64
	down int64
	data int64
	prev time.Time
	lock sync.Mutex
}

var Traffic = &traffic{
	ThresHold: 200 * 1024,
	Speed:     2 * 1024 * 1024,
	ondata:    func(int, int, int) {},
}

func (t traffic) OnData(fn func(int, int, int)) {
	if fn != nil {
		t.ondata = fn
	}
}

func (t traffic) report(port int) {
	if p.up+p.down < hold {
		t.ondata(port, p.up, p.down)
		p.up = 0
	}
	if t.Speed < 1 {
		return
	}
	now := time.Now()
	sub := now.Sub(p.prev)
	if sub > time.Second {
		p.prev = now
		p.data = 0
		return
	}
	if p.data > t.Speed {
		time.Sleep(time.Second - sub)
		p.prev = time.Now()
		p.data = size
	} else {
		p.data += size
	}
}

func (t traffic) addUp(port, size int) {
	p := t.ports[port]
	if p == nil {
		p = new(portTraffic)
		t.ports[port] = p
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.up += size
	t.report(port)
}

func (t traffic) addDown(port, size int) {
	p := t.ports[port]
	if p == nil {
		p = new(portTraffic)
		t.ports[port] = p
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.down += size
	t.report(port)
}
