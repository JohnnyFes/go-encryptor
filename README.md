# go-encryption

Пакет для шифрования чувствительных данных в Go приложениях.

## Особенности

- Шифрование AES-256 в режиме GCM
- Поддержка автоматического шифрования и расшифрования полей структур по тегам `encrypted`
- Утилиты для работы с чувствительными данными
- Поддержка base64 кодирования ключей

## Установка

```bash
go get github.com/GandzyTM/go-encryptor@latest
```

## Работа с приватным репозиторием (важно)

Если репозиторий находится в приватной организации GitHub, настройте переменную окружения `GOPRIVATE`, чтобы `go` не обращался к публичным прокси и checksum-сервисам:

```sh
export GOPRIVATE=github.com/JohnnyFes
```

- Это нужно выполнить локально и в CI/CD пайплайнах, если модуль недоступен публично.
- Остальные публичные модули (github.com, golang.org и т.д.) продолжат загружаться как обычно.
- Переменную можно добавить в shell-профиль или в конфигурацию пайплайна, например:

```yaml
# Пример для GitHub Actions
env:
  GOPRIVATE: github.com/JohnnyFes
```

## Использование

### 1. Базовое шифрование строк

```go
import (
    "github.com/JohnnyFes/go-encryptor/pkg/config"
    "github.com/JohnnyFes/go-encryptor/pkg/encryption"
)

// Создаем конфигурацию
cfg, err := config.NewConfig("your-32-byte-encryption-key")
if err != nil {
    log.Fatal(err)
}

// Создаем шифратор
encryptor, err := encryption.NewEncryptor(cfg)
if err != nil {
    log.Fatal(err)
}

// Шифруем строку
encrypted, err := encryptor.EncryptString("sensitive-data")
if err != nil {
    log.Fatal(err)
}

// Расшифровываем строку
decrypted, err := encryptor.DecryptString(encrypted)
if err != nil {
    log.Fatal(err)
}
```

### 2. Шифрование полей структуры

```go
type UserConfig struct {
    Username string `encrypted:"false"`
    Password string `encrypted:"true"`
    APIKey   string `encrypted:"true"`
    Email    string `encrypted:"true"`
}

// Создаем структуру
config := UserConfig{
    Username: "john_doe",
    Password: "secret123",
    APIKey:   "api-key-123",
    Email:    "john@example.com",
}

// Шифруем чувствительные поля
err = encryptor.EncryptFields(&config)

// Расшифровываем чувствительные поля
err = encryptor.DecryptFields(&config)
```

## Безопасность

- Используйте ключ длиной минимум 32 байта
- Храните ключ шифрования в безопасном месте (например, в vault)
- Не храните ключ в коде или конфигурационных файлах
- Регулярно ротируйте ключи шифрования

## CLI-утилита

В проекте есть отдельная CLI-утилита для шифрования строк и обновления конфигурационных файлов. Точка входа: `cmd/encryption/main.go`.

### Сборка

```bash
# Собрать бинарник
$ go build -o encryption ./cmd/encryption
```

### Установка

```bash
# Установить в $GOBIN (или $HOME/go/bin)
$ go install ./cmd/encryption
```

### Использование

```bash
# Шифрование нескольких паролей
$ ./encryption -key="your-32-byte-key" -passwords="secret123,password456,key789"

# Обновление нескольких полей в конфиге
$ ./encryption -key="your-32-byte-key" -config="config.yml" -fields="database.password,redis.password" -passwords="secret123,password456"
```

#### Поддерживаемые параметры:
- `-key` — ключ шифрования (обязателен)
- `-passwords` — список паролей для шифрования (через запятую)
- `-config` — путь к YAML/JSON конфигу
- `-fields` — список полей для обновления в конфиге (через запятую)

> Утилита шифрует значения и обновляет YAML/JSON файл на месте. Расшифровка через CLI не поддерживается — используйте пакет `pkg/encryption` в приложении.

## Примеры CLI-команд

### Шифрование одной строки (пароля)

```bash
./encryption -key="your-32-byte-key" -passwords="mysecret"
```

### Шифрование нескольких строк

```bash
./encryption -key="your-32-byte-key" -passwords="secret1,secret2,secret3"
```

### Шифрование и обновление полей в YAML/JSON конфиге

```bash
./encryption -key="your-32-byte-key" -config="config.yml" -fields="database.password,redis.password" -passwords="dbpass,redispass"
```
- Количество полей и паролей должно совпадать.
- Поддерживаются вложенные поля через точку (например, `database.password`).

### Расшифровка строки (если поддерживается)

> На данный момент CLI поддерживает только шифрование. Для расшифровки используйте библиотеку в коде Go:

```go
// Пример
plaintext, err := encryptor.DecryptString(encrypted)
```

---

> Для расширения CLI (например, добавления расшифровки через флаг) — доработайте main.go.
