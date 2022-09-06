package util

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Pool is a goroutine pool
type Pool struct {
	lock          sync.RWMutex
	goChan        chan func()
	coreNum       int
	maxNum        int
	activeNum     int
	jobNum        int
	timeout       int
	exceptionFunc func(r any)
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

//
// CreatePool
// @Description: 创建一个协程池
// @param        coreNum 核心协程数
// @param        maxNum  最大协程数
// @param        timeout 超时时间
// @return       *Pool   协程池
//
func CreatePool(coreNum int, maxNum int, timeout int) *Pool {
	ctx, cancelFunc := context.WithCancel(context.Background())
	P := &Pool{
		lock:       sync.RWMutex{},
		goChan:     make(chan func(), 5*maxNum),
		coreNum:    coreNum,
		maxNum:     maxNum,
		timeout:    timeout,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	for i := 0; i < coreNum; i++ {
		go P.work()
	}
	P.lock.Lock()
	P.activeNum = coreNum
	P.lock.Unlock()
	return P
}

//
// CheckStatus
// @Description: 检查协程池状态
// @receiver     P         协程池
// @return       activeNum 活跃协程数
// @return       jobNum    任务数
//
func (P *Pool) CheckStatus() (activeNum int, jobNum int) {
	P.lock.RLock()
	defer P.lock.RUnlock()
	return P.activeNum, P.jobNum
}

//
// CreateWork
// @Description: 创建一个任务
// @receiver     P             协程池
// @param        f             任务函数
// @param        exceptionFunc 异常处理函数
//
func (P *Pool) CreateWork(f func() (E error), exceptionFunc func(Message error)) {
	F := func() {
		if err := f(); err != nil {
			exceptionFunc(err)
			return
		}
	}
	select {
	case P.goChan <- F:
		P.lock.Lock()
		P.jobNum++
		P.lock.Unlock()
	case <-time.After(time.Duration(P.timeout) * time.Second):
		P.exceptionFunc(errors.New("goroutine超时"))
		return
	}
	P.lock.Lock()
	if P.activeNum < P.maxNum && P.jobNum > P.activeNum {
		P.activeNum++
		go P.work()
	}
	P.lock.Unlock()
}

//
// work
// @Description: 协程池工作函数
// @receiver     P 协程池
//
func (P *Pool) work() {
	defer func() {
		r := recover()
		if r != nil {
			P.exceptionFunc(r)
		}
	}()
	for {
		select {
		case <-P.ctx.Done():
			P.lock.Lock()
			P.activeNum--
			P.lock.Unlock()
			return
		case f := <-P.goChan:
			f()
			P.lock.Lock()
			P.jobNum--
			if P.activeNum > P.coreNum && P.jobNum < P.activeNum {
				P.activeNum--
				P.lock.Unlock()
				return
			}
			P.lock.Unlock()
		case <-time.After(time.Duration(P.timeout) * time.Second):
		}
	}
}

//
// SetExceptionFunc
// @Description: 设置异常处理函数
// @receiver     P 协程池
// @param        f 异常处理函数
//
func (P *Pool) SetExceptionFunc(f func(r any)) {
	P.exceptionFunc = f
}

//
// Close
// @Description: 关闭协程池
// @receiver     P 协程池
//
func (P *Pool) Close() {
	P.lock.Lock()
	P.maxNum = 0
	P.coreNum = 0
	P.cancelFunc()
	P.lock.Unlock()
	time.Sleep(time.Second * time.Duration(P.timeout))
}
