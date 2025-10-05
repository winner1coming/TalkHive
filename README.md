# 🐝 TalkHive - 多功能网页版聊天与学习平台

TalkHive 是一个专为学生与开发者设计的**网页版聊天工具**。不仅支持基础聊天，还集成了协作文档、代码编辑、笔记记录、DDL 管理、链接收藏等学习功能，让你在一个网页中实现社交与工作的无缝融合。

---

## 📌 项目特色

### 💬 聊天模块
- 私聊、群聊支持
- 聊天记录保存
- 通讯录管理与自定义标签分类

### 📚 学习与协作
- **协作文档**：基于 Yjs + Quill ，支持多人协作
- **笔记**：集成 Quill 富文本编辑
- **代码编辑器**：集成 Monaco Editor，支持语言识别、补全、发送
- **DDL 提醒器**：任务与截止时间一览无余
- **链接器**：保存常用网址、作业提交入口等，快捷访问

### 🎨 个性化
- 支持浅色 / 深色 / 暖色主题
- 字体样式和大小、聊天背景可自定义

---

## 📷 项目展示

#### 聊天：
实现了基础的聊天对话功能，用户可以和好友聊天，也可以创建群聊实现多人聊天。
![chat](./assets/collab_send_to_friend.png)
![chat](./assets/group_chat.png)
| 群聊信息管理 |  查找聊天记录 | 按成员查找 |
|--------------------|---------------------|-------------|
![chat](./assets/group_info.png) | ![chat](./assets/search_record.png)|![chat](./assets/search_record2.png)|

#### 通讯录：
不仅支持好友和群聊管理，还允许用户自定义标签分类，帮助用户快速查找和组织联系人，提升沟通效率。
![contacts](./assets/addressBook.png)
![contacts](./assets/group_apply.png)
![contacts](./assets/group_invitation.png)

#### 工作区：
工作区集成了协作文档、笔记、代码编辑、DDL 提醒和收藏功能，为用户提供一站式的办公和学习环境，提升效率。

##### -协作文档：
支持多人在线编辑，含撤销、恢复功能（只可撤销和恢复自己修改的内容），支持实时光标显示。
![chat](./assets/collab_doc.png)

用户可将文档分享到聊天，邀请好友加入协作编辑：
![chat](./assets/collab_send.png)

##### -笔记：
可编辑笔记，修改标签，为笔记添加标签，按标签查看笔记。
![note](./assets/workspace.png)
![note](./assets/workspace_note_edit.png)

##### -代码：
代码编写功能，帮助用户更清晰地记录和阅读代码片段，提升学习和开发效率。
![chat](./assets/workspace_code.png)


##### -DDL：
支持按时间优先级排序，帮助用户更好地管理任务和 deadlines。
![chat](./assets/workspace_ddl.png)

##### -回收站：
用户可以找到已删除的文件，避免因操作失误导致的数据丢失。
![chat](./assets/workspace_recycle.png)

#### 链接器：
用户可以保存老师发放的 PPT 地址、作业提交链接或者一些常用的网址，随时一键访问，告别反复查找的麻烦。
![chat](./assets/link.png)

#### 个性化设置：
支持个性化定制，用户可调整字体样式、字体大小、主题颜色、聊天背景等，打造专属使用体验。
![chat](./assets/personalization.png)

##### -主题颜色：

| 浅色（系统默认） | 深色 | 暖色 |
|----------|----------|---------------------|
| ![login](./assets/theme_light.png) | ![chat](./assets/theme_dark.png) | ![code](./assets/theme_warm.png) |

##### -自定义聊天背景：
![chat](./assets/chat_background.png)

---

## 🛠️ 技术栈

- **前端**：Vue2 + Quill + Yjs
- **后端**：go, Redis

---

## 🚀 快速启动
### 1. 克隆项目

```bash
git clone https://github.com/winner1coming/TalkHive
cd TalkHive
```
### 2. 后端配置
#### 安装go语言环境：https://go.dev/dl/

###### 在环境变量窗口中添加 GOROOT 变量，值为 Go 的安装目录（例如 C:\Go ）
###### 添加 GOPATH 变量，值为你选择的工作区路径（可以自由设置）
###### 将 Go 安装目录下的 bin 文件夹目录添加到 Path 环境变量中
###### 打开 cmd 输入 go version 查看是否成功安装

#### 安装Redis
###### https://github.com/redis-windows/redis-windows/releases

#### 配置Redis
###### 将 Redis 安装目录配置到 Path 中（例如 D:\Software\Redis-x64-5.0.14.1 ）
#### 启动 Redis ：以下两者皆可
```bash
redis-server.exe
redis-server
```
#### 调整后端配置文件（ config 文件夹中的 config.yaml 文件）：
###### 调整 database 中的 dsn （只需调整这几个即可，其余默认）：
root 为 Mysql 中的用户名称
123456 为密码
talkhive 为数据库名称

### 3.前端配置
#### 安装前端依赖：
###### 切换到前端目录：
```bash
cd Frontend
```

###### 确保已安装 Node.js 和 npm：

```bash
node -v
npm -v
```
###### 安装依赖：
```bash
npm install
```

###### 切换到yjs-ws-server目录安装依赖
```bash
cd ..
cd yjs-ws-server
npm install
```
### 3. 切换回根目录，安装依赖并启动： 
```bash
npm install
npm run dev
```

