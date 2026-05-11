# Go + Gin + GORM Convention Report

Generated: 2026-05-11
Stack: Go 1.21+, Gin v1.10+, GORM v2, golang-jwt/jwt v5

---

## 1. Project Structure (Clean Architecture)

Based on [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
project/
  cmd/               # main() entrypoints per binary (cmd/app, cmd/cli)
  internal/          # PRIVATE code — compiler-enforced non-importable
    domain/          # entities, value objects, interfaces
    repository/      # data access layer (implements domain interfaces)
    service/         # business logic layer
    handler/         # HTTP/transport layer (Gin handlers)
    middleware/      # Gin middleware (auth, logging, CORS, rate-limit)
    config/          # config structs, env binding
  pkg/               # PUBLIC library code (safe for external use)
  api/                # OpenAPI/Swagger specs, proto files
  configs/           # config.yaml, config.toml (non-code configs)
  deployments/       # Docker, K8s, CI configs
  test/               # integration test data
  go.mod / go.sum
```

**Key distinction:**
- `internal/` = private packages, compiler enforces boundaries
- `pkg/` = explicitly public libraries
- For REST API projects: `internal/` covers 90% of needs

**No `/src` directory** — that is a Java pattern. Avoid.

**File naming conventions:**
- `user_repository.go` — lowercase, snake_case
- `UserRepository` — exported types/interfaces PascalCase
- `user_service.go` — service file
- `user_handler.go` — HTTP handler file
- `user_test.go` — test file in same package
- `*_test.go` suffix is reserved for Go test files

---

## 2. Gin Framework Patterns

### Routing
```go
// Group routes for shared middleware
api := router.Group("/api/v1")
api.Use(middleware.Auth(), middleware.Logger())

api.POST("/users", handler.CreateUser)
api.GET("/users/:id", handler.GetUser)
api.PUT("/users/:id", handler.UpdateUser)
api.DELETE("/users/:id", handler.DeleteUser)
```

### Binding & Validation
```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// Gin binding tags:
// json:"field" binding:"required,min=2,email"          → JSON body
// form:"field" binding:"required"                        → form-data
// uri:"id"    binding:"required"                        → path params
// query:"page"                                          → query params
```

**Use `go-playground/validator` v10+** — Gin integrates it natively via binding tags.

### Middleware Pattern
```go
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // validate, set context, abort if invalid
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### Context Values
```go
// Set
c.Set("user_id", claims.UserID)
// Get (in downstream handler)
userID, exists := c.Get("user_id")
if !exists { /* handle missing */ }
```

### Graceful Shutdown
```go
srv := &http.Server{Addr: ":8080", Handler: router}
go srv.ListenAndServe()
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

---

## 3. GORM v2 Patterns

### Model Definition
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`           // soft delete
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:255"`
    Password  string         `gorm:"size:255"`
    Age       int            `gorm:"default:0"`
}

// Tags:
// primaryKey, uniqueIndex, index, not null, default:, size:, column:
// type:, serializer:, - (ignore field)
```

### Auto-Migrate
```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
if err != nil {
    log.Fatal("failed to connect")
}
db.AutoMigrate(&User{}, &Product{})  // safe to call repeatedly
```

### Transactions
```go
func CreateUserWithProfile(db *gorm.DB, user User, profile Profile) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        profile.UserID = user.ID
        if err := tx.Create(&profile).Error; err != nil {
            return err  // automatic rollback
        }
        return nil
    })
}
```

### Hooks
```go
// BeforeCreate, AfterCreate, BeforeUpdate, AfterUpdate, etc.
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.Password = hashPassword(u.Password)  // hash before save
    return nil
}
```

### Queries
```go
// Create
db.Create(&user)

// First / Find
db.First(&user, id)
db.Where("name = ?", name).First(&user)

// Batch
db.CreateInBatches(users, 100)

// Soft delete (uses DeletedAt)
db.Delete(&user)   // sets DeletedAt, does not hard delete

// WithContext (for cancellation)
db.WithContext(ctx).First(&user, id)
```

---

## 4. JWT Implementation (golang-jwt/jwt v5)

### Token Creation
```go
import (
    "github.com/golang-jwt/jwt/v5"
)

// Claims struct
type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

// Sign with HS256 (symmetric — same key signs and verifies)
token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
    RegisteredClaims: jwt.RegisteredClaims{
        Subject:   strconv.FormatUint(uint64(userID), 10),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
    UserID: userID,
})
signedToken, err := token.SignedString([]byte(secretKey))
```

### Token Validation
```go
// Parse with claims extraction
claims := &Claims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
    }
    return []byte(secretKey), nil
})

