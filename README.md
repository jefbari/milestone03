# LetterSquare API

REST API untuk platform review film — versi backend dari LetterSquare.
Dibangun dengan **Go 1.22**, **MySQL**, dan integrasi **Gemini AI** untuk rekomendasi film harian.

---

## Tech Stack

| Layer        | Tech                          |
|--------------|-------------------------------|
| Language     | Go 1.22                       |
| Database     | MySQL 8+                      |
| Auth         | JWT (golang-jwt/jwt v5)       |
| 3rd Party    | Gemini AI (gemini-2.5-flash)  |
| Testing      | testify (mock + assert)       |

---

## Endpoints (12 total)

### Auth
| Method | Endpoint              | Auth | Description         |
|--------|-----------------------|------|---------------------|
| POST   | /api/auth/register    | NO   | Register user baru  |
| POST   | /api/auth/login       | NO   | Login & dapat token |

### Movies
| Method | Endpoint              | Auth | Description                         |
|--------|-----------------------|------|-------------------------------------|
| GET    | /api/movies           | NO   | List semua film (search, genre, page)|
| GET    | /api/movies/:id       | NO   | Detail satu film                    |
| POST   | /api/movies           | YES  | Tambah film baru                    |
| PUT    | /api/movies/:id       | YES  | Update film                         |
| DELETE | /api/movies/:id       | YES  | Hapus film                          |

### Reviews
| Method | Endpoint                     | Auth | Description                    |
|--------|------------------------------|------|--------------------------------|
| GET    | /api/movies/:id/reviews      | NO   | Semua review untuk sebuah film |
| POST   | /api/movies/:id/reviews      | YES  | Tulis review (1x per film)     |
| GET    | /api/users/me/reviews        | YES  | Semua reviewku                 |
| PUT    | /api/reviews/:reviewId       | YES  | Update reviewku                |
| DELETE | /api/reviews/:reviewId       | YES  | Hapus reviewku                 |

### Watchlist
| Method | Endpoint                  | Auth | Description                  |
|--------|---------------------------|------|------------------------------|
| GET    | /api/watchlist            | YES  | Lihat watchlistku            |
| POST   | /api/watchlist/:movieId   | YES  | Tambah film ke watchlist     |
| DELETE | /api/watchlist/:movieId   | YES  | Hapus film dari watchlist    |

### AI Recommendation ⭐ (Gemini Integration — Q&A flow)
| Method | Endpoint                                     | Auth | Description                                              |
|--------|-----------------------------------------------|------|------------------------------------------------------------|
| POST   | /api/recommendations/start                    | YES  | Mulai sesi: dapat intro + pertanyaan pertama               |
| POST   | /api/recommendations/sessions/:id/answer      | YES  | Jawab pertanyaan → lanjut ke berikutnya, atau hasil akhir  |
| GET    | /api/recommendations/sessions/:id             | YES  | Lihat ulang sesi (termasuk hasil rekomendasi)              |

**Cara kerja:** bukan rekomendasi instan. Alurnya percakapan singkat:

1. `POST /start` → respons:
   ```json
   {
     "message": "let's find your next favorite film",
     "data": {
       "session_id": 1,
       "intro": "Your next favorite film is probably one you've never heard of.\nLet's find it.",
       "question": "What's your mood today — calm and curious, or craving something thrilling?",
       "step": 1,
       "total_steps": 3
     }
   }
   ```
2. `POST /sessions/1/answer` dengan `{"answer": "feeling adventurous"}` → dapat pertanyaan ke-2, lalu ke-3.
3. Setelah jawaban terakhir, Gemini dipanggil dengan konteks **watchlist user** + 3 jawaban tadi, dan mengembalikan rekomendasi final:
   ```json
   {
     "message": "here's what we found for you",
     "data": {
       "session_id": 1,
       "status": "completed",
       "recommendation": "1. Perfect Blue — because you liked Paprika..."
     }
   }
   ```

---

## Project Structure

```
letter-square-api/
├── cmd/api/
│   └── main.go                  # Entrypoint, dependency injection, router
├── config/
│   └── config.go                # Load env variables
├── internal/
│   ├── entity/                  # DB models (User, Movie, Review, Watchlist)
│   ├── dto/                     # Request/Response structs + validation
│   ├── apperror/                # Custom error types
│   ├── helper/                  # JWT, hashing, HTTP response utils
│   ├── repository/              # Interface + MySQL implementation
│   ├── service/                 # Business logic layer + unit tests
│   ├── handler/                 # HTTP handler layer
│   ├── middleware/              # JWT auth middleware
│   └── thirdparty/gemini/       # Gemini AI client
├── migrations/
│   └── 001_init.sql             # Schema SQL
└── .env.example
```

---

## Setup

### 1. Clone & Install deps

```bash
git clone <repo>
cd letter-square-api
go mod download
```

### 2. Setup Database

```bash
mysql -u root -p < migrations/001_init.sql
mysql -u root -p < migrations/002_recommendation_sessions.sql
```

### 3. Environment Variables

```bash
cp .env.example .env
# Edit .env sesuai konfigurasi lokal
```

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=secret
DB_NAME=letter_square
JWT_SECRET=your_jwt_secret_here
JWT_EXPIRY_HOUR=24
GEMINI_API_KEY=your_gemini_api_key_here   # dari https://aistudio.google.com
GEMINI_MODEL=gemini-2.5-flash
```

### 4. Run

```bash
go run cmd/api/main.go
```

### 5. Run Tests

```bash
go test ./internal/service/... -v
```

---

## Contoh Request

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","email":"john@example.com","password":"password123"}'
```

### Get Rekomendasi AI (Q&A flow)
```bash
# 1. Mulai sesi
curl -X POST http://localhost:8080/api/recommendations/start \
  -H "Authorization: Bearer <token>"

# 2. Jawab tiap pertanyaan (ulangi untuk 3 pertanyaan)
curl -X POST http://localhost:8080/api/recommendations/sessions/1/answer \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"answer":"feeling adventurous"}'
```

---

## Catatan Penting

> **Password Hashing**: Implementasi saat ini menggunakan SHA-256 untuk kompatibilitas environment.
> Untuk production, ganti dengan `golang.org/x/crypto/bcrypt` di `internal/helper/hash.go`.

---

## Unit Tests Coverage

| Test                           | Layer   | Status |
|-------------------------------|---------|--------|
| TestRegister_Success           | Service | PASS |
| TestRegister_EmailTaken        | Service | PASS |
| TestRegister_UsernameTaken     | Service | PASS |
| TestLogin_Success              | Service | PASS |
| TestLogin_InvalidPassword      | Service | PASS |
| TestLogin_UserNotFound         | Service | PASS |
| TestCreateReview_Success       | Service | PASS |
| TestCreateReview_MovieNotFound | Service | PASS |
| TestCreateReview_Duplicate     | Service | PASS |
| TestDeleteReview_Forbidden     | Service | PASS |
| TestUpdateReview_Success       | Service | PASS |
| TestStartSession_Success       | Service | PASS |
| TestAnswerQuestion_ProgressesToNextQuestion        | Service | PASS |
| TestAnswerQuestion_FinalAnswerGeneratesRecommendation | Service | PASS |
| TestAnswerQuestion_Forbidden    | Service | PASS |
| TestAnswerQuestion_AlreadyCompleted | Service | PASS |
| TestAnswerQuestion_EmptyAnswerRejected | Service | PASS |
