package encryption

import (
	"gitlab.pikabiduskibidi.ru/box/go-encryption/internal/interfaces"
	"gitlab.pikabiduskibidi.ru/box/go-encryption/pkg/config"
)

// EncryptorProvider реализует интерфейс interfaces.EncryptorProvider
// для предоставления экземпляров шифровальщиков

type EncryptorProvider struct{}

// NewEncryptorProvider создает новый провайдер шифровальщиков
func NewEncryptorProvider() *EncryptorProvider {
	return &EncryptorProvider{}
}

// ProvideEncryptor предоставляет новый экземпляр шифровальщика
func (p *EncryptorProvider) ProvideEncryptor(cfg *config.Config) (interfaces.Encryptor, error) {
	return NewEncryptor(cfg.Key)
}
