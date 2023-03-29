package parallel_test

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"github.com/gozelle/async/parallel"
	"github.com/gozelle/testify/require"
	"sync"
	"testing"
)

func TestRun1(t *testing.T) {
	
	var runners []async.Runner
	
	for i := 1; i <= 10; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			result = v
			return
		})
	}
	
	values := parallel.Run(context.Background(), 2, runners)
	n := 0
	for v := range values {
		n += v.(int)
	}
	require.Equal(t, 55, n)
}

func TestRunError(t *testing.T) {
	
	var runners []async.Runner
	
	for i := 1; i <= 5; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			if v == 3 {
				err = fmt.Errorf("some error")
				return
			}
			result = v
			return
		})
	}
	
	values := parallel.Run(context.Background(), 2, runners)
	var err error
	for v := range values {
		switch vv := v.(type) {
		case error:
			err = vv
		}
	}
	require.Error(t, err)
}

func TestRun2(t *testing.T) {
	var runners []parallel.Runner
	
	for i := 0; i < 100000; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			result = v
			return
		})
	}
	results := parallel.Run(context.Background(), 100, runners)
	n := 0
	for v := range results {
		n++
		switch r := v.(type) {
		case error:
			t.Error(r)
		case int:
			//t.Logf("结果: %d", r)
		default:
			t.Logf("未知类型: %v", r)
		}
	}
	require.Equal(t, n, 100000)
}

func TestChan(t *testing.T) {
	res := calc(t)
	
	for v := range res {
		t.Log(v)
	}
}

func calc(t *testing.T) <-chan int {
	results := make(chan int, 10)
	
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			results <- i
		}(i)
	}
	go func() {
		wg.Wait()
		t.Log("准备关闭")
		close(results)
	}()
	
	return results
}

func TestCloseBroadcast(t *testing.T) {
	wg := sync.WaitGroup{}
	c := make(chan int)
	
	go func() {
		c <- 1
		close(c)
	}()
	
	worker := func(index int) {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			for {
				select {
				case v, ok := <-c:
					t.Logf("worker: %d, closed: %v", v, !ok)
					if !ok {
						return
					}
				}
			}
		}()
	}
	
	worker(1)
	worker(2)
	worker(3)
	wg.Wait()
}

func TestForRange(t *testing.T) {
	c := make(chan int)
	go func() {
		close(c)
	}()
	for v := range c {
		t.Log(v)
	}
	t.Log("即将退出")
}

func TestChCloseSignal(t *testing.T) {
	c := make(chan int)
	go func() {
		close(c)
	}()
	t.Log(<-c)
	t.Log("即将退出")
}

func TestParallel(t *testing.T) {
	_ = run()
}

func run() (err error) {
	
	// 生成 runner
	runner := func(index int) parallel.Runner {
		return func(ctx context.Context) (result any, err error) {
			return index, nil
		}
	}
	
	var runners []parallel.Runner
	for i := 0; i < 10000; i++ {
		runners = append(runners, runner(i))
	}
	
	// 同时最多有 10 个并发
	results := parallel.Run(context.Background(), 10, runners)
	
	// 固定写法，用于从通道中接收处理结果
	for v := range results {
		switch r := v.(type) {
		case int:
			// 结果处理
		case error:
			err = r
			return
		default:
			err = fmt.Errorf("unknown result type: %v", v)
			return
		}
	}
	
	return
}
