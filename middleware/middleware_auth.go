package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Alphonnse/file_server/handlers"
	"github.com/Alphonnse/file_server/internal/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ApiConfigWrapper struct {
	handlers.ApiConfig // Embed the original type because of the original is
	// in the other package
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfigWrapper) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find out the token from cookie
		tokenString, err := r.Cookie("Authorization")

		if err != nil {
			fmt.Println("its from midleware", r.URL.Path)
			handlers.RespondWithError(w, 401, fmt.Sprint("Error with cookie:", err))
			return
		}

		// Decode/validate it
		// для нерусских. Здесь мы вызываешь парс, отдаем токенстринг и
		// функция возращающую интерфейс и в функции проверяя сингинг метод
		// отдаем СЕКРЕТ
		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // there might me an error
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil

		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // the invalid token

			// check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				handlers.RespondWithError(w, 401, fmt.Sprint("Error with JWT:", err))
			}

			// Find the user with token sub

			subUUID, err := uuid.Parse(claims["sub"].(string))
			if err != nil {
				handlers.RespondWithError(w, 500, fmt.Sprint("Error while parsing UUID:", err))
				return
			}
			user, err := apiCfg.DB.FindWithID(r.Context(), subUUID)
			if err != nil {
				handlers.RespondWithError(w, 500, fmt.Sprint("No matched UUID in DB:", err))
			}

			// Attach to req and continue
			handler(w, r, user)

		} else {
			handlers.RespondWithError(w, 401, fmt.Sprint("Error with JWT:", err))
		}
	}
}
