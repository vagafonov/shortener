package middleware

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/vagafonov/shortener/internal/cookie"
	"github.com/vagafonov/shortener/pkg/encrypting"
)

// Проверяет наличие cookie c идентификатором пользователя и выдает ее в случае ее отсутствия.
func (mw *middleware) WithUserIDCookie(next http.Handler, cryptoKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDCoockie, err := r.Cookie("userID")
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				mw.logger.Error().Msgf("error: %v\n", err)

				return
			}
		}

		if userIDCoockie != nil && userIDCoockie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		if userIDCoockie == nil {
			mw.logger.Debug().Msg("cookie doesn't exist. setting")
			ck := mw.setCookie(w, cryptoKey)
			r.AddCookie(ck)
		} else {
			mw.logger.Debug().Msg("cookie exist")
			// Расшифровка UUID
			binUserID, _ := hex.DecodeString(userIDCoockie.Value)
			decrypted, err := encrypting.Decrypt(binUserID, cryptoKey)
			mw.logger.Info().Str("decrypted", decrypted).Send()
			if err != nil {
				mw.setCookie(w, cryptoKey)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (mw *middleware) setCookie(w http.ResponseWriter, cryptoKey []byte) *http.Cookie {
	c := cookie.CreateCookieWithUserID(mw.logger, cryptoKey)
	mw.logger.Debug().Str("name", c.Name).Str("value", c.Value).Msg("created cookie")
	http.SetCookie(w, c)

	return c
}
