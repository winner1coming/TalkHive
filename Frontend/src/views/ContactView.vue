<template>
    <div class="contact-view">
      <Sidebar :currentTab="currentTab" @tab-selected="handleTabChange" />
      <div class="main-content">
        <component :is="currentTabComponent" />
      </div>
    </div>
  </template>
  
  <script>
  import Sidebar from "@/components/Contact_list/FunctionBar.vue";
  import FriendRequests from "@/components/Contact_list/FriendRequests.vue";
  import GroupNotifications from "@/components/Contact_list/GroupNotifications.vue";
  import FriendList from "@/components/Contact_list/FriendList.vue";
  import GroupList from "@/components/Contact_list/GroupList.vue";
  // import Blacklist from "@/components/Contact_list/Blacklist.vue";
  
  export default {
    components: {
      Sidebar,
      FriendRequests,
      GroupNotifications,
      FriendList,
      GroupList,
    },
    data() {
      return {
        currentTab: "FriendRequests", // 默认展示“好友申请”
      };
    },
    computed: {
      currentTabComponent() {
        const components = {
          friendRequests: "FriendRequests",
          groupNotifications: "GroupNotifications",
          friendList: "FriendList",
          groupList: "GroupList",
          blacklist: "Blacklist",
        };
        return components[this.currentTab];
      },
    },
    methods: {
      handleTabChange(tabName) {
        this.currentTab = tabName;
      },
    },
  };
  </script>
  
  <style scoped>
  .contact-view {
    display: flex;
    height: 100vh;
  }
  .main-content {
    flex: 1;
    padding: 20px;
    background-color: #f8f9fa;
    overflow-y: auto;
  }
  </style>
  