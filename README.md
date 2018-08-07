# env_logger

This is a super simple project which aims to help out with setting up logging correctly in your project

It exposes a single function `ConfigureLogger` which takes in a prefix (the name of your application)
as well as the logger that you wish to configure.

Currently it only supports `logrus`, but PR's are welcome to support additional loggers.



# TODO

- add an interface so that any logger can be injected as the canonical logger

``` go
type Logger struct {
  New() -> Logger // used to instantiate a new logger
  Debug() // emit debug message
  Info()
  Warn()
  Panic()
}
```
