import Vue from 'vue';
import Router from 'vue-router';

// 导入各个视图组件
import Home from '../views/Home.vue';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import MessageList from '../components/MessageList.vue';
import SingleChat from '../components/SingleChat.vue';
import GroupChat from '../components/GroupChat.vue';
import Profile from '../components/Profile.vue';
import GroupManagement from '../components/GroupManagement.vue';
import ContactList from '../components/ContactList.vue';
import FriendList from '../components/FriendList.vue';
import GroupList from '../components/GroupList.vue';
import FriendRequest from '../components/FriendRequest.vue';
import GroupRequest from '../components/GroupRequest.vue';
import PersonalHomepage from '../components/PersonalHomepage.vue';
import EditProfile from '../components/EditProfile.vue';
import SecuritySettings from '../components/SecuritySettings.vue';
import SystemSettings from '../components/SystemSettings.vue';

// 使用 Vue Router
Vue.use(Router);

// 导出路由实例
export default new Router({
  // 定义路由配置
  routes: [
    // 主页路由
    { path: '/', component: Home },
    
    // 登录页面路由
    { path: '/login', component: Login },
    
    // 注册页面路由
    { path: '/register', component: Register },
    
    // 消息列表页面路由
    { path: '/messages', component: MessageList },
    
    // 单聊页面路由，动态路由参数 :id
    { path: '/chat/:id', component: SingleChat },
    
    // 群聊页面路由，动态路由参数 :id
    { path: '/group/:id', component: GroupChat },
    
    // 个人信息页面路由
    { path: '/profile', component: Profile },
    
    // 群组管理页面路由
    { path: '/group-management', component: GroupManagement },
    
    // 联系人列表页面路由
    { path: '/contacts', component: ContactList },
    
    // 好友列表页面路由
    { path: '/friends', component: FriendList },
    
    // 群组列表页面路由
    { path: '/groups', component: GroupList },
    
    // 好友请求页面路由
    { path: '/friend-requests', component: FriendRequest },
    
    // 群组请求页面路由
    { path: '/group-requests', component: GroupRequest },
    
    // 个人主页页面路由
    { path: '/personal-homepage', component: PersonalHomepage },
    
    // 编辑资料页面路由
    { path: '/edit-profile', component: EditProfile },
    
    // 安全设置页面路由
    { path: '/security-settings', component: SecuritySettings },
    
    // 系统设置页面路由
    { path: '/system-settings', component: SystemSettings },
  ],
});