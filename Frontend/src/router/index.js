import { createRouter, createWebHistory } from 'vue-router';

// 导入各个视图组件
import Home from '@/views/Home.vue';
import LoginTH from '@/views/LoginTH.vue';
import Register from '@/components/Register.vue';
import ChatView from '@/views/ChatView.vue';
import ContactView from '@/views/ContactView.vue';
import ForgetPassword from '@/components/ForgetPassword.vue';
import ShowList from '@/views/SettingView.vue';
import Workspace from '@/views/Workspace.vue';

//import GroupChat from '@/components/GroupChat.vue';
//import Profile from '@/components/Profile.vue';
//import GroupManagement from '@/components/GroupManagement.vue';
//import ContactList from '@/components/ContactList.vue';
//import FriendList from '@/components/FriendList.vue';
//import GroupList from '@/components/GroupList.vue';
// import FriendRequest from '@/components/FriendRequest.vue';
// import GroupRequest from '@/components/GroupRequest.vue';
// import PersonalHomepage from '@/components/PersonalHomepage.vue';
// import EditProfile from '@/components/EditProfile.vue';
// import SecuritySettings from '@/components/SecuritySettings.vue';
// import SystemSettings from '@/components/SystemSettings.vue';

// 定义路由配置
const routes = [
  {path: '/', component: LoginTH},
  {path:'/loginth',component:LoginTH},
  { path: '/register', component: Register },
  {path:'/forgetpassword',component:ForgetPassword},
  { 
    path: '/home', 
    component: Home,
    meta:{requiresAuth:true},
    children:[
      { path: '/chat', component: ChatView },
      { path: '/contact', component: ContactView },
      {path:'/setlist',component:ShowList},
      {path:'/workspace',component:Workspace},
    ]
  },
  //{ path: '/profile', component: Profile },
   //{ path: '/group-management', component: GroupManagement },
  // { path: '/friends', component: FriendList },
 // { path: '/groups', component: GroupList },
  // { path: '/friend-requests', component: FriendRequest },
  // { path: '/group-requests', component: GroupRequest },
  // { path: '/personal-homepage', component: PersonalHomepage },
  // { path: '/edit-profile', component: EditProfile },
  // { path: '/security-settings', component: SecuritySettings },
  // { path: '/system-settings', component: SystemSettings },
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
});

/*router.beforeEach((to, from, next) => {
    const isLoggedIn = localStorage.getItem(('isLoggedIn')==='true');

    if(to.meta.requiresAuth){
      if(isLoggedIn){
        next();
      }
      else{
        next('/loginth');
      }
    }
    else{
      next();
    }
});*/

export default router;