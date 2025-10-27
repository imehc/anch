package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	api "github.com/imehc/anch/generated"
	"github.com/imehc/anch/repository"
	"github.com/imehc/anch/util"
)

// DiaryService 实现 DiaryAPIServicer 接口
type DiaryService struct {
	diaryRepo repository.DiaryRepository
}

// NewDiaryService 创建新的日记服务实例
func NewDiaryService(diaryRepo repository.DiaryRepository) *DiaryService {
	return &DiaryService{
		diaryRepo: diaryRepo,
	}
}

// ListDiary 获取日记列表
func (s *DiaryService) ListDiary(ctx context.Context, tag string, mood int32) (api.ImplResponse, error) {
	// 从 context 获取用户 ID
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("获取日记列表失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// 构建查询条件
	var moodFilter sql.NullInt32
	if mood > 0 {
		moodFilter = sql.NullInt32{Int32: mood, Valid: true}
	}

	// 查询日记列表
	diaries, err := s.diaryRepo.List(ctx, userID, tag, moodFilter)
	if err != nil {
		util.LogError("查询日记列表", err, map[string]any{"用户ID": userID})
		return api.Response(http.StatusInternalServerError, "查询日记列表失败"), fmt.Errorf("failed to list diaries: %w", err)
	}

	// 转换为 API 响应格式
	var apiDiaries []api.Diary
	for _, d := range diaries {
		apiDiary := api.Diary{
			Id:        int32(d.ID),
			Content:   d.Content,
			Tags:      d.Tags,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		}

		if d.Mood.Valid {
			apiDiary.Mood = d.Mood.Int32
		}

		if d.RelatedBillID.Valid {
			apiDiary.RelatedBillId = d.RelatedBillID.Int32
		}

		// 获取文件列表
		files, _ := s.diaryRepo.GetFiles(ctx, d.ID)
		if files != nil {
			for _, f := range files {
				apiDiary.Files = append(apiDiary.Files, api.File{
					Id:        int32(f.ID),
					FileType:  f.FileType,
					FileUrl:   f.FileURL,
					CreatedAt: f.CreatedAt,
				})
			}
		}

		apiDiaries = append(apiDiaries, apiDiary)
	}

	return api.Response(http.StatusOK, apiDiaries), nil
}

// CreateDiary 创建日记
func (s *DiaryService) CreateDiary(ctx context.Context, diaryCreate api.DiaryCreate) (api.ImplResponse, error) {
	// 从 context 获取用户 ID
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("创建日记失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// 转换为仓储模型
	diary := &repository.Diary{
		UserID:  userID,
		Content: diaryCreate.Content,
		Tags:    diaryCreate.Tags,
	}

	if diaryCreate.Mood > 0 {
		diary.Mood = sql.NullInt32{Int32: diaryCreate.Mood, Valid: true}
	}

	if diaryCreate.RelatedBillId > 0 {
		diary.RelatedBillID = sql.NullInt32{Int32: diaryCreate.RelatedBillId, Valid: true}
	}

	// 创建日记
	createdDiary, err := s.diaryRepo.Create(ctx, diary)
	if err != nil {
		util.LogError("创建日记", err, map[string]any{"用户ID": userID})
		return api.Response(http.StatusInternalServerError, "创建日记失败"), fmt.Errorf("failed to create diary: %w", err)
	}

	// 转换为 API 响应格式
	apiDiary := api.Diary{
		Id:        int32(createdDiary.ID),
		Content:   createdDiary.Content,
		Tags:      createdDiary.Tags,
		CreatedAt: createdDiary.CreatedAt,
		UpdatedAt: createdDiary.UpdatedAt,
	}

	if createdDiary.Mood.Valid {
		apiDiary.Mood = createdDiary.Mood.Int32
	}

	if createdDiary.RelatedBillID.Valid {
		apiDiary.RelatedBillId = createdDiary.RelatedBillID.Int32
	}

	return api.Response(http.StatusCreated, apiDiary), nil
}

// GetDiary 查询单条日记
func (s *DiaryService) GetDiary(ctx context.Context, id int32) (api.ImplResponse, error) {
	// 从 context 获取用户 ID
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("查询日记失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// 查询日记
	diary, err := s.diaryRepo.GetByID(ctx, int(id), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Warn("查询日记: 日记不存在, 日记ID=%d, 用户ID=%d", id, userID)
			return api.Response(http.StatusNotFound, "日记不存在"), err
		}
		util.LogError("查询日记", err, map[string]any{"日记ID": id, "用户ID": userID})
		return api.Response(http.StatusInternalServerError, "查询日记失败"), fmt.Errorf("failed to get diary: %w", err)
	}

	// 转换为 API 响应格式
	apiDiary := api.Diary{
		Id:        int32(diary.ID),
		Content:   diary.Content,
		Tags:      diary.Tags,
		CreatedAt: diary.CreatedAt,
		UpdatedAt: diary.UpdatedAt,
	}

	if diary.Mood.Valid {
		apiDiary.Mood = diary.Mood.Int32
	}

	if diary.RelatedBillID.Valid {
		apiDiary.RelatedBillId = diary.RelatedBillID.Int32
	}

	// 获取文件列表
	files, _ := s.diaryRepo.GetFiles(ctx, diary.ID)
	if files != nil {
		for _, f := range files {
			apiDiary.Files = append(apiDiary.Files, api.File{
				Id:        int32(f.ID),
				FileType:  f.FileType,
				FileUrl:   f.FileURL,
				CreatedAt: f.CreatedAt,
			})
		}
	}

	return api.Response(http.StatusOK, apiDiary), nil
}

// UpdateDiary 更新日记
func (s *DiaryService) UpdateDiary(ctx context.Context, id int32, diaryCreate api.DiaryCreate) (api.ImplResponse, error) {
	// 从 context 获取用户 ID
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("更新日记失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// 转换为仓储模型
	diary := &repository.Diary{
		ID:      int(id),
		UserID:  userID,
		Content: diaryCreate.Content,
		Tags:    diaryCreate.Tags,
	}

	if diaryCreate.Mood > 0 {
		diary.Mood = sql.NullInt32{Int32: diaryCreate.Mood, Valid: true}
	}

	if diaryCreate.RelatedBillId > 0 {
		diary.RelatedBillID = sql.NullInt32{Int32: diaryCreate.RelatedBillId, Valid: true}
	}

	// 更新日记
	updatedDiary, err := s.diaryRepo.Update(ctx, diary)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Warn("更新日记: 日记不存在, 日记ID=%d, 用户ID=%d", id, userID)
			return api.Response(http.StatusNotFound, "日记不存在"), err
		}
		util.LogError("更新日记", err, map[string]any{"日记ID": id, "用户ID": userID})
		return api.Response(http.StatusInternalServerError, "更新日记失败"), fmt.Errorf("failed to update diary: %w", err)
	}

	// 转换为 API 响应格式
	apiDiary := api.Diary{
		Id:        int32(updatedDiary.ID),
		Content:   updatedDiary.Content,
		Tags:      updatedDiary.Tags,
		CreatedAt: updatedDiary.CreatedAt,
		UpdatedAt: updatedDiary.UpdatedAt,
	}

	if updatedDiary.Mood.Valid {
		apiDiary.Mood = updatedDiary.Mood.Int32
	}

	if updatedDiary.RelatedBillID.Valid {
		apiDiary.RelatedBillId = updatedDiary.RelatedBillID.Int32
	}

	return api.Response(http.StatusCreated, apiDiary), nil
}

// DeleteDiary 删除日记
func (s *DiaryService) DeleteDiary(ctx context.Context, id int32) (api.ImplResponse, error) {
	// 从 context 获取用户 ID
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("删除日记失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// 删除日记
	err := s.diaryRepo.Delete(ctx, int(id), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Warn("删除日记: 日记不存在, 日记ID=%d, 用户ID=%d", id, userID)
			return api.Response(http.StatusNotFound, "日记不存在"), err
		}
		util.LogError("删除日记", err, map[string]any{"日记ID": id, "用户ID": userID})
		return api.Response(http.StatusInternalServerError, "删除日记失败"), fmt.Errorf("failed to delete diary: %w", err)
	}

	return api.Response(http.StatusNoContent, nil), nil
}
