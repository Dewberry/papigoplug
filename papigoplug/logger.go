package papigoplug

import (
	"log" // used when logger fails to initiate

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLog starts the logrus logger. This should be called exactly one time by the main program before the program uses
// other functions from papigoplug. A good default level is "info".
func InitLog(level string) {
	if Log != nil {
		Log.Fatal("log has already been initialized")
	}
	levelCoded, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatalf("unknown log level: %q. choices: %q", level, logrus.AllLevels)
	}
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(levelCoded)
}
