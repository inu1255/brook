package brook

import (
	"sync"
	"time"
)

type traffic struct {
	Threshold int
	onData    func(int, int, int)
	ports     map[int]*portTraffic
}

type portTraffic struct {
	speed int
	up    int
	down  int
	data  int
	prev  time.Time
	lock  sync.Mutex
}

var Traffic = &traffic{
	Threshold: 200 * 1024,
	onData:    func(int, int, int) {},
}

func (t traffic) OnData(fn func(int, int, int)) {
	if fn != nil {
		t.onData = fn
	}
}

func (t traffic) LimitSpeed(port, speed int) {
	p := t.getP(port)
	p.speed = speed
}

func (t traffic) getP(port int) *portTraffic {
	p := t.ports[port]
	if p == nil {
		p = new(portTraffic)
		t.ports[port] = p
	}
	return p
}

func (t traffic) addUp(port, size int) {
	p := t.getP(port)
	p.lock.Lock()
	defer p.lock.Unlock()
	p.up += size
	if p.up+p.down < t.Threshold {
		t.onData(port, p.up, p.down)
		p.up = 0
	}
	if p.speed < 1 {
		return
	}
	now := time.Now()
	sub := now.Sub(p.prev)
	if sub > time.Second {
		p.prev = now
		p.data = 0
		return
	}
	if p.data > p.speed {
		time.Sleep(time.Second - sub)
		p.prev = time.Now()
		p.data = size
	} else {
		p.data += size
	}
}

func (t traffic) addDown(port, size int) {
	p := t.ports[port]
	if p == nil {
		p = new(portTraffic)
		t.ports[port] = p
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.down += int(size)
	if p.up+p.down < t.Threshold {
		t.onData(port, p.up, p.down)
		p.up = 0
	}
	if p.speed < 1 {
		return
	}
	now := time.Now()
	sub := now.Sub(p.prev)
	if sub > time.Second {
		p.prev = now
		p.data = 0
		return
	}
	if p.data > p.speed {
		time.Sleep(time.Second - sub)
		p.prev = time.Now()
		p.data = size
	} else {
		p.data += size
	}
}
