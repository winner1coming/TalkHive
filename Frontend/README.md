以下先给出前端代码的具体框架：
    现在阶段可考虑先实现Home.vue的主界面构建的跳转逻辑设置。
    同时项目入口文件主要是登录注册页面的实现。
    各个组件可分别开发，先假设有对应的数据库

Frontend/
├── public/
│   ├── index.html  # 项目的入口 HTML 文件
├── src/
│   ├── assets/
│   │   ├── css/
│   │   │   ├── main.css  # 全局样式文件
│   │   ├── images/
│   │   │   ├── logo.png  # 项目中使用的图片文件
│   ├── components/
│   │   │   ├── Contact_list/ # 通讯录的功能区
│   │   │   │   ├── ContactList.vue  # 通讯录组件
│   │   │   │   ├── FriendList.vue  # 好友列表组件
│   │   │   │   ├── FriendRequest.vue  # 好友申请组件
│   │   │   │   ├── GroupList.vue  # 群聊列表组件
│   │   │   │   ├── GroupRequest.vue  # 群聊申请组件
│   │   │   ├── Message_list/ # 消息列表的功能区
│   │   │   │   ├── GroupChat.vue  # 群聊组件
│   │   │   │   ├── GroupManagement.vue  # 群聊管理组件
│   │   │   │   ├── MessageList.vue  # 消息列表组件
│   │   │   │   ├── SingleChat.vue  # 单聊组件
│   │   │   ├── Settings_list/ # 个人主页的功能区
│   │   │   │   ├── EditProfile.vue  # 编辑个人信息组件
│   │   │   │   ├── PersonalHomepage.vue  # 个人主页组件
│   │   │   │   ├── Profile.vue  # 个人信息组件
│   │   │   │   ├── SecuritySettings.vue  # 安全设置组件
│   │   │   │   ├── SystemSettings.vue  # 系统设置组件
│   │   │   ├── Tools_list/ # 笔记的功能区
│   │   │   │   ├── DDLs.vue #DDL实现的组件
│   │   │   │   ├── Favorites.vue #收藏组件实现
│   │   │   │   ├── Notes.vue #笔记组件实现
│   │   ├── Login.vue  # 登录组件
│   │   ├── Register.vue  # 注册组件
│   ├── router/
│   │   ├── index.js  # 路由配置文件
│   ├── services/
│   │   ├── api.js  # 与后端连接的接口文件
│   ├── store/
│   │   ├── index.js  # Vuex 状态管理配置文件
│   ├── utils/
│   │   ├── storage.js  # 本地存储工具文件
│   ├── views/
│   │   ├── Home.vue  # 主页视图组件
│   ├── App.vue  # 主应用组件
│   ├── main.js  # 项目入口文件
├── README.md  # 项目说明文件
├── package.json  # 项目依赖配置文件