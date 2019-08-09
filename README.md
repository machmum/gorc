### gorc - golang reusable component

Reusable component for golang. Provides:

* Logger
    * Log using [zap](https://github.com/uber-go/zap)
    * Logs to file based on current day `yyyy-mm-dd.log`
    * Include `trace-id` for tracking

* Request
    * `RequestID` is a string of the form "host.example.com/random-0001"
    
* String
    * Concatenate string with `Join` using `strings.Builder` 

