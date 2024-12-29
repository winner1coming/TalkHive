package routes

import (
	"TalkHive/controllers/workSpace"
	"github.com/gin-gonic/gin"
)

func SetupWorkspaceRoutes(router *gin.Engine) *gin.RouterGroup {
	r := router.Group("/workspace")

	// 工作区！！！

	// 工作区 - 笔记模块
	r.GET("/notes/list/:id", workSpace.GetNotesList)               // 返回笔记列表√
	r.GET("/notes/get/:id", workSpace.GetNote)                     // 获取笔记√
	r.POST("/notes/newnote/:id", workSpace.CreateNote)             // 新建并保存笔记（支持 MD）√
	r.PUT("/notes/editnote/:id", workSpace.EditNote)               // 编辑并保存笔记√
	r.POST("/notes/Share/:id", workSpace.ShareNote)                // 分享笔记√
	r.POST("/notes/editnotename/:id", workSpace.ChangeNoteName)    // 修改笔记名√
	r.DELETE("/notes/deletenote/:id", workSpace.DeleteNote)        // 删除笔记√
	r.GET("/notes/categories/:id", workSpace.GetTypeList)          // 获取分类列表√
	r.GET("/notes/dividenotes/:id", workSpace.GetNotesByCategory)  // 按分类查看笔记√
	r.POST("/notes/editnotetype/:id", workSpace.EditNoteType)      // 修改笔记所在分类√
	r.POST("/notes/categories/delete/:id", workSpace.DeleteType)   // 删除分类√
	r.POST("/notes/categories/edit/:id", workSpace.ChangeTypeName) //修改分类名称√
	r.POST("/notes/categories/new/:id", workSpace.CreateType)      // 新建分类√

	// 工作区 - 我的收藏
	r.GET("/favorites/list/:id", workSpace.GetFavorites)               // 返回收藏列表
	r.GET("/favorites/get/:id", workSpace.ViewFavorite)                // 查看收藏内容
	r.POST("/favorites/add/:id", workSpace.AddFavorite)                // 新增收藏
	r.POST("/favorites/delete/:id", workSpace.DeleteMultipleFavorites) // 批量删除收藏

	// 工作区 - 我的代码
	r.POST("/code/list/:id", workSpace.SearchCode)   // 返回代码文件列表√
	r.POST("/code/get/:id", workSpace.GetCode)       // 获取代码文件内容√
	r.POST("/code/new/:id", workSpace.CreateCode)    // 新建并保存代码文件√
	r.PUT("/code/edit/:id", workSpace.EditCode)      // 编辑代码文件√
	r.POST("/code/share/:id", workSpace.ShareCode)   // 分享代码文件√
	r.POST("/code/name/:id", workSpace.ChangeName)   // 修改代码文件名√
	r.POST("/code/delete/:id", workSpace.DeleteCode) // 删除代码文件√

	// 工作区 - 搜索栏
	r.GET("/search", workSpace.SearchByKeyword) // 根据关键字匹配内容

	// 工作区 - 回收站
	r.GET("/trash/list/:id", workSpace.GetTrashItems)             // 获取回收站列表√
	r.POST("/trash/restore/:id", workSpace.RestoreItem)           // 恢复回收站笔记√
	r.DELETE("/trash/permanent/:id", workSpace.DeletePermanently) // 永久删除√

	// 工作区 - DDL记录模块
	r.GET("/ddl/pending/:id", workSpace.GetUncompletedDDL) // 查看待完成DDL事件√
	r.GET("/ddl/completed/:id", workSpace.GetCompletedDDL) // 查看已完成DDL事件√
	r.POST("/ddl/create/:id", workSpace.CreateDDL)         // 新建DDL事件√
	r.PUT("/ddl/update/:id", workSpace.EditDDL)            // 编辑DDL事件√
	r.PUT("/ddl/complete/:id", workSpace.MarkDDLComplete)  // 勾选完成事件√
	r.DELETE("/ddl/delete/:id", workSpace.DeleteDDL)       // 删除DDL事件√
	r.GET("/ddl/reminders", workSpace.GetUpcomingDDL)      // 查看主页面DDL提醒

	return r
}
