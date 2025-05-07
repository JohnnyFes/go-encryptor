package config

import (
	"encoding/base64"
	"errors"
	"strings"
)

const (
	// DefaultKeyLength стандартная длина ключа шифрования
	DefaultKeyLength = 32
)

var (
	// ErrInvalidKeyLength ошибка при неверной длине ключа
	ErrInvalidKeyLength = errors.New("invalid key length")
)

// Config содержит настройки для шифрования
type Config struct {
	// Key - ключ шифрования
	Key string
	// KeyLength - требуемая длина ключа
	KeyLength int
}

// Option функция для настройки конфигурации
type Option func(*Config)

// WithKeyLength устанавливает требуемую длину ключа
func WithKeyLength(length int) Option {
	return func(c *Config) {
		c.KeyLength = length
	}
}

// NewConfig создает новую конфигурацию
func NewConfig(key string, opts ...Option) (*Config, error) {
	cfg := &Config{
		Key:       key,
		KeyLength: DefaultKeyLength,
	}

	// Применяем опции
	for _, opt := range opts {
		opt(cfg)
	}

	// Проверяем длину ключа
	if len(key) < cfg.KeyLength {
		return nil, ErrInvalidKeyLength
	}

	// Проверяем, что ключ не зашифрован
	if strings.HasPrefix(key, "ENC[") {
		return nil, errors.New("encrypted key is not allowed")
	}

	// Проверяем, что ключ в base64
	if _, err := base64.StdEncoding.DecodeString(key); err == nil {
		return cfg, nil
	}

	return cfg, nil
}
