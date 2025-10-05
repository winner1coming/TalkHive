package workSpace

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
    "log"
    "errors"
    "gorm.io/gorm"
    "strings"
	"encoding/base64"
)

// GET /api/collab/my 获取用户创建/参与过的协作文档列表
// FetchCollabDocs - 获取用户参与过的所有协作文档（包括创建的）
// FetchCollabDocs - 获取用户参与过的所有协作文档，并返回owner昵称 √
func FetchAllCollabDocs(c *gin.Context) {
    userID := c.GetHeader("User-Id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    type DocWithOwner struct {
        DocID     uint
        DocName   string
        UpdatedAt time.Time
        OwnerName string
    }

    var docs []DocWithOwner

    // 查询用户参与过的文档，并关联创建者昵称
    if err := global.Db.Table("collab_docs").
        Select("collab_docs.doc_id, collab_docs.doc_name, collab_docs.updated_at, account_infos.nickname AS owner_name").
        Joins("JOIN collab_doc_members ON collab_docs.doc_id = collab_doc_members.doc_id").
        Joins("JOIN account_infos ON collab_docs.owner_id = account_infos.account_id").
        Where("collab_doc_members.user_id = ? AND collab_docs.is_show = ?", global.ParseUint(userID), true).
        Order("collab_docs.updated_at DESC").
        Scan(&docs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document list"})
        return
    }

    // 组织返回结果
    var result []map[string]interface{}
    for _, doc := range docs {
        result = append(result, map[string]interface{}{
            "doc_id":             doc.DocID,
            "doc_name":           doc.DocName,
            "last_modified_time": doc.UpdatedAt.Format("2006-01-02 15:04"),
            "owner_name":         doc.OwnerName,
        })
    }

    c.JSON(http.StatusOK, result)
}

// func FetchAllCollabDocs(c *gin.Context) {
//     userID := c.GetHeader("User-Id")
//     if userID == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
//         return
//     }

//     var docs []models.CollabDoc

//     // 查询参与过的文档（无论是不是自己创建）
//     if err := global.Db.Table("collab_docs").
//         Select("collab_docs.*").
//         Joins("JOIN collab_doc_members ON collab_docs.doc_id = collab_doc_members.doc_id").
//         Where("collab_doc_members.user_id = ? AND collab_docs.is_show = ?", global.ParseUint(userID), true).
//         Order("collab_docs.updated_at DESC").
//         Find(&docs).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document list"})
//         return
//     }

//     // 拼接返回结果
//     var result []map[string]interface{}
//     for _, doc := range docs {
//         result = append(result, map[string]interface{}{
//             "doc_id":             doc.DocID,
//             "doc_name":           doc.DocName,
//             "last_modified_time": doc.UpdatedAt.Format("2006-01-02 15:04"),
//         })
//     }

//     c.JSON(http.StatusOK, result)
// }

// FetchCollabDocs - 获取用户创建的协作文档列表 √
func FetchMyCollabDocs(c *gin.Context) {
    userID := c.GetHeader("User-Id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    var docs []models.CollabDoc

    // 查询当前用户创建的文档（或者可以扩展为他参与的文档）
    if err := global.Db.Table("collab_docs").
        Where("owner_id = ? AND is_show = ?", global.ParseUint(userID), true).
        Order("updated_at DESC").
        Find(&docs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document list"})
        return
    }

    // 拼接返回结果
    var result []map[string]interface{}
    for _, doc := range docs {
        result = append(result, map[string]interface{}{
            "doc_id":             doc.DocID,
            "doc_name":           doc.DocName,
            "last_modified_time": doc.UpdatedAt.Format("2006-01-02 15:04"),
        })
    }

    c.JSON(http.StatusOK, result)
}

// 创建文档
func CreateCollabDoc(c *gin.Context) {
    var req struct {
        DocName string `json:"doc_name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    userID := c.GetHeader("User-Id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
        return
    }

    doc := models.CollabDoc{
        DocName: req.DocName,
        OwnerID: global.ParseUint(userID),
        IsShow:  true,
    }

    if err := global.Db.Create(&doc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create document"})
        return
    }

    // 3. 自动把创建者加入成员表 ✅ ✅
    member := models.CollabDocMember{
        DocID:    doc.DocID,
        UserID:   doc.OwnerID,
        Role:     "owner",  // 创建者默认 owner
        JoinedAt: time.Now(),
    }
    if err := global.Db.Create(&member).Error; err != nil {
        // 如果成员关系写入失败，可以打印日志，但不影响文档创建
        log.Printf("Failed to add owner as member: %v", err)
    }

    // 初始化快照为空
    global.Db.Create(&models.CollabDocSnapshot{
        DocID: doc.DocID,
        Snapshot: []byte{}, // 空状态
    })

    c.JSON(http.StatusOK, gin.H{"doc_id": doc.DocID})
}

// 获取文档快照（前端进入协作时恢复）
func GetCollabDoc(c *gin.Context) {
    var req struct {
        DocID uint `json:"doc_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    var snapshot models.CollabDocSnapshot
    if err := global.Db.First(&snapshot, "doc_id = ?", req.DocID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "No snapshot found"})
        return
    }

    // 返回 base64 字符串给前端
    encoded := base64.StdEncoding.EncodeToString(snapshot.Snapshot)

    c.JSON(http.StatusOK, gin.H{"snapshot": encoded})
}

