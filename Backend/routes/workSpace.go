package routes

import (
	"TalkHive/controllers/workSpace"
	"github.com/gin-gonic/gin"
)

func SetupWorkspaceRoutes(router *gin.Engine) *gin.RouterGroup {
	r := router.Group("/workspace")

	// 工作区！！！

	// 工作区 - 笔记模块
	r.GET("/notes/list", workSpace.GetNotesList)               // 返回笔记列表√
	r.POST("/notes/get", workSpace.GetNote)                     // 获取笔记√
	r.POST("/notes/newnote", workSpace.CreateNote)             // 新建并保存笔记（支持 MD）√
	r.POST("/notes/editnote", workSpace.EditNote)               // 编辑并保存笔记√
	r.POST("/notes/share", workSpace.ShareNote)                // 分享笔记√
	r.POST("/notes/editnotename", workSpace.ChangeNoteName)    // 修改笔记名√
	r.POST("/notes/deletenote", workSpace.DeleteNote)        // 删除笔记√
	r.GET("/notes/categories", workSpace.GetTypeList)          // 获取分类列表√
	r.GET("/notes/dividenotes", workSpace.GetNotesByCategory)  // 按分类查看笔记√
	r.POST("/notes/editnotetype", workSpace.EditNoteType)      // 修改笔记所在分类√
	r.POST("/notes/categories/delete", workSpace.DeleteType)   // 删除分类√
	r.POST("/notes/categories/edit", workSpace.ChangeTypeName) //修改分类名称√
	r.POST("/notes/categories/new", workSpace.CreateType)      // 新建分类√

	// 工作区 - 我的收藏
	r.GET("/favorites/list", workSpace.GetFavorites)               // 返回收藏列表
	r.GET("/favorites/get", workSpace.ViewFavorite)                // 查看收藏内容
	r.POST("/favorites/add", workSpace.AddFavorite)                // 新增收藏
	r.POST("/favorites/delete", workSpace.DeleteMultipleFavorites) // 批量删除收藏

	// 工作区 - 我的代码
	r.POST("/code/list", workSpace.SearchCode)   // 返回代码文件列表√
	r.POST("/code/get", workSpace.GetCode)       // 获取代码文件内容√
	r.POST("/code/new", workSpace.CreateCode)    // 新建并保存代码文件√
	r.POST("/code/edit", workSpace.EditCode)      // 编辑代码文件√
	r.POST("/code/share", workSpace.ShareCode)   // 分享代码文件√
	r.POST("/code/name", workSpace.ChangeName)   // 修改代码文件名√
	r.POST("/code/delete", workSpace.DeleteCode) // 删除代码文件√

	// 工作区 - 搜索栏
	//r.GET("/search", workSpace.SearchByKeyword) // 根据关键字匹配内容

	// 工作区 - 回收站
	r.GET("/recycle/files", workSpace.GetTrashItems)             // 获取回收站列表√
	r.POST("/recycle/restore-file", workSpace.RestoreItem)           // 恢复回收站笔记√
	r.POST("/recycle/delete-permanent", workSpace.DeletePermanently) // 永久删除√

	// 工作区 - DDL记录模块
	r.GET("/ddl/pending", workSpace.GetUncompletedDDL) // 查看待完成DDL事件√
	r.GET("/ddl/completed", workSpace.GetCompletedDDL) // 查看已完成DDL事件√
	r.POST("/ddl/create", workSpace.CreateDDL)         // 新建DDL事件√
	r.PUT("/ddl/update", workSpace.EditDDL)            // 编辑DDL事件√
	r.POST("/ddl/complete", workSpace.MarkDDLComplete)  // 勾选完成事件√
	r.POST("/ddl/delete", workSpace.DeleteDDL)       // 删除DDL事件√
	r.GET("/ddl/reminders", workSpace.GetUpcomingDDL)      // 查看主页面DDL提醒

	return r
}
