# group

errgroup：这是 Go 语言标准库golang.org/x/sync中的一个包，可以在并发操作中简化错误处理。如果你的 goroutines 中的任何一个返回错误，errgroup 可以确保所有的 goroutines 都会被取消。

```go
var g errgroup.Group

urls := []string{"http://url1", "http://url2", "http://url3"}
for _, url := range urls {
    // capture range variable
    url := url
    g.Go(func() error {
        // Fetch the URL.
        resp, err := http.Get(url)
        if err == nil {
            resp.Body.Close()
        }
        return err
    })
}
// Wait for all HTTP fetches to complete.
if err := g.Wait(); err == nil {
    fmt.Println("Successfully fetched all URLs.")
}
```