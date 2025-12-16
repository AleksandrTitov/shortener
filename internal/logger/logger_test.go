package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogLevel_Initialize(t *testing.T) {
	tests := []struct {
		name        string
		initLevel   string
		expectLevel logrus.Level
		err         bool
	}{
		{
			name:        "Ошибочный уровень логирования",
			initLevel:   "INVALID",
			expectLevel: logrus.InfoLevel,
			err:         true,
		},
		{
			name:        "Уровень логирования info",
			initLevel:   "INFO",
			expectLevel: logrus.InfoLevel,
			err:         false,
		},
		{
			name:        "Уровень логирования warn",
			initLevel:   "WARN",
			expectLevel: logrus.WarnLevel,
			err:         false,
		},
		{
			name:        "Уровень логирования error",
			initLevel:   "ERROR",
			expectLevel: logrus.ErrorLevel,
			err:         false,
		},
		{
			name:        "Уровень логирования fatal",
			initLevel:   "FATAL",
			expectLevel: logrus.FatalLevel,
			err:         false,
		},
		{
			name:        "Уровень логирования debug",
			initLevel:   "DEBUG",
			expectLevel: logrus.DebugLevel,
			err:         false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Initialize(test.initLevel, false)
			switch {
			case test.err:
				assert.Error(t, err)
			case !test.err:
				assert.NoError(t, err)
			}
			assert.Equal(t, Log.Level, test.expectLevel)
		})
	}
}
