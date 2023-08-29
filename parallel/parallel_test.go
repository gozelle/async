package parallel_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	
	"github.com/gozelle/async/parallel"
	"github.com/gozelle/atomic"
	"github.com/gozelle/testify/require"
)

func TestRun1(t *testing.T) {
	
	var runners []parallel.Runner[int]
	
	for i := 1; i <= 10; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result int, err error) {
			result = v
			time.Sleep(time.Duration(v) * time.Second)
			return
		})
	}
	
	values := parallel.Run[int](context.Background(), 2, runners)
	t.Log("begin")
	n := 0
	err := parallel.Wait[int](values, func(v int) error {
		n += v
		t.Log(v)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 55, n)
}

func TestRunError(t *testing.T) {
	
	var runners []parallel.Runner[int]
	sum := atomic.NewInt32(0)
	for i := 1; i <= 100; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result int, err error) {
			if v == 20 {
				err = fmt.Errorf("some error")
				fmt.Println("error")
				return
			} else {
				sum.Add(int32(v))
			}
			result = v
			fmt.Println(v)
			return
		})
	}
	
	values := parallel.Run[int](context.Background(), 5, runners)
	total := int32(0)
	err := parallel.Wait[int](values, func(v int) error {
		total += int32(v)
		return nil
	})
	require.Error(t, err)
	require.Equal(t, sum.Load(), total)
}

func TestOversize(t *testing.T) {
	
	var runners []parallel.Runner[int]
	sum := atomic.NewInt32(0)
	for i := 1; i <= 100; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result int, err error) {
			if v == 20 {
				err = fmt.Errorf("some error")
				fmt.Println("error")
				return
			} else {
			}
			sum.Add(int32(v))
			result = v
			fmt.Println(v)
			return
		})
	}
	
	values := parallel.Run[int](context.Background(), 200, runners)
	total := int32(0)
	err := parallel.Wait[int](values, func(v int) error {
		total += int32(v)
		return nil
	})
	require.Error(t, err)
	require.Equal(t, sum.Load(), total)
}

func TestErrorWait(t *testing.T) {
	var runners []parallel.Runner[int]
	
	for i := 1; i <= 10; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result int, err error) {
			fmt.Println("运行了", v)
			time.Sleep(1 * time.Second)
			if v == 1 {
				err = fmt.Errorf("some error")
			}
			if v == 2 {
				time.Sleep(5 * time.Second)
			}
			
			result = v
			
			return
		})
	}
	
	values := parallel.Run[int](context.Background(), 3, runners)
	err := parallel.Wait[int](values, func(v int) error {
		t.Log(v)
		return nil
	})
	t.Log(err)
	require.Error(t, err)
}

func TestRun2(t *testing.T) {
	var runners []parallel.Runner[int]
	
	for i := 0; i < 100000; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result int, err error) {
			result = v
			return
		})
	}
	results := parallel.Run[int](context.Background(), 100, runners)
	n := 0
	for v := range results {
		n++
		require.NoError(t, v.Error)
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
	runner := func(index int) parallel.Runner[int] {
		return func(ctx context.Context) (result int, err error) {
			return index, nil
		}
	}
	
	var runners []parallel.Runner[int]
	for i := 0; i < 10000; i++ {
		runners = append(runners, runner(i))
	}
	
	// 同时最多有 10 个并发
	results := parallel.Run[int](context.Background(), 10, runners)
	
	// 固定写法，用于从通道中接收处理结果
	for v := range results {
		if v.Error != nil {
			// 错误处理
		} else {
			// 处理数据
			_ = v.Value
		}
	}
	
	return
}

func TestRunWithCancel(t *testing.T) {
	
	var runners []parallel.Runner[int]
	
	for i := 1; i <= 5; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result int, err error) {
			time.Sleep(time.Duration(v) * time.Second)
			t.Log(v)
			return
		})
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancel()
	}()
	
	values := parallel.Run[int](ctx, 2, runners)
	
	err := parallel.Wait[int](values, func(v int) error {
		return nil
	})
	require.Error(t, err)
	
}
