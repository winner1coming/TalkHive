package workSpace

import (
	"TalkHive/global"
	"TalkHive/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 我的笔记！！！

// CreateNote - 创建笔记√
func CreateNote(c *gin.Context) {

	log.Printf("Received note_name:")

	// 2. 接收表单其他参数
	var temp struct {
		NoteName string `json:"note_name" binding:"required"`
		Type     string `json:"type" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data", "error": err.Error()})
		return
	}
	log.Printf("Received note_name: %s, type: %s", temp.NoteName, temp.Type)

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3. 检查该用户的分类是否存在
	if temp.Type != "" {
		var noteDivide models.NoteDivide
		if err := global.Db.Model(&models.NoteDivide{}).Where("account_id = ? AND nd_name = ?",
			global.ParseUint(userID),
			temp.Type).First(&noteDivide).Error; err != nil {
			// 如果没有找到分类，则插入新的分类
			if err := global.Db.Create(&models.NoteDivide{
				NDName:    temp.Type,
				AccountID: global.ParseUint(userID),
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create note category",
					"error": err.Error()})
				return
			}
		}
	}

	// 4.确保文件路径安全并添加后缀
	sanitizedNoteName := strings.ReplaceAll(temp.NoteName, "/", "_")
	sanitizedNoteName = strings.ReplaceAll(sanitizedNoteName, "\\", "_")
	rootDir := "D:/TalkHive/Notes/"
	filePath := filepath.Join(rootDir, sanitizedNoteName+".md")

	// 确保目录存在
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create directory", "error": err.Error()})
		return
	}

	// 5.创建空白文件
	file, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create blank note file", "error": err.Error()})
		return
	}
	defer file.Close()

	// 6.保存到数据库
	note := models.Notes{
		NoteName:  temp.NoteName,
		SaveTime:  time.Now(),
		CachePath: sanitizedNoteName + ".md", //目标路径下的相对路径
		Type:      temp.Type,
		AccountID: global.ParseUint(userID),
		IsShow:    true,
	}
	if err := global.Db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save note metadata", "error": err.Error()})
		return
	}

	// 7.返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"note_id":   note.NoteID,
		"note_name": note.NoteName,
		"save_time": note.SaveTime.Format("2006-01-02 15:04"),
		"type":      note.Type,
	})
}

// GetNotesList - 获取用户的笔记列表√
func GetNotesList(c *gin.Context) {
	// 1. 从请求中获取 user_id 参数
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	// 定义用于存储结果的结构体
	type NoteResponse struct {
		NoteID           uint      `json:"note_id"`
		NoteName         string    `json:"note_name"`
		LastModifiedTime time.Time `json:"save_time"` // 使用 time.Time 类型
		Type             string    `json:"type"`
	}

	// 查询数据库
	var notes []NoteResponse
	if err := global.Db.Table("notes").
		Select("note_id, note_name, save_time as last_modified_time, type").
		Where("account_id = ? AND is_show = ?", userID, true).
		Order("save_time DESC"). // 按时间降序排序
		Scan(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	// 返回结果时格式化时间
	var response []map[string]interface{}
	for _, note := range notes {
		response = append(response, map[string]interface{}{
			"id":           note.NoteID,
			"filename":     note.NoteName,
			"lastModified": note.LastModifiedTime.Format("2006-01-02 15:04"), // 格式化时间
			"category":     note.Type,
		})
	}

	// 返回结果
	c.JSON(http.StatusOK, response)
}

// GetNote - 获取笔记
func GetNote(c *gin.Context) {
	// 1. 从请求中获取 code_id 参数
	//userID := c.Param("id") // 使用 c.Param 获取路径参数
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	var req struct {
		NoteID uint `json:"note_id" binding:"required"` // 修改为正确的 JSON 字段名
	}

	// 2. 绑定请求体中的数据到 req 结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// 3. 数据库查询：查找指定 NoteID 且 IsShow = true 的记录
	var note models.Notes
	if err := global.Db.Model(&models.Notes{}).Where("note_id = ? AND account_id = ? AND is_show = ?",
		req.NoteID, global.ParseUint(userID), true).First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Note not found or not visible"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error", "error": err.Error()})
		return
	}

	// 4. 拼接文件的完整路径
	rootDir := "D:/TalkHive/Notes/" // 默认根目录
	filePath := filepath.Join(rootDir, note.CachePath)

	// 5. 验证文件是否存在
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open file", "error": err.Error()})
		return
	}
	defer file.Close()

	// 6. 设置正确的 HTTP 响应头
	contentType := "text/markdown" // 对应md格式的文件
	c.Header("Content-Type", contentType)

	// 7. 返回文件流
	if _, err := io.Copy(c.Writer, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send file", "error": err.Error()})
		return
	}
}

// EditNote - 编辑笔记
func EditNote(c *gin.Context) {
	// 接收JSON数据
	var requestData struct {
		NoteID   uint   `json:"NoteID" binding:"required"`   // 笔记ID
		NoteName string `json:"NoteName" binding:"required"` // 笔记名称
		Type     string `json:"Type" binding:"required"`     // 笔记类型
		Content  string `json:"Content" binding:"required"`  // 笔记内容
	}

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 获取用户ID
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 检查数据库中是否有匹配的记录
	var note models.Notes
	if err := global.Db.Model(&models.Notes{}).Where("note_id = ? AND account_id = ? AND is_show = ?",
		requestData.NoteID, global.ParseUint(userID), true).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// 更新笔记名称为 note_name
	if err := global.Db.Model(&note).Update("note_name", requestData.NoteName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note name"})
		return
	}

	// 2. 更新 Notes 表，将 Type 修改为 Type
	if err := global.Db.Model(&note).Update("type", requestData.Type).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note type"})
		return
	}

	// 更新笔记内容到指定路径
	savePath := fmt.Sprintf("D:/TalkHive/Notes/%s", note.CachePath)
	err := os.WriteFile(savePath, []byte(requestData.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save note content"})
		return
	}

	// 更新数据库的保存时间
	note.SaveTime = time.Now()
	if err := global.Db.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note metadata"})
		return
	}

	// 返回元信息给前端
	c.JSON(http.StatusOK, gin.H{
		"note_id":   note.NoteID,
		"note_name": requestData.NoteName,
		"save_time": note.SaveTime.Format("2006-01-02 15:04"),
		"type":      requestData.Type,
	})
}

// ShareNote - 将笔记分享给通讯录里的好友
func ShareNote(c *gin.Context) {
	// 解析前端传递的参数
	var request struct {
		NoteID uint `json:"note_id" binding:"required"`
		FdID   uint `json:"fd_id" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1. 检查是否存在这个好友
	var contact models.Contacts
	if err := global.Db.Model(&models.Contacts{}).Where("owner_id = ? AND contact_id = ?",
		global.ParseUint(userID), request.FdID).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	// 2. 检查用户是否拥有该笔记
	var note models.Notes
	if err := global.Db.Model(&models.Notes{}).Where("account_id = ? AND note_id = ?",
		global.ParseUint(userID), request.NoteID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not accessible"})
		return
	}

	// 3. 在 Notes 表中为好友新增一条笔记记录
	newNote := models.Notes{
		NoteName:  note.NoteName,
		SaveTime:  time.Now(),
		Type:      note.Type,
		CachePath: note.CachePath,
		AccountID: request.FdID, // 分享给好友
		IsShow:    note.IsShow,  // 是否显示，通常是共享的状态
	}

	if err := global.Db.Create(&newNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share note"})
		return
	}

	// 4. 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Note shared successfully"})
}

// ChangeNoteName - 修改笔记名称
func ChangeNoteName(c *gin.Context) {
	// 解析前端传递的参数
	var request struct {
		NoteID      uint   `json:"note_id" binding:"required"`
		OldNoteName string `json:"old_note_name" binding:"required"`
		NewNoteName string `json:"new_note_name" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1. 从 Notes 表中筛选出对应的笔记项
	var note models.Notes
	if err := global.Db.Model(&models.Notes{}).Where("note_id = ? AND is_show = true AND account_id = ?",
		request.NoteID, global.ParseUint(userID)).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not accessible"})
		return
	}

	// 2. 判断 Name 是否为 old_note_name
	if note.NoteName != request.OldNoteName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old note name does not match"})
		return
	}

	// 3. 更新笔记名称为 new_note_name
	if err := global.Db.Model(&models.Notes{}).Where("note_id = ? AND account_id = ?",
		request.NoteID, global.ParseUint(userID)).Update("name", request.NewNoteName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note name"})
		return
	}

	// 4. 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Note name updated successfully"})
}

// DeleteNote - 删除笔记√
func DeleteNote(c *gin.Context) {
	var req struct {
		NoteID uint `json:"note_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind Error:", err) // 打印绑定错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 打印调试信息
	fmt.Printf("Parsed Request: userID = %s, noteID = %d\n", userID, req.NoteID)

	// 1. 检查 Notes 表中是否存在指定笔记（account_id 和 note_id 匹配）
	var note models.Notes
	if err := global.Db.Model(&models.Notes{}).Where("account_id = ? AND note_id = ?",
		global.ParseUint(userID), req.NoteID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// 2. 更新 is_show 字段为 false（表示移动到回收站）
	note.IsShow = false
	if err := global.Db.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note status"})
		return
	}

	// 3. 向 Recycle 表插入回收站记录
	recycle := models.Recycle{
		RecycleID:   note.NoteID,
		RecycleType: "note",
		AccountID:   note.AccountID,
		RecycleTime: time.Now(),
	}
	if err := global.Db.Create(&recycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add note to recycle bin"})
		return
	}

	// 4. 返回操作成功信息
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted and moved to recycle bin successfully"})
}

// GetTypeList - 获取用户的笔记分类列表√
func GetTypeList(c *gin.Context) {
	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "User ID is required",
		})
		return
	}

	// 查询 NoteDivide 表
	var noteDivide []models.NoteDivide
	if err := global.Db.Where("account_id = ?", global.ParseUint(userID)).Find(&noteDivide).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Failed to fetch note categories",
		})
		return
	}

	// 提取分类名称列表
	var categories []string
	for _, category := range noteDivide {
		categories = append(categories, category.NDName)
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, gin.H{
		"status":     200,
		"message":    "获取成功",
		"categories": categories,
	})
}

