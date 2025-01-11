package workSpace

import (
	"TalkHive/global"
	"TalkHive/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// 我的收藏！！！

// GetFavorites - 查看收藏列表
// ①打开数据库favorites(table_name, id, account_id)，使用select选出属性account_id为用户id的数据，逐条遍历。
// 其中table_name取值为Notes\Codes\MessageInfo，id即为这条收藏的id值
// （与Notes中note_id、Codes的code_id、MessageInfo的message_id一致）
// ②	对筛选后的数据遍历，根据table_name取值，使用对应的表Notes\Codes\MessageInfo，
// 再根据收藏表该数据的id取值，在对应的索引项note_id\code_id\message_id找到该收藏在Notes\Codes\MessageInfo的具体信息。
// ③	返回值type根据该条收藏所在表决定。
// ④	Message_list_name：消息表名Notes\Codes\MessageInfo
// ⑤	如果该收藏是消息，sender_name为MessageInfo对应数据里的send_account_id；是代码或笔记，则为用户id本身
// ⑥	返回值Message_id为favorites的id值
// ⑦	返回值Object_name：如果该收藏是消息，Object_name为MessageInfo对应数据里的Content；是代码或笔记，则Name
// ⑧	Time都可以从三个表中获取
func GetFavorites(c *gin.Context) {
	// 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 查询用户收藏的所有记录
	var favorites []models.Favorites
	if err := global.Db.Where("account_id = ?", userID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch favorites"})
		return
	}

	var response []gin.H // 用于存储最终的响应数据

	// 遍历收藏记录
	for _, favorite := range favorites {
		var item gin.H
		item = gin.H{
			"message_list_name": favorite.TableName, // 收藏对应的表名
			"message_id":        favorite.ID,        // 收藏ID
		}

		switch favorite.TableName {
		case "note":
			var note models.Notes
			if err := global.Db.Where("note_id = ?", favorite.ID).First(&note).Error; err != nil {
				continue // 如果找不到对应的笔记，跳过这条记录
			}
			item["type"] = "note"
			item["object_name"] = note.NoteName
			//item["sender_name"] = fmt.Sprintf("%d", note.AccountID) // 笔记的sender_name就是用户ID
			// 获取 sender_name
			var sender models.AccountInfo
			if err := global.Db.Where("account_id = ?", note.AccountID).First(&sender).Error; err != nil {
				item["sender_name"] = fmt.Sprintf("%d", note.AccountID) // 如果获取不到昵称，返回 AccountID
			} else {
				item["sender_name"] = sender.Nickname
			}
			item["time"] = note.SaveTime.Format("2006-01-02 15:04")

		case "code":
			var code models.Codes
			if err := global.Db.Where("code_id = ?", favorite.ID).First(&code).Error; err != nil {
				continue // 如果找不到对应的代码，跳过这条记录
			}
			item["type"] = "code"
			item["object_name"] = code.Name + code.Suffix           // 代码名+后缀名
			//item["sender_name"] = fmt.Sprintf("%d", code.AccountID) // 代码的sender_name就是用户ID
			// 获取 sender_name
			var sender models.AccountInfo
			if err := global.Db.Where("account_id = ?", code.AccountID).First(&sender).Error; err != nil {
				item["sender_name"] = fmt.Sprintf("%d", code.AccountID) // 如果获取不到昵称，返回 AccountID
			} else {
				item["sender_name"] = sender.Nickname
			}
			item["time"] = code.SaveTime.Format("2006-01-02 15:04")

		case "message":
			var message models.MessageInfo
			if err := global.Db.Where("message_id = ?", favorite.ID).First(&message).Error; err != nil {
				continue // 如果找不到对应的消息，跳过这条记录
			}
			item["type"] = "message"
			item["object_name"] = message.Content
			//item["sender_name"] = fmt.Sprintf("%d", message.SendAccountID) // 消息的sender_name是发送者的ID
			// 获取 sender_name
			var sender models.AccountInfo
			if err := global.Db.Where("account_id = ?", message.SendAccountID).First(&sender).Error; err != nil {
				item["sender_name"] = fmt.Sprintf("%d", message.SendAccountID) // 如果获取不到昵称，返回 AccountID
			} else {
				item["sender_name"] = sender.Nickname
			}
			item["time"] = message.CreateTime.Format("2006-01-02 15:04")   // 消息的时间

		default:
			continue // 如果不是Notes, Codes, MessageInfo表，跳过此条记录
		}

		response = append(response, item)
	}

	// 返回用户收藏的列表
	c.JSON(http.StatusOK, response)
}

// ViewFavorite - 查看收藏内容
func ViewFavorite(c *gin.Context) {

	var seq struct {
		Type      string `json:"message_type" binding:"required"`
		ID        uint   `json:"message_id" binding:"required"`
		tableName string `json:"message_table_name" binding:"required"`
	}

	// 1. 获取前端传来的数据并验证
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3.分别调用message、code、note的获取函数
	switch seq.Type {
	case "message":
		// 如果未来需要获取 Message 的逻辑，可以在这里实现
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Message retrieval is not implemented yet"})
	case "note":
		// 调用 GetNote 并传递上下文
		c.Set("note_id", seq.ID) // 将 NoteID 放入上下文
		GetNote(c)
	case "code":
		// 调用 GetCode 并传递上下文
		c.Set("code_id", seq.ID) // 将 NoteID 放入上下文
		GetNote(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type, must be 'message', 'note', or 'code'"})
		return

	}
}

// AddFavorite - 添加收藏
func AddFavorite(c *gin.Context) {
	var seq struct {
		ID   uint   `json:"id" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	// 1. 获取前端传来的数据并验证
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3. 根据类型选择对应的表并验证数据
	switch seq.Type {
	case "message":
		var message models.MessageInfo
		if err := global.Db.Where("send_account_id = ? AND message_id = ?",
			userID, seq.ID).First(&message).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Message not found or doesn't belong to user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case "note":
		var note models.Notes
		if err := global.Db.Where("account_id = ? AND note_id = ? AND is_show = ?",
			userID, seq.ID, true).First(&note).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or doesn't belong to user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case "code":
		var code models.Codes
		if err := global.Db.Where("account_id = ? AND code_id = ? AND is_show = ?",
			userID, seq.ID, true).First(&code).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Code file not found or doesn't belong to user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type, must be 'message', 'note', or 'code'"})
		return
	}

	// 4. 检查用户是否已收藏此数据
	var favorite models.Favorites
	if err := global.Db.Where("account_id = ? AND table_name = ? AND id = ?", userID, seq.Type, seq.ID).First(&favorite).Error; err == nil {
		// 如果已经收藏，则返回提示
		c.JSON(http.StatusConflict, gin.H{"message": "This item is already in your favorites"})
		return
	}

	// 5. 新增收藏记录
	favorite = models.Favorites{
		TableName: seq.Type, // 表示收藏的类型
		ID:        seq.ID,
		AccountID: global.ParseUint(userID),
	}

	if err := global.Db.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to favorites", "message": err.Error()})
		return
	}

	// 6. 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites successfully"})
}

// DeleteFavorite - 删除收藏
func DeleteFavorite(c *gin.Context) {
	var seq struct {
		ID   uint   `json:"id" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	// 1. 获取前端传来的数据并验证
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3. 根据类型选择对应的表并验证数据
	switch seq.Type {
	case "message":
		var message models.MessageInfo
		if err := global.Db.Where("send_account_id = ? AND message_id = ?",
			userID, seq.ID).First(&message).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Message not found or doesn't belong to user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case "note":
		var note models.Notes
		if err := global.Db.Where("account_id = ? AND note_id = ? AND is_show = ?",
			userID, seq.ID, true).First(&note).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or doesn't belong to user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case "code":
		var code models.Codes
		if err := global.Db.Where("account_id = ? AND code_id = ? AND is_show = ?",
			userID, seq.ID, true).First(&code).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Code file not found or doesn't belong to user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type, must be 'message', 'note', or 'code'"})
		return
	}

	// 5. 删除 Favorites 表中对应的记录
	if err := global.Db.Where("account_id = ? AND table_name = ? AND id = ?", userID, seq.Type, seq.ID).Delete(&models.Favorites{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from favorites", "message": err.Error()})
		return
	}

	// 4. 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Deleted from favorites successfully"})
}

// DeleteMultipleFavorites - 批量删除收藏
func DeleteMultipleFavorites(c *gin.Context) {
	var seq []struct {
		ID   uint   `json:"message_id" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	// 1. 获取前端传来的数据并验证
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取用户ID
	//userID := c.Param("id")
	userID := c.GetHeader("User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3. 循环处理每个收藏项
	for _, item := range seq {
		// 处理每个收藏项的删除逻辑
		var err error
		switch item.Type {
		case "message":
			var message models.MessageInfo
			err = global.Db.Where("send_account_id = ? AND message_id = ?",
				userID, item.ID).First(&message).Error
			if err == nil {
				err = global.Db.Where("account_id = ? AND table_name = ? AND id = ?",
					userID, "message", item.ID).Delete(&models.Favorites{}).Error
			}

		case "note":
			var note models.Notes
			err = global.Db.Where("account_id = ? AND note_id = ? AND is_show = ?",
				userID, item.ID, true).First(&note).Error
			if err == nil {
				err = global.Db.Where("account_id = ? AND table_name = ? AND id = ?",
					userID, "note", item.ID).Delete(&models.Favorites{}).Error
			}

		case "code":
			var code models.Codes
			err = global.Db.Where("account_id = ? AND code_id = ? AND is_show = ?",
				userID, item.ID, true).First(&code).Error
			if err == nil {
				err = global.Db.Where("account_id = ? AND table_name = ? AND id = ?",
					userID, "code", item.ID).Delete(&models.Favorites{}).Error
			}

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type, must be 'message', 'note', or 'code'"})
			return
		}

		// 4. 错误处理
		if err != nil {
			// 如果出错，返回错误信息并停止处理
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 5. 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Deleted from favorites successfully"})
}
