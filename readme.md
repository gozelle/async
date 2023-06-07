# Async 

## 安装

```
go get -u github.com/gozelle/async
```


## parallel 并发调用

通过信号量控制同一时间的最大并发数量，如果中间有错误产生，则会取消后续排队的任务。

```go
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
			return
		}
		// 处理数据
		_ = v.Value
	}
	
	return
}
```

## race 竞赛调用

取最先成功的结果返回，取消或忽略后续排队的任务执行，如果全部任务都失败，则返回合并的错误。

```go
func main() {
	runners := []*race.Runner[int]{
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 1
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 3
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 2
				return
			},
		},
	}
	
	race.Run[int](context.Background(), runners)
}   
```