// GetNotesByCategory - 按分类查看笔记√
func GetNotesByCategory(c *gin.Context) {
	// 获取路径参数
	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var seq struct {
		Type string `json:"type" binding:"required"`
	}

	// 绑定请求体中的JSON数据
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 查询 Notes 表
	var notes []models.Notes
	query := global.Db.Model(&models.Notes{}).Where("account_id = ? AND is_show = ?",
		global.ParseUint(userID), true)
	if seq.Type != "" && seq.Type != "null" { // 未分类用 "null" 表示
		query = query.Where("type = ?", seq.Type)
	} else {
		query = query.Where("type IS NULL")
	}

	if err := query.Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes by category"})
		return
	}

	// 格式化返回结果
	var result []map[string]interface{}
	for _, note := range notes {
		// 格式化时间字段
		var formattedTime string
		if note.SaveTime.IsZero() {
			formattedTime = "" // 处理零时间，可能没有时间数据时返回空字符串
		} else {
			formattedTime = note.SaveTime.Format("2006-01-02 15:04") // 格式化为年-月-日 时:分
		}

		result = append(result, map[string]interface{}{
			"id":           note.NoteID,
			"filename":     note.NoteName,
			"lastmodified": formattedTime, // 格式化时间字段
			"category":     note.Type,
		})
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, result)
}

