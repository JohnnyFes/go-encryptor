package encryption_test

import (
	"testing"

	"github.com/JohnnyFes/go-encryptor/internal/encryption"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid key 32 bytes",
			key:     "12345678901234567890123456789012",
			wantErr: false,
		},
		{
			name:    "short key",
			key:     "short",
			wantErr: false,
		},
		{
			name:    "long key",
			key:     "this is a very long key that will be hashed to 32 bytes using SHA-256",
			wantErr: false,
		},
		{
			name:    "encrypted key",
			key:     "ENC[AES256:encrypted]",
			wantErr: true,
		},
		{
			name:    "base64 key",
			key:     "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryption.NewEncryptor(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		text    string
		wantErr bool
	}{
		{
			name:    "32 bytes key",
			key:     "12345678901234567890123456789012",
			text:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "short key",
			key:     "short",
			text:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "long key",
			key:     "this is a very long key that will be hashed to 32 bytes using SHA-256",
			text:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "empty text",
			key:     "12345678901234567890123456789012",
			text:    "",
			wantErr: false,
		},
		{
			name:    "special characters",
			key:     "12345678901234567890123456789012",
			text:    "!@#$%^&*()_+",
			wantErr: false,
		},
		{
			name:    "unicode text",
			key:     "12345678901234567890123456789012",
			text:    "Привет, мир!",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := encryption.NewEncryptor(tt.key)
			if err != nil {
				t.Fatalf("Failed to create encryptor: %v", err)
			}

			encrypted, err := encryptor.Encrypt(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !isValidEncryptedFormat(encrypted) {
					t.Errorf("Encrypt() invalid format: %s", encrypted)
					return
				}

				decrypted, err := encryptor.Decrypt(encrypted)
				if err != nil {
					t.Errorf("Decrypt() error = %v", err)
					return
				}

				if decrypted != tt.text {
					t.Errorf("Decrypt() = %v, want %v", decrypted, tt.text)
				}
			}
		})
	}
}

func TestAESEncryptor_DecryptInvalid(t *testing.T) {
	key := "12345678901234567890123456789012"
	encryptor, err := encryption.NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "invalid format",
			text:    "not encrypted",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			text:    "ENC[OTHER:data]",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			text:    "ENC[AES256:invalid base64]",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func isValidEncryptedFormat(text string) bool {
	return len(text) > 11 && text[:11] == "ENC[AES256:" && text[len(text)-1] == ']'
}
