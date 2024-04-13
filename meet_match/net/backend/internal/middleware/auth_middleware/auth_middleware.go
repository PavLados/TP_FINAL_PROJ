package auth_middleware

import (
	"net/http"
	"strings"
	"test_backend_frontend/internal/lib/api/response"
	auth_service "test_backend_frontend/internal/services/auth"
	"test_backend_frontend/pkg/auth_utils"

	"github.com/go-chi/render"
)

func JwtAuthMiddleware(next http.Handler, secret string, tokenHandler auth_utils.ITokenHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		err := tokenHandler.ValidateToken(token, auth_service.SECRET)
		if err != nil {
			if err == auth_utils.ErrParsingToken {
				render.JSON(w, r, response.Error(err.Error()))
				render.Status(r, http.StatusBadRequest)
			} else {
				render.JSON(w, r, response.Error(err.Error()))
				render.Status(r, http.StatusUnauthorized)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}
