<template>
  <div class="message-item" @contextmenu.prevent="showContextMenu($event)">
    <div class="avatar">
      <img :src="avatar" alt="avatar" />
    </div>
    <div class="message-content-wrapper">
      <div class="message-header">
        <span class="message-sender">{{ message.sender }}</span>
        <span class="message-time">{{ message.timestamp }}</span>
      </div>
      <div class="message-content" v-html="message.content"></div>
    </div>
    <div v-if="showMenu" class="context-menu" :style="{ top: menuPosition.y + 'px', left: menuPosition.x + 'px' }">
      <button @click="handleAction('delete')">删除</button>
      <button @click="handleAction('copy')">复制</button>
      <button @click="handleAction('select')">多选</button>
    </div>
  </div>
</template>

<script>
export default {
  props: {message: {
      type: Object,
      required: true
    },
    avatar: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      showMenu: false,
      menuPosition: { x: 0, y: 0 }
    };
  },
  methods: {
    showContextMenu(event) {
      this.menuPosition = { x: event.clientX, y: event.clientY };
      this.showMenu = true;
      document.addEventListener('click', this.hideContextMenu);
      document.addEventListener('contextmenu', this.hideContextMenu);
    },
    hideContextMenu(event) {
      // 检查点击是否在菜单内，如果是则不隐藏菜单
      if (this.$el.contains(event.target)) {
        return;
      }
      this.showMenu = false;
      document.removeEventListener('click', this.hideContextMenu);
      document.removeEventListener('contextmenu', this.hideContextMenu);
    },
    handleAction(action) {
      this.$emit('message-action', action, this.message);
      this.hideContextMenu();
    }
  },
  beforeDestroy() {
    document.removeEventListener('click', this.hideContextMenu);
    document.removeEventListener('contextmenu', this.hideContextMenu);
  }
};
</script>

<style scoped>
.message-item {
  display: flex;
  align-items: flex-start;
  padding: 5px;
  position: relative;
}
.avatar {
  align-self:flex-end;
}
.avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.message-content-wrapper {
  flex: 1;
  max-width: 250px;
  display: inline-flex;
  flex-direction: column;
}
.message-header {
  flex:1;
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}
.message-sender {
  color: #888;
  font-size: 0.8rem;
  text-align: left;
}
.message-time {
  color: #888;
  font-size: 0.8rem;
  text-align: right;
}
.message-content {
  flex:5;
  background-color: #75baeb;
  padding: 10px;
  border-radius: 5px;
  text-align: left;
}
.context-menu {
  position: absolute;
  background: white;
  border: 1px solid #ccc;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}
.context-menu button {
  display: block;
  width: 100%;
  padding: 5px 10px;
  text-align: left;
  background: none;
  border: none;
  cursor: pointer;
}
.context-menu button:hover {
  background: #f0f0f0;
}
</style>