package routes

import (
	"TalkHive/controllers/workSpace"
	"github.com/gin-gonic/gin"
)

func SetupWorkspaceRoutes(router *gin.Engine) *gin.RouterGroup {
	r := router.Group("/workspace")

	// 工作区！！！

	// 工作区 - 笔记模块
	r.POST("/notes/new", workSpace.CreateNote)                       // 新建笔记（支持 MD）
	r.GET("/notes", workSpace.GetAllNotes)                           // 查看所有笔记
	r.GET("/notes/:id", workSpace.GetNote)                           // 查看特定笔记
	r.PUT("/notes/:id", workSpace.EditNote)                          // 编辑笔记
	r.DELETE("/notes/:id", workSpace.DeleteNote)                     // 删除笔记
	r.GET("/notes/category/:category", workSpace.GetNotesByCategory) // 按分类查看笔记

	// 工作区 - 我的收藏
	r.GET("/favorites", workSpace.GetFavorites)     // 查看收藏列表
	r.GET("/favorites/:id", workSpace.ViewFavorite) // 查看收藏内容

	// 工作区 - 我的代码
	r.POST("/code/new", workSpace.CreateCode)      // 新建代码文件
	r.GET("/code/:id", workSpace.GetCode)          // 查看代码文件
	r.PUT("/code/:id", workSpace.EditCode)         // 编辑代码文件
	r.POST("/code/share/:id", workSpace.ShareCode) // 分享代码文件

	// 工作区 - 搜索栏
	r.GET("/search", workSpace.SearchByKeyword) // 根据关键字匹配内容

	// 工作区 - 回收站
	r.GET("/trash", workSpace.GetTrashItems)                      // 查看回收站内容
	r.POST("/trash/restore/:id", workSpace.RestoreNote)           // 恢复回收站笔记
	r.DELETE("/trash/permanent/:id", workSpace.DeletePermanently) // 永久删除

	// 工作区 - DDL记录模块
	r.POST("/ddl/new", workSpace.CreateDDL)               // 新建DDL事件
	r.GET("/ddl", workSpace.GetAllDDL)                    // 查看所有DDL事件
	r.GET("/ddl/:id", workSpace.GetDDL)                   // 查看单个DDL事件
	r.PUT("/ddl/:id", workSpace.EditDDL)                  // 编辑DDL事件
	r.PUT("/ddl/complete/:id", workSpace.MarkDDLComplete) // 勾选完成事件
	r.DELETE("/ddl/:id", workSpace.DeleteDDL)             // 删除DDL事件
	r.GET("/ddl/reminders", workSpace.GetUpcomingDDL)     // 查看主页面DDL提醒

	return r
}
