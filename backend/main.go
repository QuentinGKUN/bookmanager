package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"booksystem/internal/db"
	"booksystem/internal/handler"
	"booksystem/internal/middleware"
	"booksystem/internal/service"
)

func main() {
	// 初始化数据库
	// 支持从环境变量读取数据库路径，默认使用当前目录下的 booksystem.db
	dbPath := getEnv("DB_PATH", "booksystem.db")

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(database)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Redis
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0
	redisService, err := service.NewRedisService(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Redis features will be disabled.", err)
		redisService = nil
	}

	// 初始化Hertz服务器
	h := server.Default(
		server.WithHostPorts(":8089"),
	)

	// 注册中间件
	h.Use(middleware.CORS())

	// 注册路由
	registerRoutes(h, database, redisService)

	// 启动服务器
	h.Spin()
}

func registerRoutes(h *server.Hertz, db *gorm.DB, redisService *service.RedisService) {
	// 创建处理器
	bookHandler := handler.NewBookHandler(db)
	areaHandler := handler.NewAreaHandler(db)
	bookshelfHandler := handler.NewBookshelfHandler(db)
	shelfLayerHandler := handler.NewShelfLayerHandler(db)
	locationHandler := handler.NewLocationHandler(db)
	borrowHandler := handler.NewBorrowHandler(db, redisService)

	api := h.Group("/api/v1")
	{
		// 图书管理
		api.POST("/books", bookHandler.Create)
		api.GET("/books", bookHandler.List)
		api.GET("/books/barcode/:barcode", bookHandler.GetByBarcode)
		api.PUT("/books/:id", bookHandler.Update)
		api.DELETE("/books/:id", bookHandler.Delete)

		// 位置管理
		api.POST("/areas", areaHandler.Create)
		api.GET("/areas", areaHandler.List)
		api.PUT("/areas/:id", areaHandler.Update)
		api.DELETE("/areas/:id", areaHandler.Delete)

		api.POST("/bookshelves", bookshelfHandler.Create)
		api.GET("/bookshelves", bookshelfHandler.List)
		api.PUT("/bookshelves/:id", bookshelfHandler.Update)
		api.DELETE("/bookshelves/:id", bookshelfHandler.Delete)

		api.POST("/shelf-layers", shelfLayerHandler.Create)
		api.GET("/shelf-layers", shelfLayerHandler.List)
		api.PUT("/shelf-layers/:id", shelfLayerHandler.Update)
		api.DELETE("/shelf-layers/:id", shelfLayerHandler.Delete)

		api.GET("/locations/tree", locationHandler.GetTree)

		// 借阅管理（兼容旧接口）
		api.POST("/borrow", borrowHandler.Create)
		api.POST("/borrow/scan", borrowHandler.Scan)
		api.GET("/borrow/records", borrowHandler.List)
		api.POST("/borrow/get-borrower", borrowHandler.GetBorrowerByPhone)
		api.POST("/borrow/return", borrowHandler.Return)

		// 新的借阅API（使用Redis）
		api.POST("/borrow/user", borrowHandler.SetBorrowUser)
		api.GET("/borrow/user", borrowHandler.GetBorrowUser)
		api.POST("/borrow/book", borrowHandler.AddBorrowBook)
		api.DELETE("/borrow/book", borrowHandler.RemoveBorrowBook)
		api.POST("/borrow/complete", borrowHandler.CompleteBorrow)

		// 新的归还API（使用Redis）
		api.POST("/return/user", borrowHandler.SetReturnUser)
		api.GET("/return/user", borrowHandler.GetReturnUser)
		api.POST("/return/book", borrowHandler.AddReturnBook)
		api.DELETE("/return/book", borrowHandler.RemoveReturnBook)
		api.POST("/return/complete", borrowHandler.CompleteReturn)
	}

	// 健康检查
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
