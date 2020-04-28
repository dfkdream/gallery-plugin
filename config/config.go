package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

type Config struct {
	BoltPath      string `json:"boltPath"`
	Interpolation resize.InterpolationFunction
	Quality       int `json:"quality"`
}

func Get() *Config {
	return &Config{
		BoltPath:      getEnvStringOr("BOLT", "./gallery.db"),
		Interpolation: getEnvInterpolationOr("INTERPOLATION", resize.Lanczos3),
		Quality:       getEnvIntOr("QUALITY", 80),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("BoltPath: %s\nInterpolation: %d\nQuality: %d", c.BoltPath, c.Interpolation, c.Quality)
}

func getEnvStringOr(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInterpolationOr(key string, defaultValue resize.InterpolationFunction) resize.InterpolationFunction {
	switch key {
	case "NearestNeighbor":
		return resize.NearestNeighbor
	case "Bilinear":
		return resize.Bilinear
	case "Bicubic":
		return resize.Bicubic
	case "MitchellNetravali":
		return resize.MitchellNetravali
	case "Lanczos2":
		return resize.Lanczos2
	case "Lanczos3":
		return resize.Lanczos3
	default:
		return defaultValue
	}
}

func getEnvIntOr(key string, defaultValue int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}
