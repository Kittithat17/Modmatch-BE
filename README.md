# modmatch-be

Backend API — Go + Gin + PostgreSQL (pgx) + sqlc

## โครงสร้าง

```
cmd/api/main.go          จุดเริ่มแอป — ประกอบ dependency + routes
internal/
  config/                โหลดค่าจาก env / .env
  database/              connection pool ของ Postgres
  middleware/            CORS และ middleware อื่นๆ
  auth/                  สร้าง/ตรวจ JWT
  <domain>/              แต่ละ feature (เช่น user) ดูหัวข้อ "เพิ่ม domain ใหม่"
pkg/logger/              logger กลาง (slog)
db/migrations/           ไฟล์ SQL สร้าง/แก้ table (0001_xxx.sql, 0002_xxx.sql, ...)
db/queries/              query SQL ให้ sqlc generate โค้ด
```

## เริ่มต้น (ครั้งแรก)

```bash
cp .env.example .env
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go mod download
```

## รัน

```bash
# แบบ dev: รันแค่ postgres ใน docker แล้วรันแอปในเครื่อง
docker compose up -d postgres
go run ./cmd/api

# แบบรันทั้งระบบใน docker
docker compose up -d --build
```

เช็คว่าแอปทำงาน: `curl http://localhost:8080/health`

> port 5436 ในเครื่องชนกับโปรแกรมอื่น? แก้ `DB_PORT` ใน `.env` และเลข port ฝั่งซ้ายใน `docker-compose.yml`

## เพิ่ม domain ใหม่ (เช่น user)

ใช้แพทเทิร์น handler → service → repository

1. `db/migrations/0001_create_users.sql` — สร้าง table
2. `db/queries/users.sql` — เขียน query (`-- name: GetUser :one`)
3. เพิ่ม entry ใน `sqlc.yaml` (มีตัวอย่างอยู่ในไฟล์) แล้วรัน `sqlc generate`
4. สร้าง `internal/user/`
   - `model.go` — struct ของ domain, error, `Repository` interface
   - `repository.go` — เรียกโค้ดจาก sqlc แล้วแปลงเป็น domain struct
   - `service.go` — business logic
   - `handler.go` — รับ HTTP + `RegisterRoutes`
5. ต่อสายใน `cmd/api/main.go` แล้วเรียก `RegisterRoutes(apiV1)`

> migration ใน `db/migrations` จะรันอัตโนมัติเฉพาะตอนสร้าง volume ครั้งแรก
> ถ้า DB มีอยู่แล้ว ให้รันเอง:
> `docker exec -i modmatch-postgres psql -U postgres -d modmatch < db/migrations/0002_xxx.sql`
> หรือล้าง DB ใหม่หมดด้วย `docker compose down -v` (ข้อมูลหาย)

## CI

GitHub Actions ([.github/workflows/ci.yml](.github/workflows/ci.yml)) รันทุก push เข้า main และทุก PR
เช็ค `go mod tidy`, `gofmt`, build, vet, test — ก่อน push ลองรันในเครื่องได้:

```bash
go mod tidy && go fmt ./... && go vet ./... && go test ./...
```

## คำสั่งที่ใช้บ่อย

| คำสั่ง | ทำอะไร |
|--------|--------|
| `go run ./cmd/api` | รันแอป |
| `go build ./...` | เช็คว่า compile ผ่าน |
| `go mod tidy` | จัดระเบียบ dependency |
| `go fmt ./...` / `go vet ./...` | จัดรูปแบบ / ตรวจโค้ด |
| `sqlc generate` | generate โค้ดจาก SQL (รันทุกครั้งที่แก้ .sql) |
| `docker compose logs -f app` | ดู log แอป |
| `docker compose down` | หยุดทั้งหมด (ข้อมูลยังอยู่) |
| `docker exec -it modmatch-postgres psql -U postgres -d modmatch` | เข้า psql |
