package logger

import (
	sailor "github.com/Lyrics-you/sail-logrus-formatter/sailor"
	"github.com/sirupsen/logrus"
)

// var log = Logger()

func Logger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&sailor.Formatter{
		HideKeys:        true,
		CharStampFormat: "yy/MM/dd HH:mm:ss.SSS",
		Position:        true,
		Colors:          true,
		FieldsColors:    true,
		ShowFullLevel:   true,
	})
	return log
}
