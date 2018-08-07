package main

import (
	"runtime"
	"os"
	"strings"
	logger "github.com/Sirupsen/logrus"
)

var (
	defaultLogger *logger.Logger
	loggers = make(map[string]*logger.Logger)
)

const (
	DebugV = iota
	InfoV = iota
	WarnV = iota
)

type Logger interface {
  // New()  Logger // used to instantiate a new logger
	Debug(...interface{}) // emit debug message
	Info(...interface{})
	Warn(...interface{})
	Panic(...interface{})
}

func toEnum(s string) int {
	switch strings.ToLower(s) {
	case "warn":
		return WarnV
	case "debug":
		return DebugV
	case "info":
		return InfoV
	default:
		// TODO: add log output saying defaulting to info
		return InfoV
	}

}

func configurePackageLogger(log *logger.Logger, value int) *logger.Logger {
	switch value {
	case WarnV:
		log.SetLevel(logger.WarnLevel)
	case InfoV:
		log.SetLevel(logger.InfoLevel)
	case DebugV:
		log.SetLevel(logger.DebugLevel)
	default:
		log.SetLevel(logger.InfoLevel)
	}
	return log
}

// ConfigureLogger takes in a prefix and a logger object and configures the logger depending on environment variables. 
// Configured based on the GOLANG_DEBUG environment variable
func ConfigureLogger(newDefaultLogger *logger.Logger)  {
	levels := make(map[string]int)

	// TODO: do parsing and stuff for the env-var
	if debugRaw, ok := os.LookupEnv("GOLANG_LOG"); ok {
		packages := strings.Split(debugRaw, ",")

		for _, pkg := range packages {
			// check if a package name has been specified, if not default to main
			tmp := strings.Split(pkg, "=")
			if len(tmp) == 1 {
				levels["main"] = toEnum(tmp[0])
			} else if len(tmp) == 2 {
				levels[tmp[0]] = toEnum(tmp[1])
			} else {
				newDefaultLogger.Fatal("line: '", pkg, "' is formatted incorrectly, please refer to the documentation for correct usage")
			}

			newDefaultLogger.Warn("pkg:", pkg)
		}
	}

	for key, value := range levels {
		loggers[key] = configurePackageLogger(logger.New(), value)
	}

	// configure main logger
	if value, ok := loggers["main"]; ok {
		defaultLogger = value
	} else {
		defaultLogger = newDefaultLogger
	}
}

// Props to https://stackoverflow.com/a/35213181 for the code
func getPackage () string {

    // we get the callers as uintptrs - but we just need 1
    fpcs := make([]uintptr, 1)

    // skip 4 levels to get to the caller of whoever called getPackage()
    n := runtime.Callers(4, fpcs)
    if n == 0 {
       return "" // proper error her would be better
    }

    // get the info of the actual function that's in the pointer
    fun := runtime.FuncForPC(fpcs[0]-1)
    if fun == nil {
      return ""
    }

    name := fun.Name()
    // return its name
	return strings.Split(name, ".")[0]
}

type F func(Logger)
func printLog(f F) {
	pkg := getPackage()
	if log, ok := loggers[pkg]; ok {
		f(log)
		return
	}
	f(defaultLogger)
}

// Warn prints a warning...
func Warn(args ...interface{})  {
	lambda := func(log Logger) {
		log.Warn(args...)
	}
	printLog(lambda)
}

func Info(args ...interface{})  {
	lambda := func(log Logger) {
		log.Info(args...)
	}
	printLog(lambda)
}

func Debug(args ...interface{})  {
	lambda := func(log Logger) {
		log.Debug(args...)
	}
	printLog(lambda)
}

func main() {
	Warn("--------")
	Warn("warn")
	Info("info")
	Debug("debug")
}

func init() {
	defaultLogger = logger.New()
	ConfigureLogger(defaultLogger)
}
