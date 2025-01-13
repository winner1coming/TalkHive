<template>
  <div class="contact-view">
    <FunctionBar 
      :currentTab="currentTab" 
      :style="{ width: functionBarWidth + 'px' }"
      @tab-selected="handleTabChange" 
      @search="handleSearch"
      class="function-bar"
    />

    <!-- 拖动条 -->
    <div class="resizer" @mousedown="startResize"></div>
    
    <div class="main-content">
      <FriendRequests v-show="currentTab === 'friendRequests'" />
      <GroupNotifications v-show="currentTab === 'groupNotifications'" />
      <FriendList v-show="currentTab === 'friendList'" />
      <GroupList v-show="currentTab === 'groupList'" />
      <Blacklist v-show="currentTab === 'blacklist'" />
      <SearchList v-show="currentTab === 'searchList'" :keyword="searchKeyword"/>
    </div>
  </div>
</template>

<script>
import FunctionBar from "@/components/Contact_list/FunctionBar.vue";
import FriendRequests from "@/components/Contact_list/FriendRequests.vue";
import GroupNotifications from "@/components/Contact_list/GroupNotifications.vue";
import FriendList from "@/components/Contact_list/FriendList.vue";
import GroupList from "@/components/Contact_list/GroupList.vue";
import Blacklist from "@/components/Contact_list/Blacklist.vue";
import SearchList from "@/components/Contact_list/SearchList.vue";

export default {
  components: {
    FunctionBar,
    FriendRequests,
    GroupNotifications,
    FriendList,
    GroupList,
    Blacklist,
    SearchList,
  },
  data() {
    return {
      currentTab: "FriendRequests", // 默认展示“好友申请”
      searchKeyword: "",  // 搜索关键词
      functionBarWidth: 230,  // 功能栏的宽度
      leftComponentWidth: 80,  // 左侧组件的宽度 todo
    };
  },
  // computed: {
  //   currentTabComponent() {
  //     const components = {
  //       friendRequests: "FriendRequests",
  //       groupNotifications: "GroupNotifications",
  //       friendList: "FriendList",
  //       groupList: "GroupList",
  //       blacklist: "Blacklist",
  //       searchList: "SearchList",
  //     };
  //     return components[this.currentTab];
  //   },
  // },
  methods: {
    handleTabChange(tabName) {
      this.currentTab = tabName;
    },
    handleSearch(query) {
      this.currentTab = "searchList";
      this.searchKeyword = query;
    },
    // 拖动条的逻辑
    startResize(event) {
      this.isResizing = true;
      document.addEventListener('mousemove', this.resize);
      document.addEventListener('mouseup', this.stopResize);
    },
    resize(event) {
      if (this.isResizing) {
        this.functionBarWidth = event.clientX - this.leftComponentWidth;
      }
    },
    stopResize() {
      this.isResizing = false;
      document.removeEventListener('mousemove', this.resize);
      document.removeEventListener('mouseup', this.stopResize);
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
  align-self: flex-start;
  height: 100%;
}
.main-content {
  padding: 0px;
  margin: 0;
  align-self: flex-end;
  background-color: var(--background-color);
  color: var(--text-color);
  overflow-y: auto;
  height: 100%;
  width: 100%;
}
.resizer {
  width: 3px;
  height: 100%;
  cursor: ew-resize;
  background-color: #ccc;
}
</style>
