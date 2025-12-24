package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"booksystem/internal/db"
)

type ShelfLayerHandler struct {
	db *gorm.DB
}

func NewShelfLayerHandler(db *gorm.DB) *ShelfLayerHandler {
	return &ShelfLayerHandler{db: db}
}

// CreateRequest 创建层数请求
type CreateShelfLayerRequest struct {
	BookshelfID int64  `json:"bookshelf_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

// UpdateShelfLayerRequest 更新层数请求
type UpdateShelfLayerRequest struct {
	BookshelfID int64  `json:"bookshelf_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

// Create 创建层数
func (h *ShelfLayerHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateShelfLayerRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	shelfLayer := db.ShelfLayer{
		BookshelfID: req.BookshelfID,
		Name:        req.Name,
	}

	if err := h.db.Create(&shelfLayer).Error; err != nil {
		Error(c, 400, "创建失败: "+err.Error())
		return
	}

	Success(c, shelfLayer)
}

// List 查询层数列表
func (h *ShelfLayerHandler) List(ctx context.Context, c *app.RequestContext) {
	var shelfLayers []db.ShelfLayer
	query := h.db.Model(&db.ShelfLayer{})

	if bookshelfID := c.Query("bookshelf_id"); bookshelfID != "" {
		query = query.Where("bookshelf_id = ?", bookshelfID)
	}

	if err := query.Find(&shelfLayers).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}

	Success(c, shelfLayers)
}

// Update 更新层数
func (h *ShelfLayerHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	var req UpdateShelfLayerRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	var shelfLayer db.ShelfLayer
	if err := h.db.First(&shelfLayer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, 404, "层数不存在")
		} else {
			Error(c, 500, "查询失败: "+err.Error())
		}
		return
	}

	if err := h.db.Model(&shelfLayer).Updates(map[string]interface{}{
		"bookshelf_id": req.BookshelfID,
		"name":         req.Name,
	}).Error; err != nil {
		Error(c, 400, "更新失败: "+err.Error())
		return
	}

	h.db.First(&shelfLayer, id)
	Success(c, shelfLayer)
}

// Delete 删除层数
func (h *ShelfLayerHandler) Delete(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	if err := h.db.Delete(&db.ShelfLayer{}, id).Error; err != nil {
		Error(c, 500, "删除失败: "+err.Error())
		return
	}

	Success(c, nil)
}
