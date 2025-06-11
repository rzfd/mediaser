# Google OAuth Setup untuk MediaShar

## Overview
MediaShar sekarang mendukung login dengan Google OAuth. Users dapat login menggunakan akun Google mereka tanpa perlu membuat akun terpisah.

## Setup Google OAuth

### 1. Google Cloud Console Setup

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Buat project baru atau pilih project yang sudah ada
3. Aktifkan Google+ API dan Google OAuth2 API
4. Pergi ke **APIs & Services > Credentials**
5. Klik **Create Credentials > OAuth client ID**
6. Pilih **Web application**
7. Tambahkan authorized domains:
   - `http://localhost:3000` (untuk development)
   - `https://yourdomain.com` (untuk production)
8. Copy Client ID yang dihasilkan

### 2. Frontend Configuration

Tambahkan Google Client ID ke environment variables:

**Development:**
```bash
# Buat file .env di frontend/
echo "REACT_APP_GOOGLE_CLIENT_ID=your-google-client-id-here.apps.googleusercontent.com" >> frontend/.env
```

**Docker/Production:**
Update `docker-compose.yml`:
```yaml
services:
  frontend:
    environment:
      - REACT_APP_GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
```

### 3. Backend Configuration

Tidak ada konfigurasi khusus diperlukan untuk backend karena menggunakan Google Token verification.

## Cara Kerja

1. User klik tombol "Sign in with Google"
2. Google Identity Services tampil untuk autentikasi
3. Setelah berhasil, Google mengirim JWT credential
4. Frontend mengirim credential ke backend `/api/auth/google`
5. Backend verifikasi token dengan Google API
6. Jika user belum ada, otomatis buat akun baru
7. Return JWT token untuk sesi login

## Testing

1. Setup Google Client ID di environment
2. Start aplikasi:
   ```bash
   cd frontend
   npm start
   ```
3. Buka http://localhost:3000/login
4. Klik tombol "Sign in with Google"
5. Login dengan akun Google
6. User akan otomatis login ke MediaShar

## Features

- ✅ Login dengan Google
- ✅ Auto-create user baru dari Google account
- ✅ Multi-bahasa support (EN, ID, ZH)
- ✅ Responsive design
- ✅ Error handling
- ✅ Integrasi dengan existing authentication system

## Security

- Token verification dilakukan di backend menggunakan Google API
- User data (email, nama) diambil dari verified Google token
- Tidak ada password disimpan untuk OAuth users
- JWT token tetap digunakan untuk sesi management

## Troubleshooting

**Error: "Invalid Google token"**
- Pastikan Google Client ID sudah benar
- Cek authorized domains di Google Cloud Console

**Error: "Google button tidak muncul"**
- Pastikan Google Client ID sudah di-set di environment
- Cek console browser untuk error JavaScript

**Error: "Failed to create user"**
- Cek koneksi database
- Pastikan backend running dan accessible 