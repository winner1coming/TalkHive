package workSpace

import (
	"TalkHive/global"
	"TalkHive/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 我的笔记！！！

// CreateNote - 创建笔记
func CreateNote(c *gin.Context) {
	// 1. 接收上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No file uploaded", "error": err.Error()})
		return
	}

	// 2. 接收表单其他参数
	var temp struct {
		Name string `form:"note_name" binding:"required"`
		Type string `form:"type" binding:"omitempty"`
	}
	if err := c.ShouldBind(&temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data", "error": err.Error()})
		return
	}

	// 获取用户ID
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3. 检查该用户的分类是否存在
	if temp.Type != "" {
		var noteDivide models.NoteDivide
		if err := global.Db.Where("account_id = ? AND nd_name = ?", global.ParseUint(userID),
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

	// 4. 定义文件保存路径
	rootDir := "D:/TalkHive/Notes/"
	filepath := filepath.Join(rootDir, fmt.Sprintf("%s.md", temp.Name))
	//if _, err := os.Stat(rootDir); os.IsNotExist(err) {
	//	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create directory", "error": err.Error()})
	//		return
	//	}
	//}
	//suffix := temp.Name + ".md"
	//filePath := filepath.Join(rootDir, suffix)
	//
	//// 5. 保存文件到服务器（使用 io.Copy）
	//inFile, err := file.Open()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open uploaded file", "error": err.Error()})
	//	return
	//}
	//defer inFile.Close()
	//
	//outFile, err := os.Create(filePath)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create file on server", "error": err.Error()})
	//	return
	//}
	//defer outFile.Close()
	//
	//if _, err := io.Copy(outFile, inFile); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file", "error": err.Error()})
	//	return
	//}

	suffix := temp.Name + ".md"
	// 5. 保存文件到服务器
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file", "error": err.Error()})
		return
	}

	// 6. 将文件信息存入数据库
	note := models.Notes{
		Name:      temp.Name,
		SaveTime:  time.Now(),
		CachePath: suffix,
		Type:      temp.Type,
		AccountID: global.ParseUint(userID),
		IsShow:    true,
	}
	if err := global.Db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save note metadata", "error": err.Error()})
		return
	}

	// 7. 返回文件元信息
	c.JSON(http.StatusOK, gin.H{
		"note_id":   note.NoteID,
		"note_name": note.Name,
		"save_time": note.SaveTime.Format("2006-01-02 15:04"),
		"type":      note.Type,
	})
}

// GetNotesList - 获取用户的笔记列表
func GetNotesList(c *gin.Context) {
	// 1. 从请求中获取 user_id 参数
	id := c.Param("id") // 使用 c.Param 获取路径参数
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	// 定义用于存储结果的结构体（可以只选择需要返回的字段）
	type NoteResponse struct {
		NoteID           uint   `json:"note_id"`
		NoteName         string `json:"name"`
		LastModifiedTime string `json:"save_time"` // 注意这里使用string类型，防止直接和time.Time冲突
		Type             string `json:"type"`
	}

	// 筛选结果:Scan得到的时间为string类型
	var notes []NoteResponse
	if err := global.Db.Table("notes").
		Select("note_id, name as note_name, save_time as last_modified_time, type, is_show").
		Where("account_id = ? AND is_show = ?", id, true).
		Scan(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	// 格式化时间字段
	for i := range notes {
		// 将LastModifiedTime从string转为time.Time类型
		if parsedTime, err := time.Parse("2006-01-02 15:04:05.000", notes[i].LastModifiedTime); err == nil {
			// 格式化时间
			notes[i].LastModifiedTime = parsedTime.Format("2006-01-02 15:04")
		} else {
			// 如果解析失败，设置默认值或处理错误
			notes[i].LastModifiedTime = "Invalid date format"
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, notes)
}

// GetNote - 获取笔记
func GetNote(c *gin.Context) {
	// 1. 从请求中获取 code_id 参数
	userID := c.Param("id") // 使用 c.Param 获取路径参数
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
	if err := global.Db.Model(&models.Notes{}).Where("note_id = ? AND account_id = ? AND is_show = ?", req.NoteID, global.ParseUint(userID), true).First(&note).Error; err != nil {
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
	// 接收表单数据
	var requestData struct {
		NoteID uint `form:"note_id" binding:"required"` // 使用form标签，接收note_id字段
	}

	// 绑定FormData数据
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 获取用户ID
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1.检查数据库中是否有匹配的记录
	var note models.Notes
	if err := global.Db.Model(&models.Notes{}).Where("note_id = ? AND account_id = ? AND is_show = ?",
		requestData.NoteID, global.ParseUint(userID), true).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node file not found"})
		return
	}

	// 2.处理上传文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	// 3.定义文件存储路径
	savePath := fmt.Sprintf("D:/TalkHive/Notes/%s", note.CachePath)

	// 4.将文件保存到指定路径
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file data"})
		return
	}

	// 5.更新数据库中的文件信息
	note.SaveTime = time.Now()
	if err := global.Db.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update code metadata"})
		return
	}

	// 返回元信息给前端
	c.JSON(http.StatusOK, gin.H{
		"node_id":   note.NoteID,
		"note_name": note.Name,
		"save_time": note.SaveTime.Format("2006-01-02 15:04"),
		"type":      note.Type,
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
	userID := c.Param("id")
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
		Name:      note.Name,
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
	userID := c.Param("id")
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
	if note.Name != request.OldNoteName {
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

// DeleteNote - 删除笔记
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
	userID := c.Param("id")
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

// GetTypeList - 获取用户的笔记分类列表
func GetTypeList(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("id")
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

// GetNotesByCategory - 按分类查看笔记
func GetNotesByCategory(c *gin.Context) {
	// 获取路径参数
	// 获取用户ID
	userID := c.Param("id")
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
			"filename":     note.Name,
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
	userID := c.Param("id")
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
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1.检查 NoteDivide 表是否存在该分类
	var noteDivide models.NoteDivide
	if err := global.Db.Where("nd_name = ? AND account_id = ? AND IsShow = ?", requestData.NAName,
		global.ParseUint(userID), true).First(&noteDivide).Error; err != nil {
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
	userID := c.Param("id")
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
	userID := c.Param("id")
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
