package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis键前缀
	BorrowUserKey  = "borrow:user"  // 当前借阅用户信息
	BorrowBooksKey = "borrow:books" // 当前借阅的图书列表
	ReturnUserKey  = "return:user"  // 当前归还用户信息
	ReturnBooksKey = "return:books" // 当前归还的图书列表
	RedisTTL       = 2 * time.Hour  // Redis数据过期时间（2小时）
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(addr, password string, db int) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisService{client: client}, nil
}

// BorrowUser 借阅用户信息
type BorrowUser struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// BorrowBook 借阅图书信息
type BorrowBook struct {
	Barcode string  `json:"barcode"`
	Name    *string `json:"name,omitempty"`
}

// SetBorrowUser 设置当前借阅用户
func (r *RedisService) SetBorrowUser(ctx context.Context, user *BorrowUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, BorrowUserKey, data, RedisTTL).Err()
}

// GetBorrowUser 获取当前借阅用户
func (r *RedisService) GetBorrowUser(ctx context.Context) (*BorrowUser, error) {
	data, err := r.client.Get(ctx, BorrowUserKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var user BorrowUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// AddBorrowBook 添加借阅图书
func (r *RedisService) AddBorrowBook(ctx context.Context, book *BorrowBook) error {
	data, err := json.Marshal(book)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, BorrowBooksKey, data).Err()
}

// GetBorrowBooks 获取当前借阅图书列表
func (r *RedisService) GetBorrowBooks(ctx context.Context) ([]*BorrowBook, error) {
	dataList, err := r.client.LRange(ctx, BorrowBooksKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	books := make([]*BorrowBook, 0, len(dataList))
	for _, data := range dataList {
		var book BorrowBook
		if err := json.Unmarshal([]byte(data), &book); err != nil {
			continue
		}
		books = append(books, &book)
	}
	return books, nil
}

// RemoveBorrowBook 删除借阅图书（通过索引）
func (r *RedisService) RemoveBorrowBook(ctx context.Context, index int) error {
	// 先获取所有图书
	books, err := r.GetBorrowBooks(ctx)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(books) {
		return fmt.Errorf("index out of range")
	}

	// 删除列表并重新添加（除了要删除的那本）
	if err := r.client.Del(ctx, BorrowBooksKey).Err(); err != nil {
		return err
	}

	// 反向添加（保持顺序）
	for i := len(books) - 1; i >= 0; i-- {
		if i != index {
			data, _ := json.Marshal(books[i])
			if err := r.client.LPush(ctx, BorrowBooksKey, data).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearBorrowData 清除借阅数据
func (r *RedisService) ClearBorrowData(ctx context.Context) error {
	return r.client.Del(ctx, BorrowUserKey, BorrowBooksKey).Err()
}

// SetReturnUser 设置当前归还用户
func (r *RedisService) SetReturnUser(ctx context.Context, user *BorrowUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, ReturnUserKey, data, RedisTTL).Err()
}

// GetReturnUser 获取当前归还用户
func (r *RedisService) GetReturnUser(ctx context.Context) (*BorrowUser, error) {
	data, err := r.client.Get(ctx, ReturnUserKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var user BorrowUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// AddReturnBook 添加归还图书
func (r *RedisService) AddReturnBook(ctx context.Context, book *BorrowBook) error {
	data, err := json.Marshal(book)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, ReturnBooksKey, data).Err()
}

// GetReturnBooks 获取当前归还图书列表
func (r *RedisService) GetReturnBooks(ctx context.Context) ([]*BorrowBook, error) {
	dataList, err := r.client.LRange(ctx, ReturnBooksKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	books := make([]*BorrowBook, 0, len(dataList))
	for _, data := range dataList {
		var book BorrowBook
		if err := json.Unmarshal([]byte(data), &book); err != nil {
			continue
		}
		books = append(books, &book)
	}
	return books, nil
}

// RemoveReturnBook 删除归还图书（通过索引）
func (r *RedisService) RemoveReturnBook(ctx context.Context, index int) error {
	// 先获取所有图书
	books, err := r.GetReturnBooks(ctx)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(books) {
		return fmt.Errorf("index out of range")
	}

	// 删除列表并重新添加（除了要删除的那本）
	if err := r.client.Del(ctx, ReturnBooksKey).Err(); err != nil {
		return err
	}

	// 反向添加（保持顺序）
	for i := len(books) - 1; i >= 0; i-- {
		if i != index {
			data, _ := json.Marshal(books[i])
			if err := r.client.LPush(ctx, ReturnBooksKey, data).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearReturnData 清除归还数据
func (r *RedisService) ClearReturnData(ctx context.Context) error {
	return r.client.Del(ctx, ReturnUserKey, ReturnBooksKey).Err()
}
