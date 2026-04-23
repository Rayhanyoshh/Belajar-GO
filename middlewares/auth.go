package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Rahasia ini HARUS SAMA PERSIS dengan rahasia yang ada di Microservice SSO!
// Jika berbeda, Tracker akan menganggap token dari SSO itu palsu.
var JWTSecretKey = []byte("KUNCI_RAHASIA_SUPER_KUAT_123")

// AuthMiddleware adalah "Satpam" yang menjaga rute API kita
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Cek apakah ada karcis masuk di Header bernama "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Akses Ditolak: Anda belum Login (Header Authorization kosong)"}`, http.StatusUnauthorized)
			return
		}

		// 2. Format standar JWT adalah "Bearer <kode_token_panjang_sekali>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Artinya kata "Bearer " tidak ditemukan
			http.Error(w, `{"error": "Akses Ditolak: Format token salah. Gunakan format 'Bearer <token>'" }`, http.StatusUnauthorized)
			return
		}

		// 3. Validasi & Bongkar Isi Token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma enkripsinya menggunakan HMAC standar
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			// Berikan kunci gemboknya untuk membongkar token
			return JWTSecretKey, nil
		})

		// 4. Jika gagal dibongkar (Palsu atau sudah kedaluwarsa)
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Akses Ditolak: Token JWT Palsu atau sudah kedaluwarsa!"}`, http.StatusUnauthorized)
			return
		}

		// 5. Ekstrak data dari dalam token (Misalnya kita ambil user_id)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Memasukkan data user_id ke dalam "Tas Bawaan" (Context)
			// Agar fungsi utama (misal: AddWebsite) tahu "Siapa" yang memanggilnya
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			r = r.WithContext(ctx)
		}

		// 6. LOLOS! Silakan masuk ke dalam ruangan utama
		next(w, r)
	}
}
