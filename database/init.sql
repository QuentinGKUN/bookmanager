-- 图书管理系统数据库初始化脚本
-- MySQL 8.4.5

CREATE DATABASE IF NOT EXISTS booksystem DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE booksystem;

-- 区域表
CREATE TABLE IF NOT EXISTS `area` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '区域名称',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='区域表';

-- 书架表
CREATE TABLE IF NOT EXISTS `bookshelf` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `area_id` BIGINT NOT NULL COMMENT '所属区域ID',
  `name` VARCHAR(50) NOT NULL COMMENT '书架名称',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  UNIQUE KEY `uk_area_name` (`area_id`, `name`),
  FOREIGN KEY (`area_id`) REFERENCES `area`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书架表';

-- 层数表
CREATE TABLE IF NOT EXISTS `shelf_layer` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `bookshelf_id` BIGINT NOT NULL COMMENT '所属书架ID',
  `name` VARCHAR(50) NOT NULL COMMENT '层数名称',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  UNIQUE KEY `uk_bookshelf_name` (`bookshelf_id`, `name`),
  FOREIGN KEY (`bookshelf_id`) REFERENCES `bookshelf`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='层数表';

-- 图书表
CREATE TABLE IF NOT EXISTS `book` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `barcode` VARCHAR(100) NOT NULL UNIQUE COMMENT '一维码（唯一标识）',
  `name` VARCHAR(200) NOT NULL COMMENT '书名',
  `quantity` INT NOT NULL DEFAULT 0 COMMENT '总数量',
  `in_stock` INT NOT NULL DEFAULT 0 COMMENT '在库数量',
  `shelf_layer_id` BIGINT NULL COMMENT '位置ID（层数ID）',
  `price` DECIMAL(10,2) NULL COMMENT '价格',
  `remark` TEXT NULL COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_shelf_layer_id` (`shelf_layer_id`),
  FOREIGN KEY (`shelf_layer_id`) REFERENCES `shelf_layer`(`id`) ON DELETE SET NULL,
  CHECK (`in_stock` <= `quantity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图书表';

-- 借阅记录表
CREATE TABLE IF NOT EXISTS `borrow_record` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `borrower_name` VARCHAR(50) NOT NULL COMMENT '借阅人姓名',
  `borrower_phone` VARCHAR(20) NOT NULL COMMENT '借阅人电话',
  `borrow_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '借阅时间',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（1:借出，2:已归还）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_borrower_phone` (`borrower_phone`),
  INDEX `idx_borrow_time` (`borrow_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='借阅记录表';

-- 借阅明细表
CREATE TABLE IF NOT EXISTS `borrow_detail` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `borrow_record_id` BIGINT NOT NULL COMMENT '借阅记录ID',
  `book_id` BIGINT NULL COMMENT '图书ID（可能为空，因为图书可能未录入）',
  `barcode` VARCHAR(100) NOT NULL COMMENT '图书一维码（冗余字段，便于查询）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_borrow_record_id` (`borrow_record_id`),
  INDEX `idx_barcode` (`barcode`),
  FOREIGN KEY (`borrow_record_id`) REFERENCES `borrow_record`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`book_id`) REFERENCES `book`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='借阅明细表';




