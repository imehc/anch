package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Bill 账单数据模型
type Bill struct {
	ID              int
	UserID          int
	OriginalAmount  float64
	DiscountAmount  float64
	ActualAmount    float64
	DiscountType    sql.NullString
	PaymentMethod   string
	Category        sql.NullString
	Notes           sql.NullString
	RelatedDiaryID  sql.NullInt32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BillRepository 账单数据访问接口
type BillRepository interface {
	Create(ctx context.Context, bill *Bill) (*Bill, error)
	GetByID(ctx context.Context, id, userID int) (*Bill, error)
	List(ctx context.Context, userID int, category, paymentMethod, month string) ([]*Bill, error)
	Update(ctx context.Context, bill *Bill) (*Bill, error)
	Delete(ctx context.Context, id, userID int) error
	AddFile(ctx context.Context, billID, fileID int) error
	GetFiles(ctx context.Context, billID int) ([]*File, error)
}

type billRepository struct {
	db *sql.DB
}

// NewBillRepository 创建账单仓储实例
func NewBillRepository(db *sql.DB) BillRepository {
	return &billRepository{db: db}
}

// Create 创建账单
func (r *billRepository) Create(ctx context.Context, bill *Bill) (*Bill, error) {
	query := `
		INSERT INTO bills (user_id, original_amount, discount_amount, actual_amount,
			discount_type, payment_method, category, notes, related_diary_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		bill.UserID,
		bill.OriginalAmount,
		bill.DiscountAmount,
		bill.ActualAmount,
		bill.DiscountType,
		bill.PaymentMethod,
		bill.Category,
		bill.Notes,
		bill.RelatedDiaryID,
	).Scan(&bill.ID, &bill.CreatedAt, &bill.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return bill, nil
}

// GetByID 根据 ID 获取账单
func (r *billRepository) GetByID(ctx context.Context, id, userID int) (*Bill, error) {
	query := `
		SELECT id, user_id, original_amount, discount_amount, actual_amount,
			discount_type, payment_method, category, notes, related_diary_id,
			created_at, updated_at
		FROM bills
		WHERE id = $1 AND user_id = $2
	`

	bill := &Bill{}
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&bill.ID,
		&bill.UserID,
		&bill.OriginalAmount,
		&bill.DiscountAmount,
		&bill.ActualAmount,
		&bill.DiscountType,
		&bill.PaymentMethod,
		&bill.Category,
		&bill.Notes,
		&bill.RelatedDiaryID,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return bill, nil
}

// List 获取账单列表
func (r *billRepository) List(ctx context.Context, userID int, category, paymentMethod, month string) ([]*Bill, error) {
	query := `
		SELECT id, user_id, original_amount, discount_amount, actual_amount,
			discount_type, payment_method, category, notes, related_diary_id,
			created_at, updated_at
		FROM bills
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	// 按分类过滤
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	// 按支付方式过滤
	if paymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argIndex)
		args = append(args, paymentMethod)
		argIndex++
	}

	// 按月份过滤（格式：2025-10）
	if month != "" {
		query += fmt.Sprintf(" AND TO_CHAR(created_at, 'YYYY-MM') = $%d", argIndex)
		args = append(args, month)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*Bill
	for rows.Next() {
		bill := &Bill{}
		err := rows.Scan(
			&bill.ID,
			&bill.UserID,
			&bill.OriginalAmount,
			&bill.DiscountAmount,
			&bill.ActualAmount,
			&bill.DiscountType,
			&bill.PaymentMethod,
			&bill.Category,
			&bill.Notes,
			&bill.RelatedDiaryID,
			&bill.CreatedAt,
			&bill.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}

	return bills, rows.Err()
}

// Update 更新账单
func (r *billRepository) Update(ctx context.Context, bill *Bill) (*Bill, error) {
	query := `
		UPDATE bills
		SET original_amount = $1, discount_amount = $2, actual_amount = $3,
			discount_type = $4, payment_method = $5, category = $6, notes = $7,
			related_diary_id = $8, updated_at = NOW()
		WHERE id = $9 AND user_id = $10
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		bill.OriginalAmount,
		bill.DiscountAmount,
		bill.ActualAmount,
		bill.DiscountType,
		bill.PaymentMethod,
		bill.Category,
		bill.Notes,
		bill.RelatedDiaryID,
		bill.ID,
		bill.UserID,
	).Scan(&bill.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return bill, nil
}

// Delete 删除账单
func (r *billRepository) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM bills WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// AddFile 添加账单文件关联
func (r *billRepository) AddFile(ctx context.Context, billID, fileID int) error {
	query := `INSERT INTO bill_files (bill_id, file_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, billID, fileID)
	return err
}

// GetFiles 获取账单的所有文件
func (r *billRepository) GetFiles(ctx context.Context, billID int) ([]*File, error) {
	query := `
		SELECT f.id, f.file_type, f.file_url, f.created_at
		FROM files f
		INNER JOIN bill_files bf ON f.id = bf.file_id
		WHERE bf.bill_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		file := &File{}
		err := rows.Scan(&file.ID, &file.FileType, &file.FileURL, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}
