package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"be/internal/models"
	"be/internal/services"

	"github.com/gin-gonic/gin"
)

type ExerciseController struct {
	service *services.ExerciseService
}

func NewExerciseController(service *services.ExerciseService) *ExerciseController {
	return &ExerciseController{service}
}

func (c *ExerciseController) CreateExercise(ctx *gin.Context) {
	var bodyBytes []byte
	if rawData, err := ctx.GetRawData(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu"})
		return
	} else {
		bodyBytes = rawData
	}

	var tempMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không đúng định dạng"})
		return
	}

	validFields := map[string]bool{
		"id": true, "bodyPart": true, "equipment": true, "name": true, "target": true,
		"secondaryMuscles": true, "instructions": true, "gifUrl": true, "createdAt": true, "updatedAt": true, "sportId": true, "unitType": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var exercise models.Exercise
	if err := json.Unmarshal(bodyBytes, &exercise); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdExercise, err := c.service.Create(ctx, &exercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdExercise})
}

func (c *ExerciseController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	exercise, err := c.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "exercise not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": exercise})
}

func (c *ExerciseController) GetAllBySportId(ctx *gin.Context) {
	sportId := ctx.Param("sportId")

	exercises, err := c.service.GetAllBySportId(ctx, sportId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(exercises) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Exercise{},
			"notes":   "không có bài tập nào",
			"message": "chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": exercises,
	})
}

func (c *ExerciseController) GetAllBySportName(ctx *gin.Context) {
	sportName := ctx.Param("sportName")
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	exercises, totalCount, err := c.service.GetAllBySportName(ctx, sportName, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(exercises) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Exercise{},
			"totalCount": 0,
			"notes":      "không có bài tập nào",
			"message":    "chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       exercises,
		"totalCount": totalCount,
	})
}

func (c *ExerciseController) GetAllByBodyPart(ctx *gin.Context) {
	bodyPart := ctx.Param("bodyPart")
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	exercises, totalCount, err := c.service.GetAllByBodyPart(ctx, bodyPart, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(exercises) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Exercise{},
			"totalCount": 0,
			"notes":      "không có bài tập nào",
			"message":    "chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       exercises,
		"totalCount": totalCount,
	})
}

func (c *ExerciseController) GetAll(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	exercises, totalCount, err := c.service.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(exercises) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Exercise{},
			"totalCount": 0,
			"notes":      "không có bài tập nào",
			"message":    "chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       exercises,
		"totalCount": totalCount,
	})
}

func (c *ExerciseController) Update(ctx *gin.Context) {
	ctx.Param("id")
	var exercise models.Exercise
	if err := ctx.ShouldBindJSON(&exercise); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedExercise, err := c.service.Update(ctx, &exercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedExercise})
}

func (c *ExerciseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx, id); err != nil {
		if err.Error() == "exercise not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "exercise deleted"})
}
