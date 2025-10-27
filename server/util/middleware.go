package util

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware JWT 认证中间件
type AuthMiddleware struct {
	jwtManager *JWTManager
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtManager *JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth 验证 Access Token 的中间件
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 记录请求信息
		Info("收到请求: method=%s, path=%s, remote_addr=%s", r.Method, r.URL.Path, r.RemoteAddr)

		// 从 Header 获取 token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			Error("认证失败: 缺少授权头, path=%s, remote_addr=%s", r.URL.Path, r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "缺少授权头")
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			Error("认证失败: 授权头格式错误, path=%s, remote_addr=%s", r.URL.Path, r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "授权头格式错误")
			return
		}

		tokenString := parts[1]

		// 验证 token
		claims, err := m.jwtManager.VerifyToken(tokenString)
		if err != nil {
			Error("认证失败: 无效或过期的令牌, path=%s, remote_addr=%s, error=%v", r.URL.Path, r.RemoteAddr, err)
			respondWithError(w, http.StatusUnauthorized, "无效或过期的令牌")
			return
		}

		// 验证 token 类型必须是 access token
		if claims.TokenType != "access" {
			Error("认证失败: 令牌类型错误, path=%s, remote_addr=%s, token_type=%s", r.URL.Path, r.RemoteAddr, claims.TokenType)
			respondWithError(w, http.StatusUnauthorized, "需要访问令牌")
			return
		}

		// 将用户信息注入到 context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)

		Info("认证成功: user_id=%d, username=%s, path=%s", claims.UserID, claims.Username, r.URL.Path)

		// 继续处理请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// respondWithError 返回错误响应
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
