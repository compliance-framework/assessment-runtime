package plugins

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Register(plugin Plugin) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}
