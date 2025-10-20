package sensitive

import (
	"reflect"

	"github.com/JohnnyFes/go-encryptor/internal/interfaces"
)

// FieldEncryptor реализует интерфейс interfaces.FieldEncryptor для шифрования полей в структурах
type FieldEncryptor struct {
	encryptor interfaces.Encryptor
}

// NewFieldEncryptor создает новый экземпляр FieldEncryptor
func NewFieldEncryptor(encryptor interfaces.Encryptor) *FieldEncryptor {
	return &FieldEncryptor{
		encryptor: encryptor,
	}
}

// HandleFields обрабатывает поля структуры, шифруя или расшифровывая их
func (h *FieldEncryptor) HandleFields(data interface{}, encrypt bool) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return interfaces.ErrInvalidData
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Проверяем тег encrypted
		if fieldType.Tag.Get("encrypted") == "true" {
			if field.Kind() != reflect.String {
				continue
			}

			value := field.String()
			var result string
			var err error

			if encrypt {
				result, err = h.encryptor.Encrypt(value)
			} else {
				result, err = h.encryptor.Decrypt(value)
			}

			if err != nil {
				return err
			}

			field.SetString(result)
		}
	}

	return nil
}
