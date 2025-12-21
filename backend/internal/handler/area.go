package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"booksystem/internal/db"
)

type AreaHandler struct {
	db *gorm.DB
}

func NewAreaHandler(db *gorm.DB) *AreaHandler {
	return &AreaHandler{db: db}
}

// CreateRequest 创建区域请求
type CreateAreaRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateAreaRequest 更新区域请求
type UpdateAreaRequest struct {
	Name string `json:"name" binding:"required"`
}

// Create 创建区域
func (h *AreaHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateAreaRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	area := db.Area{Name: req.Name}
	if err := h.db.Create(&area).Error; err != nil {
		Error(c, 400, "创建失败: "+err.Error())
		return
	}

	Success(c, area)
}

// List 查询区域列表
func (h *AreaHandler) List(ctx context.Context, c *app.RequestContext) {
	var areas []db.Area
	if err := h.db.Find(&areas).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}
	Success(c, areas)
}

// Update 更新区域
func (h *AreaHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	var req UpdateAreaRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	var area db.Area
	if err := h.db.First(&area, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, 404, "区域不存在")
		} else {
			Error(c, 500, "查询失败: "+err.Error())
		}
		return
	}

	if err := h.db.Model(&area).Update("name", req.Name).Error; err != nil {
		Error(c, 400, "更新失败: "+err.Error())
		return
	}

	h.db.First(&area, id)
	Success(c, area)
}

// Delete 删除区域
func (h *AreaHandler) Delete(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	if err := h.db.Delete(&db.Area{}, id).Error; err != nil {
		Error(c, 500, "删除失败: "+err.Error())
		return
	}

	Success(c, nil)
}



