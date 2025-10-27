package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// Diary 日记数据模型
type Diary struct {
	ID            int
	UserID        int
	Content       string
	Mood          sql.NullInt32
	Tags          []string
	RelatedBillID sql.NullInt32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// DiaryFile 日记文件关联
type DiaryFile struct {
	DiaryID int
	FileID  int
}

// DiaryRepository 日记数据访问接口
type DiaryRepository interface {
	Create(ctx context.Context, diary *Diary) (*Diary, error)
	GetByID(ctx context.Context, id, userID int) (*Diary, error)
	List(ctx context.Context, userID int, tag string, mood sql.NullInt32) ([]*Diary, error)
	Update(ctx context.Context, diary *Diary) (*Diary, error)
	Delete(ctx context.Context, id, userID int) error
	AddFile(ctx context.Context, diaryID, fileID int) error
	GetFiles(ctx context.Context, diaryID int) ([]*File, error)
}

type diaryRepository struct {
	db *sql.DB
}

// NewDiaryRepository 创建日记仓储实例
func NewDiaryRepository(db *sql.DB) DiaryRepository {
	return &diaryRepository{db: db}
}

// Create 创建日记
func (r *diaryRepository) Create(ctx context.Context, diary *Diary) (*Diary, error) {
	query := `
		INSERT INTO diaries (user_id, content, mood, tags, related_bill_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		diary.UserID,
		diary.Content,
		diary.Mood,
		pq.Array(diary.Tags),
		diary.RelatedBillID,
	).Scan(&diary.ID, &diary.CreatedAt, &diary.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return diary, nil
}

// GetByID 根据 ID 获取日记
func (r *diaryRepository) GetByID(ctx context.Context, id, userID int) (*Diary, error) {
	query := `
		SELECT id, user_id, content, mood, tags, related_bill_id, created_at, updated_at
		FROM diaries
		WHERE id = $1 AND user_id = $2
	`

	diary := &Diary{}
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&diary.ID,
		&diary.UserID,
		&diary.Content,
		&diary.Mood,
		pq.Array(&diary.Tags),
		&diary.RelatedBillID,
		&diary.CreatedAt,
		&diary.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return diary, nil
}

// List 获取日记列表
func (r *diaryRepository) List(ctx context.Context, userID int, tag string, mood sql.NullInt32) ([]*Diary, error) {
	query := `
		SELECT id, user_id, content, mood, tags, related_bill_id, created_at, updated_at
		FROM diaries
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	// 按标签过滤
	if tag != "" {
		query += " AND $" + string(rune('0'+argIndex)) + " = ANY(tags)"
		args = append(args, tag)
		argIndex++
	}

	// 按情绪过滤
	if mood.Valid {
		query += " AND mood = $" + string(rune('0'+argIndex))
		args = append(args, mood.Int32)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diaries []*Diary
	for rows.Next() {
		diary := &Diary{}
		err := rows.Scan(
			&diary.ID,
			&diary.UserID,
			&diary.Content,
			&diary.Mood,
			pq.Array(&diary.Tags),
			&diary.RelatedBillID,
			&diary.CreatedAt,
			&diary.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		diaries = append(diaries, diary)
	}

	return diaries, rows.Err()
}

// Update 更新日记
func (r *diaryRepository) Update(ctx context.Context, diary *Diary) (*Diary, error) {
	query := `
		UPDATE diaries
		SET content = $1, mood = $2, tags = $3, related_bill_id = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		diary.Content,
		diary.Mood,
		pq.Array(diary.Tags),
		diary.RelatedBillID,
		diary.ID,
		diary.UserID,
	).Scan(&diary.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return diary, nil
}

// Delete 删除日记
func (r *diaryRepository) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM diaries WHERE id = $1 AND user_id = $2`
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

// AddFile 添加日记文件关联
func (r *diaryRepository) AddFile(ctx context.Context, diaryID, fileID int) error {
	query := `INSERT INTO diary_files (diary_id, file_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, diaryID, fileID)
	return err
}

// GetFiles 获取日记的所有文件
func (r *diaryRepository) GetFiles(ctx context.Context, diaryID int) ([]*File, error) {
	query := `
		SELECT f.id, f.file_type, f.file_url, f.created_at
		FROM files f
		INNER JOIN diary_files df ON f.id = df.file_id
		WHERE df.diary_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, diaryID)
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
