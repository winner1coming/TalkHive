前端代码的具体框架：
Frontend
├── README.md
├── package-lock.json
├── package.json
├── public
│   └── index.html
├── src
│   ├── App.vue
│   ├── assets
│   │   ├── css     # 样式
│   │   ├── icon    # 图标
│   │   ├── images  # 图片
│   │   └── sounds  # 音效
│   ├── components
│   │   ├── AccountLogin.vue    # 登录界面
│   │   ├── Chat_list       #聊天的功能区
│   │   │   ├── AddFriendGroup.vue
│   │   │   ├── BuildGroup.vue
│   │   │   ├── ChatBox.vue
│   │   │   ├── ChatBoxLazyLoad.vue
│   │   │   ├── ChatList.vue
│   │   │   ├── CodeEdit.vue
│   │   │   ├── Emoji.vue
│   │   │   ├── FriendManagement.vue
│   │   │   ├── GroupChat.vue
│   │   │   ├── GroupManagement.vue
│   │   │   ├── InviteMember.vue
│   │   │   ├── MemberSelect.vue
│   │   │   ├── MessageInput.vue
│   │   │   ├── MessageItem.vue
│   │   │   └── SingleChat.vue
│   │   ├── Contact_list    # 通讯录的功能区
│   │   │   ├── BlackList.vue
│   │   │   ├── BlackListManagement.vue
│   │   │   ├── DivideAdd.vue
│   │   │   ├── DivideDelete.vue
│   │   │   ├── DivideManagement.vue
│   │   │   ├── DivideMove.vue
│   │   │   ├── FriendList.vue
│   │   │   ├── FriendRequests.vue
│   │   │   ├── FunctionBar.vue
│   │   │   ├── GroupList.vue
│   │   │   ├── GroupNotifications.vue
│   │   │   ├── HeaderToggle.vue
│   │   │   ├── SearchList.vue
│   │   │   └── itemList.vue
│   │   ├── ForgetPassword.vue
│   │   ├── Register.vue
│   │   ├── Settings_list   # 设置的功能区
│   │   │   ├── ChangeEmail.vue
│   │   │   ├── ChangePassword.vue
│   │   │   ├── ChatBackground.vue
│   │   │   ├── EditProfile.vue
│   │   │   ├── FontSize.vue
│   │   │   ├── FriendPermission.vue
│   │   │   ├── Profile.vue
│   │   │   ├── SecuritySettings.vue
│   │   │   ├── SoundSetting.vue
│   │   │   ├── SystemSettings.vue
│   │   │   └── ThemeSetting.vue
│   │   ├── SmsLogin.vue
│   │   ├── WorkSpace   # 个人工作空间的功能区
│   │   │   ├── Code.vue
│   │   │   ├── CodeEditor.vue
│   │   │   ├── CollabEditor.vue
│   │   │   ├── DdlList.vue
│   │   │   ├── Favorites.vue
│   │   │   ├── NoteEditor.vue
│   │   │   ├── Notes.vue
│   │   │   ├── OnlineCollaboration.vue
│   │   │   ├── QuillEditor.vue
│   │   │   ├── Recycle.vue
│   │   │   └── SelectFriend.vue
│   │   └── base    # 基础组件
│   │       ├── BuildGroup.vue
│   │       ├── ContextMenu.vue
│   │       ├── EditableText.vue
│   │       ├── EventBus.vue
│   │       ├── GroupProfileCard.vue
│   │       ├── Notification.vue
│   │       ├── PersonProfileCard.vue
│   │       ├── SearchBar.vue
│   │       ├── SearchLink.vue
│   │       ├── SwitchButton.vue
│   │       ├── ToggleContent.vue
│   │       └── Windows.vue
│   ├── main.ts
│   ├── router
│   │   └── index.js    # 路由配置文件
│   ├── services        # 与后端连接的接口文件
│   │   ├── api.js
│   │   ├── chatList.js
│   │   ├── contactList.js
│   │   ├── loginth.js
│   │   ├── mock
│   │   ├── settingView.js
│   │   └── workspace_api.js
│   ├── store   # Vuex 状态管理配置文件
│   │   └── index.js
│   ├── utils   # 本地存储工具文件
│   │   └── storage.js
│   └── views   # 主页视图组件
│       ├── ChatView.vue
│       ├── ContactView.vue
│       ├── Home.vue
│       ├── Link.vue
│       ├── LoginTH.vue
│       ├── SettingView.vue
│       └── WorkSpace.vue
└── vite.config.js
