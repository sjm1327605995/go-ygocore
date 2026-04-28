package core

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// CardQueryer 定义了卡片数据查询接口，可由外部实现或使用内置的 sqlx 实现
type CardQueryer interface {
	GetCardData(code uint32) (*CardData, error)
}

// SQLxCardReader 使用 sqlx 查询卡片数据的 CardReader 实现
type SQLxCardReader struct {
	db *sqlx.DB
}

// NewSQLxCardReader 创建基于 sqlx 的卡片数据读取器
// db 由外部提供，sql 语句查询的表结构需包含 CardData 对应的字段
func NewSQLxCardReader(db *sqlx.DB) *SQLxCardReader {
	return &SQLxCardReader{db: db}
}

// DB 返回内部的 *sqlx.DB，便于外部管理连接生命周期
func (r *SQLxCardReader) DB() *sqlx.DB {
	return r.db
}

// GetCardData 根据卡片代码查询卡片数据
func (r *SQLxCardReader) GetCardData(code uint32) (*CardData, error) {
	var card CardData
	query := `SELECT id as code, alias, setcode as set_code, type, level, attribute, race, atk as attack, def as defense, lscale, rscale, link_marker FROM datas WHERE id = ?`
	err := r.db.Get(&card, query, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get card data for code %d: %w", code, err)
	}
	return &card, nil
}

// ToCardReader 将 CardQueryer 适配为 core.CardReader 类型
func ToCardReader(queryer CardQueryer) CardReader {
	return func(code uint32) *CardData {
		card, err := queryer.GetCardData(code)
		if err != nil {
			return nil
		}
		return card
	}
}
