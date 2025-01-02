<template>
  <div class="function-bar">
    <!-- 搜索框-->
    <SearchBar 
      @search="handleSearch" 
      @button-click="showContextMenu($event)"
      :isImmidiate="true"
      ref="searchBar"
    />
    <ul class="menu">
      <li
        v-for="item in menuItems"
        :key="item.name"
        :class="{ active: currentTab === item.name }"
        @click="selectTab(item.name)"
      >
        {{ item.label }}
      </li>
    </ul>
    <!-- 添加好友弹窗 -->
    <AddFriendGroup
      v-if="isAddModalVisible"
      @close="isAddModalVisible = false"
    />
    <!-- 新建群聊弹窗 -->
    <BuildGroup
      v-if="isBuildModalVisible"
      @close="isBuildModalVisible = false"
    />
    <ContextMenu ref="contextMenu"  @select-item="handleMenu" />
  </div>
</template>
  
<script>
import SearchBar from '@/components/base/SearchBar.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import AddFriendGroup from '@/components/Chat_list/AddFriendGroup.vue';
import BuildGroup from '@/components/Chat_list/BuildGroup.vue';
export default {
  props: {
    currentTab: String, // 当前选中的标签
  },
  components: {
    SearchBar,
    ContextMenu,
    AddFriendGroup,
    BuildGroup,
  },
  data() {
    return {
      isAddModalVisible: false,
      isBuildModalVisible: false,
      menuItems: [
        { name: "friendRequests", label: "好友申请" },
        { name: "groupNotifications", label: "群聊通知" },
        { name: "friendList", label: "好友列表" },
        { name: "groupList", label: "群组列表" },
        { name: "blacklist", label: "黑名单" },
      ],
    };
  },
  methods: {
    selectTab(tabName) {
      this.$refs.searchBar.clear(); // 清空搜索框
      this.$emit("tab-selected", tabName); // 触发父组件事件
    },
    // 搜索消息
    async handleSearch(keyword) {
      this.$emit('search', keyword);
    },
    // 显示的菜单
    showContextMenu(event) {
      this.menuType = 'new';
      const items = [
        '添加好友/群聊',
        '新建群聊',
      ];
      this.$refs.contextMenu.show(event, items, null, null, null);
    },
    // 处理菜单点击事件
    async handleMenu(option) {
      if(option === '添加好友/群聊') {
        this.isAddModalVisible = true;
      }else if(option === '新建群聊') {
        this.isBuildModalVisible = true;
      }
    },
  },
};
</script>
  
<style scoped>
.function-bar {
  width: 250px;
  background-color: #f8f9fa;
  border-right: 1px solid #ddd;
}
.menu {
  list-style: none;
  padding: 0;
}
.menu li {
  padding: 15px;
  cursor: pointer;
  border-bottom: 1px solid #ddd;
}
.menu li:hover {
  background-color: #e9ecef;
}
.menu li.active {
  background-color: #007bff;
  color: #fff;
}
</style>
  