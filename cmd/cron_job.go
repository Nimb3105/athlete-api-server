package main

import (
	"be/internal/services"
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// InitCronJobs khởi tạo và chạy các cron job định kỳ
func InitCronJobs(scheduleService *services.TrainingScheduleService) {
	c := cron.New(cron.WithSeconds()) // dùng WithSeconds nếu muốn cron chính xác tới giây

	_, err := c.AddFunc("@every 2m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		updatedCount, updatedExercises, err := scheduleService.AutoMarkOverdue(ctx)
		if err != nil {
			log.Printf("[CRON] Lỗi cập nhật trạng thái quá hạn: %v", err)
			return
		}

		if updatedCount > 0 {
			log.Printf("[CRON] Đã cập nhật %d lịch tập thành 'chưa hoàn thành'", updatedCount)
		}

		if updatedExercises > 0 {
			log.Printf("[CRON] Đã cập nhật %d bài tập thành 'chưa hoàn thành'", updatedExercises)
		}
	})
	if err != nil {
		log.Fatalf("[CRON] Không thể tạo cron job: %v", err)
	}

	c.Start()
	log.Println("[CRON] Đã khởi chạy cron job.")
}