// EditNoteType - 修改笔记所在分类
func EditNoteType(c *gin.Context) {
	// 接收请求参数
	var requestData struct {
		NoteID  uint   `json:"note_id" binding:"required"`        // 对应输入字段 "note_id"
		OldName string `json:"old_type_name" binding:"omitempty"` // 对应输入字段 "old_type_name"
		NewName string `json:"new_type_name" binding:"omitempty"` // 对应输入字段 "new_type_name"
	}

	// 尝试绑定 JSON 参数
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 参数验证：OldName 和 NewName 不允许同时为空
	if requestData.NewName == "" && requestData.OldName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both cannot be empty"})
		return
	}

	// 获取用户 ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1. 查询 Notes 表，检查笔记是否存在
	var note models.Notes
	query := global.Db.Model(&models.Notes{}).Where("account_id = ? AND note_id = ?",
		global.ParseUint(userID), requestData.NoteID)

	// 如果 OldName 不为空，则需要验证分类是否匹配
	if requestData.OldName != "" {
		query = query.Where("type = ?", requestData.OldName)
	}

	// 查询数据库
	if err := query.First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or old_type_name mismatch"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
		return
	}

	// 2. 更新 Notes 表，将 Type 修改为 new_type_name
	if err := global.Db.Model(&note).Update("type", requestData.NewName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note type"})
		return
	}

	// 3. 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Note type updated successfully"})
}

