# Async 


## parallel 并发调用

通过信号量控制同一时间的最大并发数量，如果中间有错误产生，则会取消后续排队的任务。

```go
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
```

## race 竞赛调用

特性：取最先成功的结果返回，取消或忽略后续排队的任务执行。

```go
func main(){
	runners := []*race.Runner{
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result any, err error) {
				result = 1
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result any, err error) {
				result = 3
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result any, err error) {
				result = 2
				return
			},
		},
	}
	
	r, err := race.Run(context.Background(), runners)
	if err != nil{
	    return 	
    }
}    
```