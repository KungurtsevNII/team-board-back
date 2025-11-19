# Prettylog

A golang `log/slog` handler for pretty console output, designed to be used during development time only.

This software is "borrowed" (copy-pasted) from this article by [dusted.codes](https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package) with improvements from this repository by [dusted-go](https://github.com/dusted-go/logging).

This is mainly for me to have a single source to easily retrieve this code in my own projects. If anyone is seriously interested in this software I highly recommend visiting dusted-go's repo instead.

An example (for me): 

```go
prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
    Level:       slog.LevelInfo,
    AddSource:   false,
    ReplaceAttr: nil,
})
logger := slog.New(prettyHandler)
```
