<template>
  <div class="contact-view">
    <FunctionBar 
      :currentTab="currentTab" 
      @tab-selected="handleTabChange" 
      class="function-bar"
    />
    <div class="main-content">
      <component :is="currentTabComponent" />
    </div>
  </div>
</template>

<script>
import FunctionBar from "@/components/Contact_list/FunctionBar.vue";
import FriendRequests from "@/components/Contact_list/FriendRequests.vue";
import GroupNotifications from "@/components/Contact_list/GroupNotifications.vue";
import FriendList from "@/components/Contact_list/FriendList.vue";
import GroupList from "@/components/Contact_list/GroupList.vue";
// import Blacklist from "@/components/Contact_list/Blacklist.vue";

export default {
  components: {
    FunctionBar,
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
  display: inline-flex;
  height: 100%;
  width: 100%;
}
.function-bar {
  flex:3;
  align-self: flex-start;
  height: 100%;
}
.main-content {
  padding: 0px;
  margin: 0;
  flex: 7;
  align-self: flex-end;
  background-color: #f8f9fa;
  overflow-y: auto;
  height: 100%;
  
}
</style>
