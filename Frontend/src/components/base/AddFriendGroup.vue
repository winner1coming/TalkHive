<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <h2>加好友/群</h2>
      <SearchBar :isImmidiate="true" @search="search" @button-click="search"/>
    </div>
  </div>
</template>

<script>
import { addFriendGroup, searchFriendGroup } from '@/services/api';
import SearchBar from '@/components/base/SearchBar.vue';
export default {
  data() {
    return {
      friendUsername: '',
    };
  },
  methods: {
    search(query) {
      // 调用API搜索用户
      console.log('searching...', query);
    },
    addFriend() {
      // 调用API添加好友
      this.$emit('add-friend', this.friendUsername);
      this.close();
    },
    close() {
      this.$emit('close');
    },
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000; /* 确保在最上层 */
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 300px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.modal-actions button {
  margin-left: 10px;
}
</style>