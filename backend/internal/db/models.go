package db

import (
	"time"

	"gorm.io/gorm"
)

// Area 区域表
type Area struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null;unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Bookshelf 书架表
type Bookshelf struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AreaID    int64     `gorm:"not null;index" json:"area_id"`
	Area      Area      `gorm:"foreignKey:AreaID" json:"area,omitempty"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ShelfLayer 层数表
type ShelfLayer struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BookshelfID int64     `gorm:"not null;index" json:"bookshelf_id"`
	Bookshelf   Bookshelf `gorm:"foreignKey:BookshelfID" json:"bookshelf,omitempty"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Book 图书表
type Book struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Barcode      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"barcode"`
	Name         string    `gorm:"type:varchar(200);not null" json:"name"`
	Quantity     int       `gorm:"not null;default:0" json:"quantity"`
	InStock      int       `gorm:"not null;default:0" json:"in_stock"`
	ShelfLayerID *int64    `gorm:"index" json:"shelf_layer_id,omitempty"`
	ShelfLayer   ShelfLayer `gorm:"foreignKey:ShelfLayerID" json:"shelf_layer,omitempty"`
	Price        *float64  `gorm:"type:decimal(10,2)" json:"price,omitempty"`
	Remark       *string   `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BorrowRecord 借阅记录表
type BorrowRecord struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BorrowerName string         `gorm:"type:varchar(50);not null" json:"borrower_name"`
	BorrowerPhone string        `gorm:"type:varchar(20);not null;index" json:"borrower_phone"`
	BorrowTime   time.Time      `gorm:"not null;index" json:"borrow_time"`
	Status       int8           `gorm:"not null;default:1" json:"status"` // 1:借出，2:已归还
	Details      []BorrowDetail `gorm:"foreignKey:BorrowRecordID" json:"details,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// BorrowDetail 借阅明细表
type BorrowDetail struct {
	ID            int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	BorrowRecordID int64      `gorm:"not null;index" json:"borrow_record_id"`
	BookID        *int64       `gorm:"index" json:"book_id,omitempty"`
	Book          *Book        `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Barcode       string       `gorm:"type:varchar(100);not null;index" json:"barcode"`
	CreatedAt     time.Time    `json:"created_at"`
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Area{},
		&Bookshelf{},
		&ShelfLayer{},
		&Book{},
		&BorrowRecord{},
		&BorrowDetail{},
	)
}



