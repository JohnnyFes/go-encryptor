# go-encryption

Пакет для шифрования чувствительных данных в Go приложениях.

## Особенности

- Шифрование AES-256 в режиме GCM
- Поддержка автоматического расшифрования конфигурации
- Утилиты для работы с чувствительными данными
- Поддержка base64 кодирования ключей

## Установка

```bash
go get gitlab.pikabiduskibidi.ru/box/go-encryption@v1.0.0
```

## Работа с приватным репозиторием (важно)

Если вы используете приватный GitLab, для корректной работы go get/go mod рекомендуется установить переменную окружения GOPRIVATE:

```sh
export GOPRIVATE=gitlab.pikabiduskibidi.ru
```

- Это отключает проверку через публичные прокси и checksum-сервисы для всех модулей с этим доменом.
- Публичные пакеты (github.com, golang.org и т.д.) продолжают работать как обычно.
- Можно добавить в .bashrc/.zshrc или в CI/CD pipeline:

```yaml
# Пример для GitLab CI
default:
  before_script:
    - export GOPRIVATE=gitlab.pikabiduskibidi.ru
```

## Использование

### 1. Базовое шифрование строк

```go
import (
    "watchdog/go-encryption/pkg/config"
    "watchdog/go-encryption/pkg/encryption"
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

В проекте есть отдельная CLI-утилита для шифрования строк и работы с конфигурациями. Точка входа: `cmd/encryption/main.go`.

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

# Обновление нескольких полей в конфиге (пример, требует доработки)
$ ./encryption -key="your-32-byte-key" -config="config.yml" -fields="database.password,redis.password" -passwords="secret123,password456"
```

#### Поддерживаемые параметры:
- `-key` — ключ шифрования (обязателен)
- `-passwords` — список паролей для шифрования (через запятую)
- `-config` — путь к YAML/JSON конфигу
- `-fields` — список полей для обновления в конфиге (через запятую)

> Для расширения функциональности до работы с файлами конфигурации доработайте соответствующий блок в `main.go`.

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