if err != nil || !token.Valid {
    return nil, err
}
```

### Refresh Token Pattern
- Access token: short-lived (15 min)
- Refresh token: long-lived (7 days), stored in DB or httpOnly cookie
- On expiry: client sends refresh token, issue new access token
- Revoke by deleting from DB or marking as used

**DO NOT use `alg=none`** — `jwt.UnsafeAllowNoneSignatureType` must never be used in production.

### Gin JWT Middleware
```go
func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

---

## 5. Error Handling (Go Idiomatic)

```go
// Sentinel errors
var (
    ErrUserNotFound   = errors.New("user not found")
    ErrInvalidInput   = errors.New("invalid input")
    ErrUnauthorized  = errors.New("unauthorized")
)

// Wrap with context
if err := db.First(&user, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return ErrUserNotFound
    }
    return fmt.Errorf("fetch user: %w", err)
}

// Check error type
if errors.Is(err, ErrUserNotFound) { /* handle */ }

// Custom error types
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Cause   error  `json:"-"`
}
func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Cause }
```

**HTTP error responses in Gin:**
```go
c.JSON(404, gin.H{
    "error":   "user_not_found",
    "message": "User with ID 123 not found",
})
```

---

## 6. Repository / Service / Handler Layers

### Domain (internal/domain/)
```go
type User struct {
    ID       uint
    Name     string
    Email    string
}
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    Update(user *User) error
    Delete(id uint) error
}
```

### Repository (internal/repository/)
```go
type userRepository struct {
    db *gorm.DB
}
func NewUserRepository(db *gorm.DB) *userRepository {
    return &userRepository{db: db}
}
// Implement domain.UserRepository interface
func (r *userRepository) GetByID(id uint) (*domain.User, error) {
    var user domain.User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### Service (internal/service/)
```go
type userService struct {
    repo domain.UserRepository
}
func NewUserService(repo domain.UserRepository) *userService {
    return &userService{repo: repo}
}
func (s *userService) Create(req *CreateUserRequest) (*domain.User, error) {
    // Business logic, validation
    user := &domain.User{Name: req.Name, Email: req.Email}
    if err := s.repo.Create(user); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    return user, nil
}
```

### Handler (internal/handler/)
```go
type UserHandler struct {
    service *service.UserService
}
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    user, err := h.service.Create(&req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, user)
}
```

**Dependency injection wiring in main.go:**
```go
db := initDB()
repo := repository.NewUserRepository(db)
svc := service.NewUserService(repo)
handler := handler.NewUserHandler(svc)
router := gin.Default()
registerRoutes(router, handler)
```

---

## 7. Channel-Based Worker Patterns

### Basic Worker Pool
```go
type Job struct {
    ID   int
    Data string
}

func worker(id int, jobs <-chan Job, results chan<- string) {
    for job := range jobs {
        result := process(job)
        results <- result
    }
}

func StartWorkerPool(numWorkers int, numJobs int) {
    jobs := make(chan Job, numJobs)
    results := make(chan string, numJobs)

    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("job-%d", j)}
    }
    close(jobs)

    for a := 1; a <= numJobs; a++ {
        <-results
    }
}
```

### Background Worker with Context Cancellation
```go
func StartEmailWorker(ctx context.Context, emailChan <-chan Email) {
    for {
        select {
        case <-ctx.Done():
            log.Println("email worker shutting down")
            return
        case email := <-emailChan:
            sendEmail(email)
        }
    }
}

// Graceful shutdown
ctx, cancel := context.WithCancel(context.Background())
go StartEmailWorker(ctx, emailChan)
// On shutdown signal:
cancel()
```

**Key patterns:**
- `make(chan Type, bufferSize)` — buffered channel for backpressure
- `for job := range ch` — auto-close loop
- `select { case ...; default: }` — non-blocking select
- `context.Context` — cancellation and timeout propagation

---

## 8. Environment Config (.env)

**Use `godotenv` for local dev:**
```bash
# .env
DATABASE_URL=sqlite.db
JWT_SECRET=your-secret-key
JWT_EXPIRY=15m
PORT=8080
GIN_MODE=debug
```

```go
import (
    "github.com/joho/godotenv"
    "os"
)

func LoadConfig() (string, error) {
    _ = godotenv.Load()  // loads .env, ignore error if missing
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET not set")
    }
    return secret, nil
}
```

**Config struct pattern:**
```go
type Config struct {
    DatabaseURL string
    JWTSecret   string
    JWTExpiry   time.Duration
    GINMode     string
    Port        string
}

func (c *Config) GetPort() string {
    if c.Port == "" {
        return "8080"
    }
    return c.Port
}
```

**Never hardcode secrets.** Load from env or env file at startup. Validate required env vars on app startup — fail fast.

---

## 9. Sources

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- [gorm.io/docs](https://gorm.io/docs)
- [golang-jwt.github.io/jwt/usage/create](https://golang-jwt.github.io/jwt/usage/create/)
- [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)
