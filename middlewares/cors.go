package middlewares

import "net/http"

// CORSMiddleware bertugas "Membuka Pintu" agar browser (Frontend) 
// diizinkan untuk berkomunikasi dengan API ini dari port yang berbeda.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mengizinkan semua origin (domain). Di Production sungguhan, 
		// ganti "*" dengan nama domain Frontend Anda (misal: "https://my-frontend.com")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Browser biasanya mengirim request "OPTIONS" dulu untuk bertanya "Bolehkah saya masuk?"
		// Kita harus langsung menjawab 200 OK agar browser mau mengirim request aslinya (GET/POST).
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Lanjutkan ke handler tujuan
		next.ServeHTTP(w, r)
	})
}
