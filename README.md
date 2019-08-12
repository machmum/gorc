### gorc - golang reusable component

Reusable component for golang. Provides:

* Logger ([docs](https://godoc.org/github.com/machmum/gorc/logger))
    * Log using [zap](https://github.com/uber-go/zap)
    * use `NewSugaredLogger` for fast and more verbose logger
    * use `NewLogger` for fast, leveled, and structured logging (strongly-typed logger)
    * Logs to file based on current day `yyyy-mm-dd.log`
    * Include `trace-id` for tracking

* Request ([docs](https://godoc.org/github.com/machmum/gorc/request))
    * `RequestID` is a string of the form "host.example.com/random-0001"
    
* String ([docs](https://godoc.org/github.com/machmum/gorc/stringc))
    * Concatenate string with `StringBuilder` using `strings.Builder` 

