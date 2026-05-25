import { createRouter, createWebHistory } from 'vue-router';

// 导入各个视图组件,懒加载
const Home = () => import('@/views/Home.vue');
const LoginTH = () => import('@/views/LoginTH.vue');
const Register = () => import('@/components/Register.vue');
const ChatView = () => import('@/views/ChatView.vue');
const ContactView = () => import('@/views/ContactView.vue');
const ForgetPassword = () => import('@/components/ForgetPassword.vue');
const ShowList = () => import('@/views/SettingView.vue');
const WorkSpace = () => import('@/views/WorkSpace.vue');
const Favorites = () => import('@/components/WorkSpace/Favorites.vue');
const Notes = () => import("@/components/WorkSpace/Notes.vue");
const NoteEditor = () => import("@/components/WorkSpace/NoteEditor.vue");
const CodeEditor = () => import("@/components/WorkSpace/CodeEditor.vue");
const Code = () => import("@/components/WorkSpace/Code.vue");
const DdlList = () => import("@/components/WorkSpace/DdlList.vue");
const Recycle = () => import("@/components/WorkSpace/Recycle.vue");
const OnlineCollaboration = () => import('@/components/WorkSpace/OnlineCollaboration.vue');
const CollabEditor = () => import('@/components/WorkSpace/CollabEditor.vue');
const QuillEditor = () => import('@/components/WorkSpace/QuillEditor.vue');


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
      { path: '/chat', name: 'chat', component: ChatView },
      { path: '/contact', component: ContactView },
      {path:'/setlist',component:ShowList},
      {
        path: "/workspace",
        component: WorkSpace,
        children: [
          {
            path:"collabdocs",
            component: OnlineCollaboration,
          },
          {
            path:"collabdocs/editor",
            component: CollabEditor,
          },
          {
            path: "favorites",
            component: Favorites,
          },
          {
            path: "notes",
            component: Notes, 
          },
          {
            path: "code",
            component: Code,
          },
          {
            path: "ddl",
            component: DdlList,
          },
          {
            path: "notes/editor",
            component: QuillEditor,
          },
          {
            path: "code/editor",
            component: CodeEditor,
          },
          {
            path: "recycle",
            component: Recycle,
          },
          // 添加其他子路由
        ]
      }
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