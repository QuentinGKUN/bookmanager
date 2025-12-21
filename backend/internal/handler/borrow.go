package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"booksystem/internal/db"
)

type BorrowHandler struct {
	db *gorm.DB
}

func NewBorrowHandler(db *gorm.DB) *BorrowHandler {
	return &BorrowHandler{db: db}
}

// CreateBorrowRequest 创建借阅记录请求
type CreateBorrowRequest struct {
	BorrowerName  string   `json:"borrower_name" binding:"required"`
	BorrowerPhone string   `json:"borrower_phone" binding:"required"`
	Barcodes      []string `json:"barcodes" binding:"required"`
}

// ScanBorrowRequest 扫码借阅请求
type ScanBorrowRequest struct {
	Barcode string `json:"barcode" binding:"required"`
}

// ReturnBorrowRequest 归还请求
type ReturnBorrowRequest struct {
	Barcode       string `json:"barcode" binding:"required"`
	BorrowerPhone string `json:"borrower_phone" binding:"required"`
}

// GetBorrowerByPhoneRequest 根据电话查询归还人请求
type GetBorrowerByPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// Create 创建借阅记录
func (h *BorrowHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateBorrowRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建借阅记录
	record := db.BorrowRecord{
		BorrowerName:  req.BorrowerName,
		BorrowerPhone: req.BorrowerPhone,
		BorrowTime:    time.Now(),
		Status:        1,
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		Error(c, 500, "创建借阅记录失败: "+err.Error())
		return
	}

	// 创建借阅明细并更新图书在库数量
	details := make([]db.BorrowDetail, 0, len(req.Barcodes))
	for _, barcode := range req.Barcodes {
		// 查找图书（可能不存在）
		var book db.Book
		tx.Where("barcode = ?", barcode).First(&book)

		// 创建借阅明细
		detail := db.BorrowDetail{
			BorrowRecordID: record.ID,
			Barcode:        barcode,
		}
		if book.ID != 0 {
			detail.BookID = &book.ID
			// 更新图书在库数量
			if book.InStock > 0 {
				tx.Model(&book).Update("in_stock", gorm.Expr("in_stock - 1"))
			}
		}
		details = append(details, detail)
	}

	if err := tx.Create(&details).Error; err != nil {
		tx.Rollback()
		Error(c, 500, "创建借阅明细失败: "+err.Error())
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		Error(c, 500, "提交失败: "+err.Error())
		return
	}

	// 查询完整的借阅记录
	h.db.Preload("Details.Book").First(&record, record.ID)

	// 构建响应
	type BookInfo struct {
		ID      *int64  `json:"id,omitempty"`
		Barcode string  `json:"barcode"`
		Name    *string `json:"name,omitempty"`
	}

	books := make([]BookInfo, len(details))
	for i, detail := range details {
		books[i] = BookInfo{
			ID:      detail.BookID,
			Barcode: detail.Barcode,
		}
		if detail.Book != nil {
			books[i].Name = &detail.Book.Name
		}
	}

	Success(c, map[string]interface{}{
		"id":            record.ID,
		"borrower_name": record.BorrowerName,
		"borrower_phone": record.BorrowerPhone,
		"borrow_time":   record.BorrowTime,
		"books":         books,
	})
}

// Scan 扫码借阅（单本，用于前端实时扫码）
func (h *BorrowHandler) Scan(ctx context.Context, c *app.RequestContext) {
	var req ScanBorrowRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 查找图书（可能不存在）
	var book db.Book
	h.db.Where("barcode = ?", req.Barcode).First(&book)

	// 如果图书存在，更新在库数量
	if book.ID != 0 && book.InStock > 0 {
		h.db.Model(&book).Update("in_stock", gorm.Expr("in_stock - 1"))
		h.db.First(&book, book.ID)
	}

	// 返回图书信息（即使不存在也返回一维码）
	result := map[string]interface{}{
		"barcode": req.Barcode,
	}

	if book.ID != 0 {
		result["id"] = book.ID
		result["name"] = book.Name
		result["in_stock"] = book.InStock
	}

	Success(c, result)
}

