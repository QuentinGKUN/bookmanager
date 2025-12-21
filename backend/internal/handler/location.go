package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"booksystem/internal/db"
)

type LocationHandler struct {
	db *gorm.DB
}

func NewLocationHandler(db *gorm.DB) *LocationHandler {
	return &LocationHandler{db: db}
}

// GetTree 获取位置三级联动树
func (h *LocationHandler) GetTree(ctx context.Context, c *app.RequestContext) {
	var areas []db.Area
	if err := h.db.Find(&areas).Error; err != nil {
		Error(c, 500, "查询失败: "+err.Error())
		return
	}

	type ShelfLayerNode struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	type BookshelfNode struct {
		ID          int64           `json:"id"`
		Name        string          `json:"name"`
		ShelfLayers []ShelfLayerNode `json:"shelf_layers"`
	}

	type AreaNode struct {
		ID          int64            `json:"id"`
		Name        string           `json:"name"`
		Bookshelves []BookshelfNode `json:"bookshelves"`
	}

	result := make([]AreaNode, 0, len(areas))

	for _, area := range areas {
		var bookshelves []db.Bookshelf
		h.db.Where("area_id = ?", area.ID).Find(&bookshelves)

		bookshelfNodes := make([]BookshelfNode, 0, len(bookshelves))
		for _, bookshelf := range bookshelves {
			var shelfLayers []db.ShelfLayer
			h.db.Where("bookshelf_id = ?", bookshelf.ID).Find(&shelfLayers)

			shelfLayerNodes := make([]ShelfLayerNode, 0, len(shelfLayers))
			for _, layer := range shelfLayers {
				shelfLayerNodes = append(shelfLayerNodes, ShelfLayerNode{
					ID:   layer.ID,
					Name: layer.Name,
				})
			}

			bookshelfNodes = append(bookshelfNodes, BookshelfNode{
				ID:          bookshelf.ID,
				Name:        bookshelf.Name,
				ShelfLayers: shelfLayerNodes,
			})
		}

		result = append(result, AreaNode{
			ID:          area.ID,
			Name:        area.Name,
			Bookshelves: bookshelfNodes,
		})
	}

	Success(c, result)
}

