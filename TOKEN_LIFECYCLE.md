# Token Lifecycle Documentation

## Giới thiệu
Hệ thống quản lý vòng đời Access Token và Refresh Token để giảm thiểu yêu cầu đăng nhập lại thường xuyên.

## Token Configuration
- **Access Token**: Hết hạn trong **15 phút** 
- **Refresh Token**: Hết hạn trong **7 ngày**

## Quy Trình Token Lifecycle

### 1. Đăng Ký (Register)
```
POST /register
Body: {
  "email": "user@example.com",
  "password": "password123"
}
Response: Gửi email xác minh (verification_token)
```

### 2. Xác Minh Tài Khoản (Verify)
```
GET /verify?token=<verification_token>
Response: Tài khoản được kích hoạt
```

### 3. Đăng Nhập (Login)
```
POST /login
Body: {
  "email": "user@example.com",
  "password": "password123"
}
Response: {
  "access_token": "eyJhbGc...",
  "refresh_token": "uuid-string"
}
```

### 4. Sử Dụng Access Token
```
GET /protected-endpoint
Header: Authorization: Bearer <access_token>
```

### 5. Làm Mới Access Token (Refresh)
Khi access token hết hạn:
```
POST /refresh
Body: {
  "refresh_token": "<refresh_token>"
}
Response: {
  "access_token": "new-access-token"
}
```

## Flow Diagram
```
1. Register → Email xác minh
   ↓
2. Verify Email → Tài khoản active
   ↓
3. Login → Access Token (15 phút) + Refresh Token (7 ngày)
   ↓
4. Sử dụng Access Token → Request API
   ↓
5. Access Token hết hạn?
   YES → Dùng Refresh Token để lấy cái mới
   NO  → Tiếp tục sử dụng
   ↓
6. Refresh Token hết hạn? → Login lại
```

## Database Schema
### users table
- id (PK)
- email (UNIQUE)
- password_hash
- verification_token
- is_verified
- status (pending, active, inactive)
- created_at, updated_at

### refresh_tokens table
- id (PK)
- user_id (FK → users)
- refresh_token (UNIQUE)
- expires_at
- is_revoked (default: false)
- created_at, updated_at

## Bảo Mật
1. Access token hết hạn nhanh (15 phút) → giảm rủi ro nếu bị lộ
2. Refresh token hết hạn chậm (7 ngày) → nhưng được lưu trong DB
3. Nếu refresh token bị lộ, có thể revoke thủ công
4. JWT Secret được mã hóa cứng (nên đưa vào environment variables)

## Ví Dụ Testing

### 1. Đăng Ký
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@gmail.com","password":"123456"}'
```

### 2. Xác Minh (thay token từ response trên)
```bash
curl "http://localhost:8080/verify?token=<verification_token>"
```

### 3. Đăng Nhập
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@gmail.com","password":"123456"}'
```

### 4. Làm Mới Token (thay refresh_token từ response ở bước 3)
```bash
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

## Lưu Ý
- Cần chạy SQL từ `migrations/001_create_tables.sql` để tạo schema
- Nếu `is_revoked` column không tồn tại, chạy: 
  ```sql
  ALTER TABLE refresh_tokens ADD COLUMN is_revoked BOOLEAN DEFAULT FALSE;
  ALTER TABLE refresh_tokens ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
  ```
- Để bảo mật tốt hơn, nên lưu refresh token trong HTTP-only cookie thay vì body
