package util

import (
	"sync"
)

type Pool struct {
	lock          sync.Mutex
	goChan        chan func()
	coreNum       int
	maxNum        int
	activeNum     int
	jobNum        int
	exceptionFunc func(r any)
}

func CreatePool(coreNum int, maxNum int) *Pool {
	P := &Pool{
		lock:    sync.Mutex{},
		goChan:  make(chan func(), 5*maxNum),
		coreNum: coreNum,
		maxNum:  maxNum,
	}
	for i := 0; i < coreNum; i++ {
		go P.work()
	}
	return P
}

func (P *Pool) CheckStatus() (activeNum int, jobNum int) {
	P.lock.Lock()
	defer P.lock.Unlock()
	return P.activeNum, P.jobNum
}

func (P *Pool) CreateWork(f func() (E error), exceptionFunc func(Message error)) {
	F := func() {
		if err := f(); err != nil {
			exceptionFunc(err)
			return
		}
	}
	P.goChan <- F
	P.lock.Lock()
	P.jobNum++
	if P.activeNum < P.maxNum && P.jobNum > P.activeNum {
		go P.work()
	}
	P.lock.Unlock()
}

func (P *Pool) work() {
	defer func() {
		r := recover()
		if r != nil {
			P.exceptionFunc(r)
		}
	}()
	P.lock.Lock()
	P.activeNum++
	P.lock.Unlock()
	for {
		f := <-P.goChan
		f()
		P.lock.Lock()
		P.jobNum--
		if P.activeNum > P.coreNum {
			P.activeNum--
			P.lock.Unlock()
			break
		}
		P.lock.Unlock()
	}
}

func (P *Pool) SetExceptionFunc(f func(r any)) {
	P.exceptionFunc = f
}
