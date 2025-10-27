package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// User 数据库用户模型
type User struct {
	ID             int
	Username       string
	Email          string
	PasswordHash   string
	Phone          sql.NullString
	AvatarURL      sql.NullString
	Role           string
	Status         string
	DisabledReason sql.NullString
	DisabledAt     sql.NullTime
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// UserRepository 用户仓储接口
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

// userRepository 用户仓储实现
type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByUsername 根据用户名查询用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, phone, avatar_url,
		       role, status, disabled_reason, disabled_at, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.AvatarURL,
		&user.Role,
		&user.Status,
		&user.DisabledReason,
		&user.DisabledAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", username)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}

	return user, nil
}

// GetByEmail 根据邮箱查询用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, phone, avatar_url,
		       role, status, disabled_reason, disabled_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.AvatarURL,
		&user.Role,
		&user.Status,
		&user.DisabledReason,
		&user.DisabledAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", email)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	return user, nil
}

// GetByID 根据 ID 查询用户
func (r *userRepository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, phone, avatar_url,
		       role, status, disabled_reason, disabled_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.AvatarURL,
		&user.Role,
		&user.Status,
		&user.DisabledReason,
		&user.DisabledAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}

	return user, nil
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, phone, avatar_url, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.AvatarURL,
		user.Role,
		user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, phone = $3, avatar_url = $4,
		    role = $5, status = $6, disabled_reason = $7, disabled_at = $8
		WHERE id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Phone,
		user.AvatarURL,
		user.Role,
		user.Status,
		user.DisabledReason,
		user.DisabledAt,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
