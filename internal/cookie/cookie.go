package cookie

import (
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/pkg/encrypting"
)

const maxAgeCookie = 3600

// CreateCookieWithUserID create cookie.
func CreateCookieWithUserID(l *zerolog.Logger, cryptoKey []byte) *http.Cookie {
	// Генерация UUID
	newUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	// Преобразование UUID в строку
	uuidString := newUUID.String()
	l.Debug().Msg("Сгенерированный UUID:" + uuidString)

	// Шифрование UUID
	encrypted, err := encrypting.Encrypt(uuidString, cryptoKey)
	if err != nil {
		panic(err)
	}
	l.Debug().Msg("Зашифрованный UUID:" + hex.EncodeToString(encrypted))

	return &http.Cookie{
		Name:     "userID",
		Value:    hex.EncodeToString(encrypted),
		Path:     "/",
		MaxAge:   maxAgeCookie,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	}
}

// Decrypt.
func Decrypt(enc string, cryptoKey []byte) (*string, error) {
	binid, err := hex.DecodeString(enc)
	if err != nil {
		return nil, err
	}
	decrypted, err := encrypting.Decrypt(binid, cryptoKey)
	if err != nil {
		return nil, err
	}

	return &decrypted, nil
}
