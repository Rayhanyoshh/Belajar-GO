# Stage 1: Pembangun (Builder)
# Kita menggunakan image Go resmi untuk men-compile kode
FROM golang:alpine AS builder

# Mengatur direktori kerja di dalam kontainer
WORKDIR /app

# Menyalin daftar kebutuhan library (go.mod) dan mengunduhnya
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh kode kita ke dalam kontainer
COPY . .

# Melakukan kompilasi kode Go menjadi satu file utuh (binary) bernama "tracker-api"
# CGO_ENABLED=0 membuat binary-nya mandiri tanpa butuh C library (sangat ringan)
RUN CGO_ENABLED=0 GOOS=linux go build -o tracker-api main.go

# ==========================================

# Stage 2: Hasil Akhir (Production)
# Kita menggunakan image Alpine yang sizenya hanya 5 MB!
FROM alpine:latest

WORKDIR /app

# Menambahkan sertifikat SSL agar API bisa memanggil HTTPS dengan aman
RUN apk --no-cache add ca-certificates

# Mengambil HANYA file binary (tracker-api) dari Stage 1
# Kode mentah kita tidak akan dibawa ke server sehingga aman dari pencurian kode!
COPY --from=builder /app/tracker-api .

# Membuka port 8080
EXPOSE 8080

# Perintah yang dijalankan saat server Docker menyala
CMD ["./tracker-api"]
