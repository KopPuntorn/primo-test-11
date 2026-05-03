# Primo Product API


## โครงสร้างโปรเจ็ค
- `cmd/api`: จุดเริ่มต้นของแอปพลิเคชัน (main.go)
- `internal/domain`: เก็บ Entity, Interface และ Error ต่างๆ
- `internal/usecase`: ส่วนของ Business Logic
- `internal/repository`: ส่วนที่ติดต่อกับฐานข้อมูล
- `internal/handler`: ส่วนที่รับ HTTP Request (Echo)
- `internal/database`: การตั้งค่าและเชื่อมต่อฐานข้อมูล
- `docs`: ไฟล์ Swagger documentation ที่ถูก Generate ขึ้นมา

## วิธีการติดตั้งและรันโปรเจ็ค

### 1. การเตรียมสภาพแวดล้อม
- ติดตั้ง [Go 1.25+](https://golang.org/dl/)
- ติดตั้ง [Docker Desktop](https://www.docker.com/products/docker-desktop)

### 2. ตั้งค่า Environment Variables
คัดลอกไฟล์ `.env.example` ไปเป็น `.env`
```bash
cp .env.example .env
```

### 3. เริ่มทำงานฐานข้อมูล (Docker)
รัน PostgreSQL ผ่าน Docker Compose:
```bash
docker-compose up -d postgres
```

### 4. รันแอปพลิเคชัน
รันแอปพลิเคชันในเครื่อง (Local):
```bash
go run cmd/api/main.go
```

## การใช้งาน Swagger
หลังจากรันแอปพลิเคชันแล้ว สามารถเข้าดู API Documentation ได้ที่:
http://localhost:5000/api-docs/index.html

## การทดสอบ (Testing)
โปรเจ็คนี้ใช้การทดสอบ 3 ระดับ:

### 1. Unit Test
เป็นการทดสอบ Logic ภายใน โดยใช้ Mock (ไม่ต้องต่อฐานข้อมูล):
```bash
# รันเทสต์ทั้งหมด
go test ./...

# รันเทสต์แบบละเอียดเฉพาะส่วน
go test -v ./internal/usecase/...
go test -v ./internal/handler/...
```

### 2. Integration & E2E Test
เป็นการทดสอบที่ต้องเชื่อมต่อกับฐานข้อมูลจริง (PostgreSQL)

#### **เตรียมฐานข้อมูลสำหรับเทสต์:**
สร้างฐานข้อมูลชื่อ `product_test` (ทำเพียงครั้งเดียว):
```bash
docker exec -it primo-test-11-postgres-1 psql -U postgres -c "CREATE DATABASE product_test;"
```

#### **การรัน Integration Test (Repository):**

**Bash (Linux/macOS/Git Bash):**
```bash
DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=product_test go test -v ./internal/repository/...
```

#### **การรัน E2E Test (API Flow):**
เป็นการทดสอบ Flow ของ API ตั้งแต่ HTTP Request จนถึง Database:
**Bash (Linux/macOS/Git Bash):**
```bash
DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=product_test go test -v ./internal/handler/product_e2e_test.go ./internal/handler/product_handler.go ./internal/handler/product_handler_test.go
```


## API Endpoints
- `POST /product`: สร้างสินค้าใหม่
- `PATCH /product/:id`: แก้ไขข้อมูลสินค้าบางส่วน
