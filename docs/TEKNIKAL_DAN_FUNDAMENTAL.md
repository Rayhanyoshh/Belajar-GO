# Dokumentasi Teknikal & Fundamental: GoTracker

Dokumen ini dirancang khusus untuk Anda baca sebagai referensi pribadi. Ini berisi rangkuman konsep-konsep fundamental Golang yang telah kita terapkan di dalam proyek pemantau website (**GoTracker**) ini.

---

## 1. Arsitektur Proyek (Standard Go Layout)

Aplikasi ini tidak digabungkan ke dalam satu file `main.go`, melainkan dipecah menjadi beberapa *package* (folder) sesuai dengan standar industri (*Clean Architecture* versi sederhana):

*   **`/models`**: Berisi **cetakan data** (Struct). Struct adalah cara Go merepresentasikan sebuah objek (mirip dengan Class di bahasa lain, tapi Go bukan murni OOP).
*   **`/database`**: Tempat kita mengatur koneksi (*Connection Pool*) ke PostgreSQL menggunakan *driver* `github.com/lib/pq`.
*   **`/handlers`**: Otak dari rute HTTP. Tempat kita meletakkan fungsi-fungsi yang menangani permintaan (*request*) dari *user* (seperti GET dan POST).
*   **`/workers`**: Tempat menyimpan fungsi *Cron Job* (pekerja otomatis) yang berjalan sendiri di latar belakang tanpa menunggu permintaan dari *user*.
*   **`main.go`**: Titik awal (*entry point*) aplikasi. Semakin sedikit kode di sini, semakin bagus.

---

## 2. Fundamental 1: REST API (`net/http`)

Go memiliki pustaka bawaan (*standard library*) bernama `net/http` yang sangat kuat, bahkan sering digunakan tanpa memerlukan *framework* tambahan seperti Express.js (Node.js) atau Laravel (PHP).

**Cara Kerjanya:**
1.  **Router / Mux**: `mux := http.NewServeMux()` berfungsi sebagai "petunjuk jalan". Jika ada permintaan masuk, ia akan mengarahkannya ke fungsi (*handler*) yang tepat.
2.  **Handler**: Memiliki 2 parameter abadi:
    *   `w http.ResponseWriter`: Alat untuk mengirim balasan (response) ke user.
    *   `r *http.Request`: Alat untuk membaca data yang dikirim user (seperti membaca JSON atau parameter URL).

---

## 3. Fundamental 2: Concurrency (Nilai Jual Tertinggi Go)

Ini adalah alasan mengapa perusahaan besar beralih ke Go. Concurrency bukanlah sekadar "kecepatan", melainkan **efisiensi memanajemen banyak tugas sekaligus**.

Di dalam proyek ini, kita mempraktikkan 3 komponen utama *Concurrency*:

### A. Goroutine (`go func()`)
Ibarat Anda punya 10 pekerjaan rumah. Alih-alih mengerjakannya satu per satu, Anda menyewa 10 asisten untuk mengerjakan masing-masing tugas dalam waktu bersamaan. Menambahkan kata `go` di depan sebuah fungsi akan langsung membuat fungsi tersebut berjalan di latar belakang (paralel) dengan sangat ringan (jauh lebih ringan dari *Thread* milik Java atau OS).

### B. WaitGroup (`sync.WaitGroup`)
Jika Anda punya 10 asisten yang bekerja secara bersamaan, Anda butuh "Mandor" agar Anda tidak menutup kantor sebelum mereka semua selesai bekerja.
*   `wg.Add(1)` : Mandor mencatat ada 1 tugas baru.
*   `wg.Done()` : Asisten lapor tugasnya sudah selesai.
*   `wg.Wait()` : Kita menunggu sampai catatan Mandor menjadi 0 (semua asisten selesai).

### C. Channels (`chan`)
Karena para asisten (Goroutines) bekerja di ruang terpisah, mereka butuh alat komunikasi yang aman untuk menyerahkan hasil pekerjaannya ke fungsi utama. **Channel** adalah "pipa paralon" yang aman dari tabrakan data (*thread-safe*).
*   `resultsChan <- data`: Memasukkan data ke dalam pipa.
*   `data := <-resultsChan`: Mengambil data dari pipa.

---

## 4. Fundamental 3: Pekerja Latar Belakang (*Background Worker*)

Pada file `workers/checker.go`, kita menggunakan `time.Ticker`. Ini berbeda dengan `time.Sleep`. 
*   `time.Sleep`: Menghentikan seluruh program untuk sementara waktu.
*   `time.Ticker`: Sebuah jam weker yang berbunyi secara konstan di latar belakang, memicu fungsi lain berjalan, **tanpa menghentikan** fungsi utama (server HTTP) yang sedang berjalan melayani pengguna lain.

---

## 5. Cheat Sheet (Perintah Penting Terminal)

Sebagai pengembang Go, ini adalah perintah (*command*) yang akan sangat sering Anda gunakan:

*   `go run main.go` : Mengkompilasi dan menjalankan program secara langsung (untuk *development*).
*   `go build` : Membuat program menjadi aplikasi siap pakai (.exe di Windows). Aplikasi hasil *build* ini tidak lagi membutuhkan *Go* terinstal di komputer yang menjalankannya!
*   `go get <url>` : Mendownload *library* pihak ketiga dari internet.
*   `go mod tidy` : Membersihkan file `go.mod` (menghapus library yang tidak terpakai dan mendownload library yang kurang).
*   `swag init` : Menghasilkan (*generate*) ulang dokumentasi Swagger UI berdasarkan komentar di kode Anda.
