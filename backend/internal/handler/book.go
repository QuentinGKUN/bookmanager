package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"booksystem/internal/db"
)

type BookHandler struct {
	db *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{db: db}
}

// CreateRequest 创建图书请求
type CreateBookRequest struct {
	Barcode      string   `json:"barcode" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Quantity     int      `json:"quantity" binding:"required"`
	InStock      *int     `json:"in_stock"`
	ShelfLayerID *int64   `json:"shelf_layer_id"`
	Price        *float64 `json:"price"`
	Remark       *string  `json:"remark"`
}

// UpdateBookRequest 更新图书请求
type UpdateBookRequest struct {
	Name         *string  `json:"name"`
	Quantity     *int     `json:"quantity"`
	InStock      *int     `json:"in_stock"`
	ShelfLayerID *int64   `json:"shelf_layer_id"`
	Price        *float64 `json:"price"`
	Remark       *string  `json:"remark"`
}

// Create 创建图书
func (h *BookHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateBookRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 如果未指定在库数量，默认等于数量
	inStock := req.Quantity
	if req.InStock != nil {
		inStock = *req.InStock
	}

	// 验证在库数量不能大于总数量
	if inStock > req.Quantity {
		Error(c, 400, "在库数量不能大于总数量")
		return
	}

	book := db.Book{
		Barcode:      req.Barcode,
		Name:         req.Name,
		Quantity:     req.Quantity,
		InStock:      inStock,
		ShelfLayerID: req.ShelfLayerID,
		Price:        req.Price,
		Remark:       req.Remark,
	}

	if err := h.db.Create(&book).Error; err != nil {
		Error(c, 400, "创建失败: "+err.Error())
		return
	}

	Success(c, book)
}

// List 查询图书列表
func (h *BookHandler) List(ctx context.Context, c *app.RequestContext) {
	var books []db.Book
	query := h.db.Model(&db.Book{})

	// 书名模糊匹配
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 一维码精确匹配
	if barcode := c.Query("barcode"); barcode != "" {
		query = query.Where("barcode = ?", barcode)
	}

	// 位置查询
	if areaID := c.Query("area_id"); areaID != "" {
		query = query.Joins("JOIN shelf_layer ON book.shelf_layer_id = shelf_layer.id").
			Joins("JOIN bookshelf ON shelf_layer.bookshelf_id = bookshelf.id").
			Where("bookshelf.area_id = ?", areaID)
	}
	if bookshelfID := c.Query("bookshelf_id"); bookshelfID != "" {
		query = query.Joins("JOIN shelf_layer ON book.shelf_layer_id = shelf_layer.id").
			Where("shelf_layer.bookshelf_id = ?", bookshelfID)
	}
	if shelfLayerID := c.Query("shelf_layer_id"); shelfLayerID != "" {
		query = query.Where("shelf_layer_id = ?", shelfLayerID)
	}

	// 在库状态
	if status := c.Query("in_stock_status"); status != "" {
		if status == "1" {
			query = query.Where("in_stock > 0")
		} else if status == "2" {
			query = query.Where("in_stock < quantity")
		}
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Preload("ShelfLayer.Bookshelf.Area").Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}

	// 构建响应数据
	type BookResponse struct {
		ID             int64   `json:"id"`
		Barcode        string  `json:"barcode"`
		Name           string  `json:"name"`
		Quantity       int     `json:"quantity"`
		InStock        int     `json:"in_stock"`
		ShelfLayerID   *int64  `json:"shelf_layer_id"`
		ShelfLayerName *string `json:"shelf_layer_name"`
		Price          *float64 `json:"price"`
		Remark         *string `json:"remark"`
	}

	list := make([]BookResponse, len(books))
	for i, book := range books {
		var shelfLayerName *string
		if book.ShelfLayerID != nil && book.ShelfLayer.ID != 0 {
			name := book.ShelfLayer.Bookshelf.Area.Name + "-" + book.ShelfLayer.Bookshelf.Name + "-" + book.ShelfLayer.Name
			shelfLayerName = &name
		}
		list[i] = BookResponse{
			ID:             book.ID,
			Barcode:        book.Barcode,
			Name:           book.Name,
			Quantity:       book.Quantity,
			InStock:        book.InStock,
			ShelfLayerID:   book.ShelfLayerID,
			ShelfLayerName: shelfLayerName,
			Price:          book.Price,
			Remark:         book.Remark,
		}
	}

	Success(c, map[string]interface{}{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetByBarcode 根据一维码查询图书
func (h *BookHandler) GetByBarcode(ctx context.Context, c *app.RequestContext) {
	barcode := c.Param("barcode")
	var book db.Book
	if err := h.db.Where("barcode = ?", barcode).Preload("ShelfLayer.Bookshelf.Area").First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, 404, "图书不存在")
		} else {
			Error(c, 500, "查询失败: "+err.Error())
		}
		return
	}
	Success(c, book)
}

// Update 更新图书信息
func (h *BookHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	var req UpdateBookRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	var book db.Book
	if err := h.db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, 404, "图书不存在")
		} else {
			Error(c, 500, "查询失败: "+err.Error())
		}
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.InStock != nil {
		updates["in_stock"] = *req.InStock
	}
	if req.ShelfLayerID != nil {
		updates["shelf_layer_id"] = *req.ShelfLayerID
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	// 验证在库数量不能大于总数量
	if req.InStock != nil && req.Quantity != nil && *req.InStock > *req.Quantity {
		Error(c, 400, "在库数量不能大于总数量")
		return
	}
	if req.InStock != nil && req.Quantity == nil && *req.InStock > book.Quantity {
		Error(c, 400, "在库数量不能大于总数量")
		return
	}

	if err := h.db.Model(&book).Updates(updates).Error; err != nil {
		Error(c, 500, "更新失败: "+err.Error())
		return
	}

	h.db.First(&book, id)
	Success(c, book)
}

// Delete 删除图书
func (h *BookHandler) Delete(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Error(c, 400, "无效的ID")
		return
	}

	// 检查是否存在借阅记录
	var count int64
	h.db.Model(&db.BorrowDetail{}).Where("book_id = ?", id).Count(&count)
	if count > 0 {
		Error(c, 400, "该图书存在借阅记录，无法删除")
		return
	}

	if err := h.db.Delete(&db.Book{}, id).Error; err != nil {
		Error(c, 500, "删除失败: "+err.Error())
		return
	}

	Success(c, nil)
}



