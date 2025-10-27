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

// BillService 实现 BillAPIServicer 接口
type BillService struct {
	billRepo repository.BillRepository
}

// NewBillService 创建新的账单服务实例
func NewBillService(billRepo repository.BillRepository) *BillService {
	return &BillService{
		billRepo: billRepo,
	}
}

// ListBill 获取账单列表
func (s *BillService) ListBill(ctx context.Context, category, paymentMethod, month string) (api.ImplResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("获取账单列表失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	bills, err := s.billRepo.List(ctx, userID, category, paymentMethod, month)
	if err != nil {
		util.LogError("查询账单列表", err, map[string]any{"用户ID": userID})
		return api.Response(http.StatusInternalServerError, "查询账单列表失败"), fmt.Errorf("failed to list bills: %w", err)
	}

	var apiBills []api.Bill
	for _, b := range bills {
		apiBill := s.convertToAPIBill(b)
		files, _ := s.billRepo.GetFiles(ctx, b.ID)
		for _, f := range files {
			apiBill.Images = append(apiBill.Images, api.File{
				Id:        int32(f.ID),
				FileType:  f.FileType,
				FileUrl:   f.FileURL,
				CreatedAt: f.CreatedAt,
			})
		}
		apiBills = append(apiBills, apiBill)
	}

	return api.Response(http.StatusOK, apiBills), nil
}

// CreateBill 创建账单
func (s *BillService) CreateBill(ctx context.Context, billCreate api.BillCreate) (api.ImplResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("获取账单列表失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	bill := &repository.Bill{
		UserID:         userID,
		OriginalAmount: float64(billCreate.Amount.OriginalAmount),
		DiscountAmount: float64(billCreate.Amount.DiscountAmount),
		ActualAmount:   float64(billCreate.Amount.ActualAmount),
		PaymentMethod:  billCreate.PaymentMethod,
	}

	if billCreate.DiscountType != "" {
		bill.DiscountType = sql.NullString{String: billCreate.DiscountType, Valid: true}
	}
	if billCreate.Category != "" {
		bill.Category = sql.NullString{String: billCreate.Category, Valid: true}
	}
	if billCreate.Notes != "" {
		bill.Notes = sql.NullString{String: billCreate.Notes, Valid: true}
	}
	if billCreate.RelatedDiaryId > 0 {
		bill.RelatedDiaryID = sql.NullInt32{Int32: billCreate.RelatedDiaryId, Valid: true}
	}

	createdBill, err := s.billRepo.Create(ctx, bill)
	if err != nil {
		util.LogError("创建账单", err, map[string]any{"用户ID": userID})
		return api.Response(http.StatusInternalServerError, "创建账单失败"), fmt.Errorf("failed to create bill: %w", err)
	}

	apiBill := s.convertToAPIBill(createdBill)
	return api.Response(http.StatusCreated, apiBill), nil
}

// GetBill 查询单条账单
func (s *BillService) GetBill(ctx context.Context, id int32) (api.ImplResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("获取账单列表失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	bill, err := s.billRepo.GetByID(ctx, int(id), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Warn("查询账单: 账单不存在, 账单ID=%d, 用户ID=%d", id, userID)
			return api.Response(http.StatusNotFound, "账单不存在"), err
		}
		util.LogError("查询账单", err, map[string]any{"账单ID": id, "用户ID": userID})
		return api.Response(http.StatusInternalServerError, "查询账单失败"), fmt.Errorf("failed to get bill: %w", err)
	}

	apiBill := s.convertToAPIBill(bill)
	files, _ := s.billRepo.GetFiles(ctx, bill.ID)
	if files != nil {
		for _, f := range files {
			apiBill.Images = append(apiBill.Images, api.File{
				Id:        int32(f.ID),
				FileType:  f.FileType,
				FileUrl:   f.FileURL,
				CreatedAt: f.CreatedAt,
			})
		}
	}

	return api.Response(http.StatusOK, apiBill), nil
}

// UpdateBill 更新账单
func (s *BillService) UpdateBill(ctx context.Context, id int32, billCreate api.BillCreate) (api.ImplResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("更新账单失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	bill := &repository.Bill{
		ID:             int(id),
		UserID:         userID,
		OriginalAmount: float64(billCreate.Amount.OriginalAmount),
		DiscountAmount: float64(billCreate.Amount.DiscountAmount),
		ActualAmount:   float64(billCreate.Amount.ActualAmount),
		PaymentMethod:  billCreate.PaymentMethod,
	}

	if billCreate.DiscountType != "" {
		bill.DiscountType = sql.NullString{String: billCreate.DiscountType, Valid: true}
	}
	if billCreate.Category != "" {
		bill.Category = sql.NullString{String: billCreate.Category, Valid: true}
	}
	if billCreate.Notes != "" {
		bill.Notes = sql.NullString{String: billCreate.Notes, Valid: true}
	}
	if billCreate.RelatedDiaryId > 0 {
		bill.RelatedDiaryID = sql.NullInt32{Int32: billCreate.RelatedDiaryId, Valid: true}
	}

	updatedBill, err := s.billRepo.Update(ctx, bill)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Warn("更新账单: 账单不存在, 账单ID=%d, 用户ID=%d", id, userID)
			return api.Response(http.StatusNotFound, "账单不存在"), err
		}
		util.LogError("更新账单", err, map[string]any{"账单ID": id, "用户ID": userID})
		return api.Response(http.StatusInternalServerError, "更新账单失败"), fmt.Errorf("failed to update bill: %w", err)
	}

	apiBill := s.convertToAPIBill(updatedBill)
	return api.Response(http.StatusCreated, apiBill), nil
}

// DeleteBill 删除账单
func (s *BillService) DeleteBill(ctx context.Context, id int32) (api.ImplResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("删除账单失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	err := s.billRepo.Delete(ctx, int(id), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Warn("删除账单: 账单不存在, 账单ID=%d, 用户ID=%d", id, userID)
			return api.Response(http.StatusNotFound, "账单不存在"), err
		}
		util.LogError("删除账单", err, map[string]any{"账单ID": id, "用户ID": userID})
		return api.Response(http.StatusInternalServerError, "删除账单失败"), fmt.Errorf("failed to delete bill: %w", err)
	}

	return api.Response(http.StatusNoContent, nil), nil
}

// convertToAPIBill 转换为 API 响应格式
func (s *BillService) convertToAPIBill(b *repository.Bill) api.Bill {
	apiBill := api.Bill{
		Id: int32(b.ID),
		Amount: api.Amount{
			OriginalAmount: float32(b.OriginalAmount),
			DiscountAmount: float32(b.DiscountAmount),
			ActualAmount:   float32(b.ActualAmount),
		},
		PaymentMethod: b.PaymentMethod,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}

	if b.DiscountType.Valid {
		apiBill.DiscountType = b.DiscountType.String
	}
	if b.Category.Valid {
		apiBill.Category = b.Category.String
	}
	if b.Notes.Valid {
		apiBill.Notes = b.Notes.String
	}
	if b.RelatedDiaryID.Valid {
		apiBill.RelatedDiaryId = b.RelatedDiaryID.Int32
	}

	return apiBill
}
