package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"booksystem/internal/db"
)

type BookshelfHandler struct {
	db *gorm.DB
}

func NewBookshelfHandler(db *gorm.DB) *BookshelfHandler {
	return &BookshelfHandler{db: db}
}

// CreateRequest 创建书架请求
type CreateBookshelfRequest struct {
	AreaID int64  `json:"area_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

// UpdateBookshelfRequest 更新书架请求
type UpdateBookshelfRequest struct {
	AreaID int64  `json:"area_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

// Create 创建书架
func (h *BookshelfHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateBookshelfRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	bookshelf := db.Bookshelf{
		AreaID: req.AreaID,
		Name:   req.Name,
	}

	if err := h.db.Create(&bookshelf).Error; err != nil {
		Error(c, 400, "创建失败: "+err.Error())
		return
	}

	Success(c, bookshelf)
}

// List 查询书架列表
func (h *BookshelfHandler) List(ctx context.Context, c *app.RequestContext) {
	var bookshelves []db.Bookshelf
	query := h.db.Model(&db.Bookshelf{})

	if areaID := c.Query("area_id"); areaID != "" {
		query = query.Where("area_id = ?", areaID)
	}

	if err := query.Find(&bookshelves).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}

	Success(c, bookshelves)
}

// Update 更新书架
func (h *BookshelfHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	var req UpdateBookshelfRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	var bookshelf db.Bookshelf
	if err := h.db.First(&bookshelf, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, 404, "书架不存在")
		} else {
			Error(c, 500, "查询失败: "+err.Error())
		}
		return
	}

	if err := h.db.Model(&bookshelf).Updates(map[string]interface{}{
		"area_id": req.AreaID,
		"name":    req.Name,
	}).Error; err != nil {
		Error(c, 400, "更新失败: "+err.Error())
		return
	}

	h.db.First(&bookshelf, id)
	Success(c, bookshelf)
}

// Delete 删除书架
func (h *BookshelfHandler) Delete(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	if err := h.db.Delete(&db.Bookshelf{}, id).Error; err != nil {
		Error(c, 500, "删除失败: "+err.Error())
		return
	}

	Success(c, nil)
}



