package cache

import (
	"time"
)

var cc = make(map[string]*entry)

//传参函数
type Func func(key string) (interface{}, error)

//返回数据
type result struct {
	value interface{}
	err   error
}

//计算解析
type entry struct {
	res     result
	outtime int64
	ready   chan struct{}
}

//中间数据
type request struct {
	key     string
	respone chan<- result
}

type Memo struct {
	respones chan request
}

func New(f Func) *Memo {
	memo := &Memo{
		respones: make(chan request),
	}

	go memo.server(f)
	go memo.out()
	return memo
}

func (m *Memo) Get(key string) (interface{}, error) {
	res := make(chan result)
	m.respones <- request{
		key:     key,
		respone: res,
	}
	cc := <-res
	return cc.value, cc.err
}

func (m *Memo) server(f Func) {
	for res := range m.respones {
		e := cc[res.key]
		if e == nil {
			e = &entry{
				outtime: time.Now().Unix(),
				ready:   make(chan struct{}),
			}
			cc[res.key] = e
			go e.call(f, res.key)
		}
		go e.deliver(res.respone)
	}
}

func (e *Memo) out() {
	tick := time.NewTicker(10 * time.Minute)
	out := 5 * time.Minute
	for {
		select {
		case <-tick.C:
			for k, v := range cc {
				if time.Since(time.Unix(v.outtime, 0)) >= out {
					delete(cc, k)
				}
			}
		}
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}
func (e *entry) deliver(respone chan<- result) {
	<-e.ready
	respone <- e.res
}
