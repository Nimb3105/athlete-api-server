package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ForeignKeyCheckConfig định nghĩa cấu hình cho một collection cần kiểm tra
type ForeignKeyCheckConfig struct {
	Collection *mongo.Collection // Collection để kiểm tra
	Filter     bson.M            // Bộ lọc để tìm bản ghi liên quan
	Name       string            // Tên tiếng Việt của collection (dùng cho thông báo lỗi)
}

// CheckForeignKeyConstraints kiểm tra các ràng buộc khóa ngoại
// Trả về lỗi bằng tiếng Việt nếu có bản ghi liên quan
func CheckForeignKeyConstraints(ctx context.Context, configs []ForeignKeyCheckConfig) error {
	for _, config := range configs {
		count, err := config.Collection.CountDocuments(ctx, config.Filter)
		if err != nil {
			return fmt.Errorf("lỗi khi kiểm tra %s: %w", config.Name, err)
		}
		if count > 0 {
			return fmt.Errorf("không thể xóa vì còn liên quan đến %s", config.Name)
		}
	}
	return nil
}