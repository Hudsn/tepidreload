package tepidreload

import (
	"io/fs"
	"os"
	"strings"
)

type Config struct {
	TickIntervalMS    int
	WatchPath         fs.FS
	ExcludeDirs       []string
	ExcludeFiles      []string
	ExcludeExtensions []string
}

type ConfigFunc func(*Config)

func defaultConfig() Config {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	myFS := os.DirFS(pwd)

	return Config{
		TickIntervalMS:    250,
		WatchPath:         myFS,
		ExcludeDirs:       []string{},
		ExcludeFiles:      []string{},
		ExcludeExtensions: []string{},
	}
}

func NewConfig(configFuncs ...ConfigFunc) Config {
	config := defaultConfig()
	for _, configFunc := range configFuncs {
		configFunc(&config)
	}

	return config
}

func WithInterval(intervalMS int) ConfigFunc {
	return func(config *Config) {
		config.TickIntervalMS = intervalMS
	}
}

func WithWatchPath(path string) ConfigFunc {
	return func(config *Config) {
		targetFS := os.DirFS(path)
		config.WatchPath = targetFS
	}
}

func WithEmbedFS(fs fs.FS) ConfigFunc {
	return func(config *Config) {
		config.WatchPath = fs
	}
}

func WithExcludeDirs(dirs ...string) ConfigFunc {
	return func(config *Config) {
		config.ExcludeDirs = dirs
	}
}
func WithExcludeFiles(files ...string) ConfigFunc {
	return func(config *Config) {
		config.ExcludeFiles = files
	}
}
func WithExcludeExtensions(extensions ...string) ConfigFunc {
	return func(config *Config) {
		cleaned := []string{}
		for _, ext := range extensions {
			tempExt := strings.TrimPrefix(ext, ".")
			cleaned = append(cleaned, "."+tempExt)
		}
		config.ExcludeExtensions = cleaned
	}
}