// DeleteType - 删除分类√
func DeleteType(c *gin.Context) {
	// 接收请求参数
	var requestData struct {
		NAName string `json:"type_name" binding:"required"`
	}

	// 参数验证
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1.检查 NoteDivide 表是否存在该分类
	var noteDivide models.NoteDivide
	if err := global.Db.Where("nd_name = ? AND account_id = ?", requestData.NAName,
		global.ParseUint(userID)).First(&noteDivide).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
		return
	}

	// 2.更新 Notes 表，将对应分类的 Type 重置为 null
	if err := global.Db.Model(&models.Notes{}).
		Where("type = ? AND account_id = ?", requestData.NAName, global.ParseUint(userID)).
		Update("type", gorm.Expr("null")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset note types"})
		return
	}

	// 3.删除 NoteDivide 表中对应的分类数据项
	if err := global.Db.Delete(&noteDivide).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	// 4.返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted and notes reset to 'uncategorized'"})
}

// ChangeTypeName - 修改笔记分类名称√
func ChangeTypeName(c *gin.Context) {

	var requestData struct {
		OldName string `json:"old_type_name"`
		NewName string `json:"new_type_name"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1.查询 NoteDivide 表，检查是否存在对应的分类
	var noteDivide models.NoteDivide
	if err := global.Db.Model(&models.NoteDivide{}).Where("nd_name = ? AND account_id = ?", requestData.OldName,
		global.ParseUint(userID)).First(&noteDivide).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Old category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
		return
	}

	// 2.检查新分类名是否已存在
	var existingNoteDivide models.NoteDivide
	if err := global.Db.Model(&models.NoteDivide{}).Where("nd_name = ? AND account_id = ?", requestData.NewName,
		global.ParseUint(userID)).First(&existingNoteDivide).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "New category name already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	// 3.更新分类名称
	if err := global.Db.Model(&noteDivide).Update("nd_name", requestData.NewName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category name"})
		return
	}

	// 4.返回成功信息
	c.JSON(http.StatusOK, gin.H{"message": "Category name updated successfully"})
}

// CreateType - 创建新的笔记分类√
func CreateType(c *gin.Context) {
	var requestData struct {
		NDName string `json:"type_name"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing type_name or user_id"})
		return
	}

	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 检查是否已经存在相同的分类名
	var existingNoteDivide models.NoteDivide
	if err := global.Db.Model(&models.NoteDivide{}).Where("nd_name = ? AND account_id = ?", requestData.NDName,
		global.ParseUint(userID)).First(&existingNoteDivide).Error; err == nil {
		// 如果存在分类名，返回错误
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果查询出错且不是记录未找到错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	// 如果不存在，则创建新的分类
	newNoteDivide := models.NoteDivide{
		NDName:    requestData.NDName,
		AccountID: global.ParseUint(userID),
	}
	if err := global.Db.Create(&newNoteDivide).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new category"})
		return
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{"message": "Category created successfully"})
}
