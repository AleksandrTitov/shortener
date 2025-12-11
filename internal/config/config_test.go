package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFlags_NewConfig(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		addr     string
		baseHTTP string
	}{
		{
			name:     "Значения по умолчанию",
			addr:     defaultAddr,
			baseHTTP: defaultBaseHTTP,
		},
		{
			name:     "Определяем адрес через '-a'",
			addr:     ":8081",
			flags:    []string{"-a", ":8081"},
			baseHTTP: defaultBaseHTTP,
		},
		{
			name:     "Определяем адрес через '--a'",
			addr:     "0.0.0.0:8081",
			flags:    []string{"--a", "0.0.0.0:8081"},
			baseHTTP: defaultBaseHTTP,
		},
		{
			name:     "Определяем базовый URL через '-b'",
			addr:     defaultAddr,
			flags:    []string{"-b", "http://tt.go"},
			baseHTTP: "http://tt.go",
		},
		{
			name:     "Определяем базовый URL '--b'",
			addr:     defaultAddr,
			flags:    []string{"--b", "https://tt.go"},
			baseHTTP: "https://tt.go",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []string{
				"shorter",
			}
			args = append(args, test.flags...)
			os.Args = args
			flag.CommandLine = flag.NewFlagSet(args[0], flag.ExitOnError)
			cfg := NewConfig()
			assert.Equal(t, cfg.Addr, test.addr)
			assert.Equal(t, cfg.BaseHTTP, test.baseHTTP)
		})
	}
}

func TestEnv_NewConfig(t *testing.T) {
	tests := []struct {
		name     string
		envs     map[string]string
		addr     string
		baseHTTP string
	}{
		{
			name:     "Значения по умолчанию",
			envs:     map[string]string{},
			addr:     defaultAddr,
			baseHTTP: defaultBaseHTTP,
		},
		{
			name: "Определяем порт через 'SERVER_ADDRESS'",

			envs: map[string]string{
				"SERVER_ADDRESS": ":8081",
			},
			addr:     ":8081",
			baseHTTP: defaultBaseHTTP,
		},
		{
			name: "Определяем ip:порт через 'SERVER_ADDRESS'",
			addr: "0.0.0.0:8081",
			envs: map[string]string{
				"SERVER_ADDRESS": "0.0.0.0:8081",
			},
			baseHTTP: defaultBaseHTTP,
		},
		{
			name: "Определяем базовый URL через 'BASE_URL'",
			envs: map[string]string{
				"BASE_URL": "http://tt.go",
			},
			addr:     defaultAddr,
			baseHTTP: "http://tt.go",
		},
		{
			name: "Определяем базовый URL и адрес через 'BASE_URL' и 'SERVER_ADDRESS'",
			envs: map[string]string{
				"BASE_URL":       "https://tt.go",
				"SERVER_ADDRESS": "0.0.0.0:8081",
			},
			addr:     "0.0.0.0:8081",
			baseHTTP: "https://tt.go",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.envs {
				err := os.Setenv(k, v)
				assert.NoError(t, err)
			}
			args := []string{
				"shorter",
			}
			os.Args = args
			flag.CommandLine = flag.NewFlagSet(args[0], flag.ExitOnError)

			cfg := NewConfig()
			for k, _ := range test.envs {
				err := os.Unsetenv(k)
				assert.NoError(t, err)
			}

			assert.Equal(t, cfg.Addr, test.addr)
			assert.Equal(t, cfg.BaseHTTP, test.baseHTTP)
		})
	}
}
