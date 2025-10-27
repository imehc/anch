package service

import (
	"context"
	"errors"
	"net/http"

	api "github.com/imehc/anch/generated"
)

// StatsService 实现 StatsAPIServicer 接口
type StatsService struct {
	// TODO: 添加 bill repository 依赖用于统计
}

// NewStatsService 创建新的统计服务实例
func NewStatsService() *StatsService {
	return &StatsService{}
}

// GetMonthlyStats 获取月度统计
func (s *StatsService) GetMonthlyStats(ctx context.Context, month string) (api.ImplResponse, error) {
	// 验证用户认证
	_, ok := ctx.Value("user_id").(int)
	if !ok {
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// TODO: 实现月度统计逻辑
	// 返回空的统计数据
	stats := api.Stats{
		Month: month,
	}

	return api.Response(http.StatusOK, stats), nil
}

// GetCategoryStats 获取分类支出占比
func (s *StatsService) GetCategoryStats(ctx context.Context, month string) (api.ImplResponse, error) {
	_, ok := ctx.Value("user_id").(int)
	if !ok {
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// TODO: 实现分类统计逻辑
	categoryStats := make(map[string]float32)

	return api.Response(http.StatusOK, categoryStats), nil
}

// GetDiscountStats 获取优惠类型占比
func (s *StatsService) GetDiscountStats(ctx context.Context, month string) (api.ImplResponse, error) {
	_, ok := ctx.Value("user_id").(int)
	if !ok {
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// TODO: 实现优惠统计逻辑
	discountStats := make(map[string]float32)

	return api.Response(http.StatusOK, discountStats), nil
}

// GetTrendStats 获取每日支出和情绪趋势
func (s *StatsService) GetTrendStats(ctx context.Context, month string) (api.ImplResponse, error) {
	_, ok := ctx.Value("user_id").(int)
	if !ok {
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// TODO: 实现趋势统计逻辑
	stats := api.Stats{
		Month: month,
	}

	return api.Response(http.StatusOK, stats), nil
}
