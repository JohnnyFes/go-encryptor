package encryption

import (
	"testing"

	"gitlab.pikabiduskibidi.ru/box/go-encryption/pkg/config"
	"gitlab.pikabiduskibidi.ru/box/go-encryption/pkg/encryption"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid key",
			key:     "12345678901234567890123456789012",
			wantErr: false,
		},
		{
			name:    "invalid key length",
			key:     "short",
			wantErr: true,
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
			cfg, err := config.NewConfig(tt.key)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			_, err = encryption.NewEncryptor(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptor_EncryptDecrypt(t *testing.T) {
	key := "12345678901234567890123456789012"
	cfg, err := config.NewConfig(key)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	encryptor, err := encryption.NewEncryptor(cfg)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "simple text",
			text:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "empty text",
			text:    "",
			wantErr: false,
		},
		{
			name:    "special characters",
			text:    "!@#$%^&*()_+",
			wantErr: false,
		},
		{
			name:    "unicode text",
			text:    "Привет, мир!",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := encryptor.EncryptString(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Проверяем формат зашифрованных данных
				if !isValidEncryptedFormat(encrypted) {
					t.Errorf("EncryptString() invalid format: %s", encrypted)
					return
				}

				// Расшифровываем и проверяем результат
				decrypted, err := encryptor.DecryptString(encrypted)
				if err != nil {
					t.Errorf("DecryptString() error = %v", err)
					return
				}

				if decrypted != tt.text {
					t.Errorf("DecryptString() = %v, want %v", decrypted, tt.text)
				}
			}
		})
	}
}

func TestEncryptor_DecryptInvalid(t *testing.T) {
	key := "12345678901234567890123456789012"
	cfg, err := config.NewConfig(key)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	encryptor, err := encryption.NewEncryptor(cfg)
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
			_, err := encryptor.DecryptString(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func isValidEncryptedFormat(text string) bool {
	return len(text) > 11 && text[:11] == "ENC[AES256:" && text[len(text)-1] == ']'
}
