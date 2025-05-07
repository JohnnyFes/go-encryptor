package interfaces

import "errors"

var (
	// ErrInvalidKey ошибка при неверном ключе
	ErrInvalidKey = errors.New("invalid encryption key")
	// ErrInvalidData ошибка при неверных данных
	ErrInvalidData = errors.New("invalid data")
	// ErrEncryptionFailed ошибка при шифровании
	ErrEncryptionFailed = errors.New("encryption failed")
	// ErrDecryptionFailed ошибка при расшифровании
	ErrDecryptionFailed = errors.New("decryption failed")
	// ErrInvalidConfig ошибка при неверной конфигурации
	ErrInvalidConfig = errors.New("invalid configuration")
)
