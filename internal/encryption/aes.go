package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// AESEncryptor реализует шифрование данных с использованием AES-256.
// Структура содержит блок шифра AES, который используется для:
// - Шифрования чувствительных данных в конфигурации
// - Расшифровки данных при загрузке конфигурации
// - Обеспечения безопасности паролей и других конфиденциальных данных
type AESEncryptor struct {
	block cipher.Block
}

// NewEncryptor создает новый экземпляр AES шифровальщика
func NewEncryptor(key string) (*AESEncryptor, error) {
	// Если ключ не в base64, пробуем использовать его как есть
	var keyBytes []byte
	if strings.HasPrefix(key, "ENC[") {
		return nil, fmt.Errorf("encryption key cannot be encrypted")
	}

	// Пробуем декодировать как base64
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		// Если не получилось, используем ключ как есть
		keyBytes = []byte(key)
	}

	// Если ключ длиннее 32 байт, используем SHA-256 для получения ключа нужной длины
	if len(keyBytes) > 32 {
		hash := sha256.Sum256(keyBytes)
		keyBytes = hash[:]
	}

	// Если ключ короче 32 байт, дополняем его нулями
	if len(keyBytes) < 32 {
		newKey := make([]byte, 32)
		copy(newKey, keyBytes)
		keyBytes = newKey
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return &AESEncryptor{
		block: block,
	}, nil
}

// Encrypt шифрует данные
func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	// Создаем GCM
	aesGCM, err := cipher.NewGCM(e.block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Создаем nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Шифруем данные
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Кодируем в base64 и добавляем префикс
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return fmt.Sprintf("ENC[AES256:%s]", encoded), nil
}

// Decrypt расшифровывает данные
func (e *AESEncryptor) Decrypt(encrypted string) (string, error) {
	// Проверяем префикс
	if !strings.HasPrefix(encrypted, "ENC[AES256:") || !strings.HasSuffix(encrypted, "]") {
		return "", fmt.Errorf("invalid encrypted data format")
	}

	// Извлекаем зашифрованные данные
	encrypted = encrypted[11 : len(encrypted)-1]

	// Декодируем base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Создаем GCM
	aesGCM, err := cipher.NewGCM(e.block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Извлекаем nonce
	nonce := ciphertext[:aesGCM.NonceSize()]
	ciphertext = ciphertext[aesGCM.NonceSize():]

	// Расшифровываем данные
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