// List 查询借阅记录列表
func (h *BorrowHandler) List(ctx context.Context, c *app.RequestContext) {
	var records []db.BorrowRecord
	query := h.db.Model(&db.BorrowRecord{})

	// 借阅人姓名（精确匹配）
	if name := c.Query("borrower_name"); name != "" {
		query = query.Where("borrower_name = ?", name)
	}

	// 借阅人电话（精确匹配）
	if phone := c.Query("borrower_phone"); phone != "" {
		query = query.Where("borrower_phone = ?", phone)
	}

	// 图书一维码（精确匹配）
	if barcode := c.Query("barcode"); barcode != "" {
		query = query.Joins("JOIN borrow_detail ON borrow_record.id = borrow_detail.borrow_record_id").
			Where("borrow_detail.barcode = ?", barcode)
	}

	// 时间范围
	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("borrow_time >= ?", startTime)
	}
	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("borrow_time <= ?", endTime)
	}

	// 只查询正在借阅的记录（根据需求文档2.3.2）
	if c.Query("status") == "" {
		query = query.Where("status = 1")
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Preload("Details.Book").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}

	// 构建响应数据
	type BookInfo struct {
		Barcode string  `json:"barcode"`
		Name    *string `json:"name,omitempty"`
	}

	type RecordResponse struct {
		ID            int64      `json:"id"`
		BorrowerName string     `json:"borrower_name"`
		BorrowerPhone string    `json:"borrower_phone"`
		BorrowTime   time.Time  `json:"borrow_time"`
		Status        int8       `json:"status"`
		Books         []BookInfo `json:"books"`
	}

	list := make([]RecordResponse, len(records))
	for i, record := range records {
		books := make([]BookInfo, len(record.Details))
		for j, detail := range record.Details {
			books[j] = BookInfo{
				Barcode: detail.Barcode,
			}
			if detail.Book != nil {
				books[j].Name = &detail.Book.Name
			}
		}
		list[i] = RecordResponse{
			ID:            record.ID,
			BorrowerName:  record.BorrowerName,
			BorrowerPhone: record.BorrowerPhone,
			BorrowTime:    record.BorrowTime,
			Status:        record.Status,
			Books:         books,
		}
	}

	Success(c, map[string]interface{}{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetBorrowerByPhone 根据电话查询归还人信息
func (h *BorrowHandler) GetBorrowerByPhone(ctx context.Context, c *app.RequestContext) {
	var req GetBorrowerByPhoneRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 查询该电话的借阅记录（状态为借出）
	var records []db.BorrowRecord
	if err := h.db.Where("borrower_phone = ? AND status = 1", req.Phone).
		Preload("Details.Book").
		Find(&records).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}

	if len(records) == 0 {
		Error(c, 404, "未找到该电话的借阅记录")
		return
	}

	// 取第一条记录作为归还人信息（电话唯一确定归还人）
	// 合并所有记录中的图书明细
	record := records[0]
	type BookInfo struct {
		Barcode string  `json:"barcode"`
		Name    *string `json:"name,omitempty"`
	}

	books := make([]BookInfo, 0)
	for _, r := range records {
		for _, detail := range r.Details {
			bookInfo := BookInfo{
				Barcode: detail.Barcode,
			}
			if detail.Book != nil {
				bookInfo.Name = &detail.Book.Name
			}
			books = append(books, bookInfo)
		}
	}

	Success(c, map[string]interface{}{
		"borrower_name":  record.BorrowerName,
		"borrower_phone": record.BorrowerPhone,
		"borrow_time":    record.BorrowTime,
		"books":          books,
	})
}

// Return 归还图书
func (h *BorrowHandler) Return(ctx context.Context, c *app.RequestContext) {
	var req ReturnBorrowRequest
	if err := c.BindAndValidate(&req); err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 查找该一维码的借阅明细（状态为借出，且归还人电话匹配）
	var detail db.BorrowDetail
	if err := h.db.Joins("JOIN borrow_record ON borrow_detail.borrow_record_id = borrow_record.id").
		Where("borrow_detail.barcode = ? AND borrow_record.status = 1 AND borrow_record.borrower_phone = ?", req.Barcode, req.BorrowerPhone).
		First(&detail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Error(c, 404, "未找到该图书的借阅记录或归还人信息不匹配")
		} else {
			Error(c, 500, "查询失败: "+err.Error())
		}
		return
	}

	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新图书在库数量（如果图书存在）
	if detail.BookID != nil {
		var book db.Book
		if err := tx.First(&book, *detail.BookID).Error; err == nil {
			tx.Model(&book).Update("in_stock", gorm.Expr("in_stock + 1"))
		}
	}

	// 更新借阅记录状态为已归还
	var record db.BorrowRecord
	if err := tx.First(&record, detail.BorrowRecordID).Error; err != nil {
		tx.Rollback()
		Error(c, 500, "查询借阅记录失败: "+err.Error())
		return
	}

	// 删除该借阅明细（表示已归还）
	if err := tx.Delete(&detail).Error; err != nil {
		tx.Rollback()
		Error(c, 500, "删除借阅明细失败: "+err.Error())
		return
	}

	// 检查该借阅记录是否还有其他未归还的图书
	var remainingCount int64
	tx.Model(&db.BorrowDetail{}).Where("borrow_record_id = ?", detail.BorrowRecordID).Count(&remainingCount)

	// 如果没有其他未归还的图书，更新借阅记录状态为已归还
	if remainingCount == 0 {
		tx.Model(&record).Update("status", 2)
	}

	if err := tx.Commit().Error; err != nil {
		Error(c, 500, "提交失败: "+err.Error())
		return
	}

	Success(c, map[string]interface{}{
		"message": "归还成功",
	})
}