// 保存协作快照（前端定时）
func SaveCollabSnapshot(c *gin.Context) {
    var req struct {
        DocID    uint   `json:"doc_id" binding:"required"`
        Snapshot string `json:"snapshot" binding:"required"` // base64 编码的字符串
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    data, err := base64.StdEncoding.DecodeString(req.Snapshot)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid snapshot format"})
        return
    }

    global.Db.Model(&models.CollabDocSnapshot{}).
        Where("doc_id = ?", req.DocID).
        Updates(map[string]interface{}{
            "snapshot": data,
            "updated_at": time.Now(),
        })

    c.JSON(http.StatusOK, gin.H{"message": "Saved"})
}

// JoinCollabDocMember - 加入协作文档成员
// func JoinCollabDocMember(c *gin.Context) {
// 	// 获取用户ID
// 	userID := c.GetHeader("User-Id")
// 	if userID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
// 		return
// 	}

// 	// 接收请求体：文档ID
// 	var req struct {
// 		DocID uint `json:"doc_id" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	uid := global.ParseUint(userID)

// 	// 检查是否已有参与记录
// 	var existing models.CollabDocMember
// 	err := global.Db.Where("doc_id = ? AND user_id = ?", req.DocID, uid).First(&existing).Error
// 	if err == nil {
// 		// 已存在，不重复插入，直接返回成功
// 		c.JSON(http.StatusOK, gin.H{"message": "Already a member"})
// 		return
// 	}

// 	// 如果数据库查询其他错误
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}

// 	// 若不存在，则插入新记录
// 	member := models.CollabDocMember{
// 		DocID:    req.DocID,
// 		UserID:   uid,
// 		Role:     "editor",        // 默认角色，可按需求调整
// 		JoinedAt: time.Now(),
// 	}

// 	if err := global.Db.Create(&member).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Joined document successfully",
// 		"doc_id":  req.DocID,
// 		"user_id": uid,
// 	})
// }
func JoinCollabDocMember(c *gin.Context) {
    // 1. user id
    userID := c.GetHeader("User-Id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }
    uid := global.ParseUint(userID)

    // 2. bind body
    var req struct {
        DocID uint `json:"doc_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 3. 查询是否存在（使用 res := ... 以便同时读取 res.Error 和 res.RowsAffected）
    var existing models.CollabDocMember
    res := global.Db.Where("doc_id = ? AND user_id = ?", req.DocID, uid).First(&existing)

    if res.Error == nil {
        // 找到已有成员
        c.JSON(http.StatusOK, gin.H{"message": "Already a member"})
        return
    }

    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        // not found -> 尝试插入
        member := models.CollabDocMember{
            DocID:    req.DocID,
            UserID:   uid,
            Role:     "editor",
            JoinedAt: time.Now(),
        }

        if err := global.Db.Create(&member).Error; err != nil {
            // 可能存在并发插入导致的唯一键冲突（Duplicate entry），把它当作已加入处理
            if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key") {
                c.JSON(http.StatusOK, gin.H{"message": "Already a member"})
                return
            }
            log.Printf("failed to create collab member: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "Joined document successfully",
            "doc_id":  req.DocID,
            "user_id": uid,
        })
        return
    }

    // 4. 其他数据库错误
    log.Printf("db error when checking member: %v", res.Error)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
}
