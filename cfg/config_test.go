package cfg

import (
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestConvertHeaders(t *testing.T) {
	type args struct {
		headers []string
		header  *http.Header
	}
	tests := []struct {
		name   string
		args   args
		verify func(header *http.Header)
	}{
		{
			name: "WhenEmptyArray_ThenDoNothing",
			args: args{
				headers: []string{},
				header:  &http.Header{},
			},
			verify: func(header *http.Header) {
				assert.Empty(t, header)
			},
		},
		{
			name: "WhenInvalidEntry_ThenIgnore",
			args: args{
				headers: []string{"invalid"},
				header:  &http.Header{},
			},
			verify: func(header *http.Header) {
				assert.Empty(t, header)
			},
		},
		{
			name: "WhenValidEntry_ThenParse",
			args: args{
				headers: []string{"Authentication= Bearer <token>"},
				header:  &http.Header{},
			},
			verify: func(header *http.Header) {
				assert.Equal(t, "Bearer <token>", header.Get("Authentication"))
			},
		},
		{
			name: "GivenValidEntry_WhenSpacesAroundValues_ThenTrim",
			args: args{
				headers: []string{"  Authentication =   Bearer <token>  "},
				header:  &http.Header{},
			},
			verify: func(header *http.Header) {
				assert.Equal(t, "Bearer <token>", header.Get("Authentication"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConvertHeaders(tt.args.headers, tt.args.header)
			tt.verify(tt.args.header)
		})
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		envs   map[string]string
		want   *Configuration
		fs     flag.FlagSet
		verify func(c *Configuration)
	}{
		{
			name: "GivenNoFlags_ThenReturnDefaultConfig",
			args: []string{},
			verify: func(c *Configuration) {
				assert.Equal(t, "info", c.Log.Level)
			},
		},
		{
			name: "GivenLogFlags_WhenVerboseEnabled_ThenSetLoggingLevelToDebug",
			args: []string{"-v"},
			verify: func(c *Configuration) {
				assert.Equal(t, "debug", c.Log.Level)
				assert.Equal(t, true, c.Log.Verbose)
			},
		},
		{
			name: "GivenLogFlags_WhenLogLevelSpecified_ThenOverrideLogLevel",
			args: []string{"--log.level=warn"},
			verify: func(c *Configuration) {
				assert.Equal(t, "warn", c.Log.Level)
			},
		},
		{
			name: "GivenLogFlags_WhenInvalidLogLevelSpecified_ThenSetLoggingLevelToInfo",
			args: []string{"--log.level=invalid"},
			verify: func(c *Configuration) {
				assert.Equal(t, "info", c.Log.Level)
			},
		},
		{
			name: "GivenLogLevel_WhenVerboseEnabled_ThenSetLoggingLevelToDebug",
			args: []string{"--log.level=fatal", "-v"},
			verify: func(c *Configuration) {
				assert.Equal(t, "debug", c.Log.Level)
				assert.Equal(t, true, c.Log.Verbose)
			},
		},
		{
			name: "GivenFlags_WhenBindAddrSpecified_ThenOverridePort",
			args: []string{"--bindAddr", ":9090"},
			verify: func(c *Configuration) {
				assert.Equal(t, ":9090", c.BindAddr)
			},
		},
		{
			name: "GivenHeaderFlags_WhenMultipleHeadersSpecified_ThenFillArray",
			args: []string{"--isg.header", "key1=value1", "--isg.header", "KEY2= value2"},
			verify: func(c *Configuration) {
				assert.Contains(t, c.ISG.Headers, "key1=value1")
				assert.Contains(t, c.ISG.Headers, "KEY2= value2")
			},
		},
		{
			name: "GivenHeaderEnvVar_WhenMultipleHeadersSpecified_ThenTreatAsOne",
			envs: map[string]string{
				"ISG_HEADER": "key1=value1, KEY2= value2",
			},
			verify: func(c *Configuration) {
				// Unlike viper, Koanf does not support multiple entries in single ENV var via comma separation.
				// Use the config file or CLI flags to provide those
				assert.Contains(t, c.ISG.Headers, "key1=value1, KEY2= value2")
			},
		},
		{
			name: "GivenUrlFlag_ThenOverrideDefault",
			args: []string{"--isg.url", "myurl"},
			verify: func(c *Configuration) {
				assert.Equal(t, "myurl", c.ISG.URL)
			},
		},
		{
			name: "GivenTimeoutFlag_WhenSpecified_ThenOverrideDefault",
			args: []string{"--isg.timeout", "3"},
			verify: func(c *Configuration) {
				assert.Equal(t, 3*time.Second, c.ISG.Timeout)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.envs)
			result := ParseConfig("version", "commit", "date", &tt.fs, tt.args)
			tt.verify(result)
			unsetEnv(tt.envs)
		})
	}
}

func setEnv(m map[string]string) {
	for key, value := range m {
		os.Setenv(key, value)
	}
}

func unsetEnv(m map[string]string) {
	for key, _ := range m {
		os.Unsetenv(key)
	}
}
