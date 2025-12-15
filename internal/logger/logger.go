package logger

import "github.com/sirupsen/logrus"

var Log = logrus.New()

func Initialize(logLevel string, jsonFormatter bool) error {
	switch jsonFormatter {
	case true:
		Log.SetFormatter(&logrus.JSONFormatter{})
	case false:
		Log.SetFormatter(&logrus.TextFormatter{})
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Log.SetLevel(logrus.InfoLevel)
		return err
	}
	Log.SetLevel(level)

	return nil
}
