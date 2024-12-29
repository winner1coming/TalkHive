package workSpace

import (
	"TalkHive/global"
	"TalkHive/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// DDL记录！！！

// CreateDDL - 新建DDL事件
func CreateDDL(c *gin.Context) {

	// 使用 ShouldBindJSON 将前端传来的参数绑定到结构体 CreateDDLRequest。
	// 添加 binding:"required" 标签，确保必需字段存在，否则返回 400 错误。
	// 定义请求参数结构
	// CreateDDLRequest 结构体定义，增加 `Urgency` 字段的验证
	type CreateDDLRequest struct {
		Deadline    string `json:"deadline" binding:"required"`     // 任务截止日期（字符串，精确到时分）
		TaskContent string `json:"task_content" binding:"required"` // 任务内容
		Urgency     int    `json:"important"`                       // 是否重要，保持一致
	}

	var req CreateDDLRequest
	// 解析并绑定 JSON 请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	// 从 URL 路径中获取用户 ID（:id）
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	location, err := time.LoadLocation("Asia/Shanghai") // 指定中国时区（也可以是其他时区）
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location"})
		return
	}

	deadline, err := time.ParseInLocation("2006-01-02 15:04", req.Deadline, location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// 转换urgency
	urgency := req.Urgency == 1

	// 这里可以补充查询用户id是否正确的逻辑

	// 创建新的 DDL 数据项:DDLID 在 MySQL 里可以直接设置自增
	ddl := models.DDLS{
		AccountID:   global.ParseUint(userID), // 使用从路径中提取的用户ID
		DDLDate:     deadline,
		Task:        req.TaskContent,
		IsCompleted: false,   // 默认未完成
		Urgency:     urgency, // 根据前端请求传递的值
	}

	// 保存到数据库
	if err := global.Db.Create(&ddl).Error; err != nil {
		log.Println("Database error:", err) // 打印数据库错误，便于调试
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DDL"})
		return
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{"message": "DDL created successfully"})
}

// GetUncompletedDDL - 获取用户待完成的DDL列表
func GetUncompletedDDL(c *gin.Context) {
	// 1.获取用户 ID：string
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 2.查询 DDLS 表中待完成的任务
	var ddls []models.DDLS
	if err := global.Db.Where("account_id = ? AND is_completed = ?", global.ParseUint(userID), false).Find(&ddls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DDLs"})
		return
	}

	// 3.格式化返回结果
	var result []map[string]interface{}
	for _, ddl := range ddls {
		result = append(result, map[string]interface{}{
			"task_id":      ddl.DDLID,
			"deadline":     ddl.DDLDate.Format("2006-01-02 15:04"), // 格式化为日期+时间
			"task_content": ddl.Task,
			"important":    ddl.Urgency,
		})
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, result)
}

// GetCompletedDDL - 查看已完成DDL事件
func GetCompletedDDL(c *gin.Context) {
	// 1.获取用户 ID：string
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 2.查询 DDLS 表中已完成的任务
	var ddls []models.DDLS
	if err := global.Db.Where("account_id = ? AND is_completed = ?", global.ParseUint(userID), true).Find(&ddls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DDLs"})
		return
	}

	// 3.格式化返回结果
	var result []map[string]interface{}
	for _, ddl := range ddls {
		result = append(result, map[string]interface{}{
			"task_id":      ddl.DDLID,
			"deadline":     ddl.DDLDate.Format("2006-01-02 15:04"), // 格式化为日期+时间
			"task_content": ddl.Task,
			"important":    ddl.Urgency,
		})
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, result)
}

func EditDDL(c *gin.Context) {
	// 定义输入结构体
	type input struct {
		DDLID   uint   `json:"task_id" binding:"required"` // DDLID
		DDLDate string `json:"ddlDate" binding:"required"` // DDL截止日期
		Task    string `json:"task" binding:"required"`    // DDL任务内容
		Urgency int    `json:"Urgency" binding:"required"` // DDL紧急程度 (0 或 1)
	}

	var req input
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Bind error:", err) // 打印错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	location, err := time.LoadLocation("Asia/Shanghai") // 指定中国时区（也可以是其他时区）
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location"})
		return
	}

	ddlDate, err := time.ParseInLocation("2006-01-02 15:04", req.DDLDate, location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// 转换 Urgency 字段为布尔类型
	urgencyBool := req.Urgency == 1

	// 查询并更新 DDL 数据
	var DDL models.DDLS
	if err := global.Db.Where("account_id = ? AND ddl_id = ?", global.ParseUint(userID), req.DDLID).First(&DDL).Error; err != nil {
		log.Println("Query error:", err) // 打印错误信息
		c.JSON(http.StatusNotFound, gin.H{"error": "DDL not found"})
		return
	}

	// 更新 DDL 数据
	DDL.DDLDate = ddlDate
	DDL.Task = req.Task
	DDL.Urgency = urgencyBool

	if err := global.Db.Save(&DDL).Error; err != nil {
		log.Println("Save error:", err) // 打印错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DDL"})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "DDL updated successfully"})
}

// MarkDDLComplete - 更新DDL为已完成
func MarkDDLComplete(c *gin.Context) {
	// 获取请求参数
	var req struct {
		TaskID uint `json:"task_id" binding:"required"` // DDL ID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 检查是否存在指定的 DDL 数据项
	var ddl models.DDLS
	if err := global.Db.Where("account_id = ? AND ddl_id = ?", userID, req.TaskID).First(&ddl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "DDL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DDL"})
		}
		return
	}

	// 检查是否已经完成
	if ddl.IsCompleted {
		c.JSON(http.StatusOK, gin.H{"message": "This DDL is already marked as completed"})
		return
	}

	// 更新状态为已完成
	if err := global.Db.Model(&ddl).Update("is_completed", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark DDL as completed"})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "DDL marked as completed successfully"})
}

// DeleteDDL - 删除DDL事件
func DeleteDDL(c *gin.Context) {
	// 获取请求参数
	var req struct {
		TaskID uint `json:"task_id" binding:"required"` // DDL ID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 检查是否存在指定的 DDL 数据项
	var ddl models.DDLS
	if err := global.Db.Where("account_id = ? AND ddl_id = ?", userID, req.TaskID).First(&ddl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "DDL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DDL"})
		}
		return
	}

	// 删除指定的 DDL 数据项
	if err := global.Db.Delete(&ddl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DDL"})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "DDL deleted successfully"})
}

// GetUpcomingDDL - 查看主页面DDL提醒
func GetUpcomingDDL(c *gin.Context) {
	var upcomingDDL []models.DDLS
	if err := global.Db.Where("is_complete = ? AND deadline >= ?", false, time.Now()).Order("deadline asc").Limit(5).Find(&upcomingDDL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming DDLs"})
		return
	}
	c.JSON(http.StatusOK, upcomingDDL)
}
