package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/imehc/anch/config"
	"github.com/imehc/anch/db"
	api "github.com/imehc/anch/generated"
	"github.com/imehc/anch/repository"
	"github.com/imehc/anch/service"
	"github.com/imehc/anch/util"

	"github.com/go-chi/chi/v5"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "Path to config file")
	useEnv := flag.Bool("env", false, "Use environment variables instead of config file")
	flag.Parse()

	// 加载配置
	var cfg *config.Config
	var err error

	if *useEnv {
		log.Println("Loading configuration from environment variables...")
		cfg = config.LoadFromEnv()
	} else {
		log.Printf("Loading configuration from %s...\n", *configPath)
		cfg, err = config.Load(*configPath)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	}

	// 初始化日志系统
	log.Println("Initializing logger...")
	if err := util.InitLogger(util.LogConfig{
		LogDir:     "logs",
		MaxSizeMB:  10,
		MaxBackups: 30,
		MaxAgeDays: 30,
		Compress:   true,
		Level:      "info",
	}); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer util.CloseLogger()

	// 初始化数据库连接
	util.Info("Connecting to PostgreSQL database...")
	dbConn, err := db.NewPostgresDB(db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(dbConn)

	// 初始化 JWT 管理器
	jwtManager := util.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	// 初始化认证中间件
	authMiddleware := util.NewAuthMiddleware(jwtManager)

	// 初始化仓储层
	userRepo := repository.NewUserRepository(dbConn)
	diaryRepo := repository.NewDiaryRepository(dbConn)
	billRepo := repository.NewBillRepository(dbConn)

	// 初始化服务层
	authService := service.NewAuthService(userRepo, jwtManager)
	diaryService := service.NewDiaryService(diaryRepo)
	billService := service.NewBillService(billRepo)
	statsService := service.NewStatsService()

	// 创建 API 控制器
	authAPI := api.NewAuthAPIController(authService)
	diaryAPI := api.NewDiaryAPIController(diaryService)
	billAPI := api.NewBillAPIController(billService)
	statsAPI := api.NewStatsAPIController(statsService)

	// 创建新的路由器以应用选择性中间件
	router := chi.NewRouter()
	router.Use(api.Logger) // 使用生成的 Logger 中间件

	// 公开路由（不需要认证）
	publicRoutes := []string{"Login", "RefreshToken"}
	for _, route := range authAPI.OrderedRoutes() {
		if slices.Contains(publicRoutes, route.Name) {
			router.Method(route.Method, route.Pattern, route.HandlerFunc)
		}
	}

	// 受保护的路由（需要认证）
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		// Auth protected routes
		for _, route := range authAPI.OrderedRoutes() {
			if !slices.Contains(publicRoutes, route.Name) {
				r.Method(route.Method, route.Pattern, route.HandlerFunc)
			}
		}

		// Diary routes
		for _, route := range diaryAPI.OrderedRoutes() {
			r.Method(route.Method, route.Pattern, route.HandlerFunc)
		}

		// Bill routes
		for _, route := range billAPI.OrderedRoutes() {
			r.Method(route.Method, route.Pattern, route.HandlerFunc)
		}

		// Stats routes
		for _, route := range statsAPI.OrderedRoutes() {
			r.Method(route.Method, route.Pattern, route.HandlerFunc)
		}
	})

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	util.Info("Server started at http://%s", addr)
	util.Info("API documentation: http://%s/api", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
