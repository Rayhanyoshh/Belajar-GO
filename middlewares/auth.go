package middlewares

import (
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Rahasia ini HARUS SAMA PERSIS dengan rahasia yang ada di Microservice SSO!
// Jika berbeda, Tracker akan menganggap token dari SSO itu palsu.
var JWTSecretKey = []byte("KUNCI_RAHASIA_SUPER_KUAT_123")

// AuthMiddleware adalah "Satpam" yang menjaga rute API kita
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak: Anda belum Login (Header Authorization kosong)"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak: Format token salah. Gunakan format 'Bearer <token>'"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return JWTSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak: Token JWT Palsu atau sudah kedaluwarsa!"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
