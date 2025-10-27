package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/imehc/anch/generated"
	"github.com/imehc/anch/repository"
	"github.com/imehc/anch/util"
)

// AuthService 实现 AuthAPIServicer 接口
type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *util.JWTManager
}

// NewAuthService 创建新的认证服务实例
func NewAuthService(userRepo repository.UserRepository, jwtManager *util.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, loginRequest api.LoginRequest) (api.ImplResponse, error) {
	// 1. 查询用户（支持用户名或邮箱登录）
	var user *repository.User
	var err error

	// 尝试用用户名查询
	user, err = s.userRepo.GetByUsername(ctx, loginRequest.Username)
	if err != nil {
		// 如果用户名查询失败，尝试用邮箱查询
		user, err = s.userRepo.GetByEmail(ctx, loginRequest.Username)
		if err != nil {
			util.Error("登录失败: 用户不存在, 用户名/邮箱=%s", loginRequest.Username)
			return api.Response(http.StatusUnauthorized, "用户名或密码错误"), errors.New("invalid credentials")
		}
	}

	// 2. 检查用户状态
	if user.Status != "active" {
		reason := "账户未激活"
		if user.DisabledReason.Valid {
			reason = user.DisabledReason.String
		}
		util.Error("登录失败: 账户未激活, 用户ID=%d, 状态=%s, 原因=%s", user.ID, user.Status, reason)
		return api.Response(http.StatusForbidden, fmt.Sprintf("账户状态: %s, 原因: %s", user.Status, reason)), errors.New("account not active")
	}

	// 3. 验证密码
	if err := util.VerifyPassword(user.PasswordHash, loginRequest.Password); err != nil {
		util.Error("登录失败: 密码错误, 用户ID=%d, 用户名=%s", user.ID, user.Username)
		return api.Response(http.StatusUnauthorized, "用户名或密码错误"), errors.New("invalid password")
	}

	// 4. 生成 JWT token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		util.LogError("生成访问令牌", err, map[string]any{"用户ID": user.ID})
		return api.Response(http.StatusInternalServerError, "生成访问令牌失败"), fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		util.LogError("生成刷新令牌", err, map[string]any{"用户ID": user.ID})
		return api.Response(http.StatusInternalServerError, "生成刷新令牌失败"), fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. 返回认证响应
	auth := api.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int32(3600), // 1 hour
	}

	util.Info("用户登录成功, 用户ID=%d, 用户名=%s", user.ID, user.Username)
	return api.Response(http.StatusOK, auth), nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context) (api.ImplResponse, error) {
	// 从 context 中获取用户 ID（由认证中间件设置）
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		util.Error("获取当前用户失败: 未授权访问")
		return api.Response(http.StatusUnauthorized, "未授权"), errors.New("user not authenticated")
	}

	// 查询用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		util.LogError("查询用户信息", err, map[string]any{"用户ID": userID})
		return api.Response(http.StatusNotFound, "用户不存在"), fmt.Errorf("failed to get user: %w", err)
	}

	// 转换为 API 响应格式
	apiUser := api.User{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}

	// 可选字段
	if user.Phone.Valid {
		apiUser.Phone = user.Phone.String
	}
	if user.AvatarURL.Valid {
		apiUser.AvatarUrl = user.AvatarURL.String
	}
	if user.DisabledReason.Valid {
		apiUser.DisabledReason = user.DisabledReason.String
	}
	if user.DisabledAt.Valid {
		apiUser.DisabledAt = user.DisabledAt.Time
	}
	apiUser.CreatedAt = user.CreatedAt
	apiUser.UpdatedAt = user.UpdatedAt

	return api.Response(http.StatusOK, apiUser), nil
}

// RefreshToken 使用 refresh token 刷新 access token
func (s *AuthService) RefreshToken(ctx context.Context, req api.RefreshTokenRequest) (api.ImplResponse, error) {
	// 1. 验证 refresh token
	claims, err := s.jwtManager.VerifyToken(req.RefreshToken)
	if err != nil {
		util.Error("刷新令牌失败: 无效或过期的令牌, 错误=%v", err)
		return api.Response(http.StatusUnauthorized, "无效或过期的刷新令牌"), errors.New("invalid refresh token")
	}

	// 2. 验证 token 类型必须是 refresh token
	if claims.TokenType != "refresh" {
		util.Error("刷新令牌失败: 令牌类型错误, 类型=%s, 用户ID=%d", claims.TokenType, claims.UserID)
		return api.Response(http.StatusUnauthorized, "需要刷新令牌"), errors.New("invalid token type")
	}

	// 3. 验证用户是否仍然有效
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		util.LogError("查询用户", err, map[string]any{"用户ID": claims.UserID})
		return api.Response(http.StatusUnauthorized, "用户不存在"), errors.New("user not found")
	}

	// 4. 检查用户状态
	if user.Status != "active" {
		util.Error("刷新令牌失败: 账户未激活, 用户ID=%d, 状态=%s", user.ID, user.Status)
		return api.Response(http.StatusForbidden, "账户未激活"), errors.New("account not active")
	}

	// 5. 生成新的 access token
	newAccessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		util.LogError("生成访问令牌", err, map[string]any{"用户ID": user.ID})
		return api.Response(http.StatusInternalServerError, "生成访问令牌失败"), fmt.Errorf("failed to generate access token: %w", err)
	}

	// 6. 生成新的 refresh token
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		util.LogError("生成刷新令牌", err, map[string]any{"用户ID": user.ID})
		return api.Response(http.StatusInternalServerError, "生成刷新令牌失败"), fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 7. 返回新的 token
	auth := api.Auth{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int32(3600), // 1 hour
	}

	util.Info("令牌刷新成功, 用户ID=%d, 用户名=%s", user.ID, user.Username)
	return api.Response(http.StatusOK, auth), nil
}
