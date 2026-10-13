package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Ambil token dari header "Authorization"
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "akses ditolak, token tidak ditemukan",
			})
		}

		// 2. Hilangkan kata "Bearer " jika ada
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Validasi token dengan JWT_SECRET
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "token tidak valid atau sudah kedaluwarsa",
			})
		}

		// 4. Ambil data payload (claims) dari dalam token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "gagal membaca data token",
			})
		}

		// 5. Simpan user_id ke dalam context (untuk dipakai di Handler nanti)
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "format user_id pada token tidak ditemukan atau tidak valid",
			})
		}

		c.Set("user_id", int(userIDFloat))

		// 6. Lanjut ke proses handler berikutnya
		return next(c)
	}
}
