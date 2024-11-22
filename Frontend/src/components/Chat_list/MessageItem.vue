<template>
  <div class="message-item" @contextmenu.prevent="showContextMenu($event)">
    <div class="message-header">
      <span class="message-sender">{{ message.sender }}</span>
      <span class="message-time">{{ message.timestamp }}</span>
    </div>
    <div class="message-content" v-html="message.content"></div>
    <div v-if="showMenu" class="context-menu" :style="{ top: menuPosition.y + 'px', left: menuPosition.x + 'px' }">
      <button @click="handleAction('delete')">删除</button>
      <button @click="handleAction('copy')">复制</button>
      <button @click="handleAction('select')">多选</button>
    </div>
  </div>
</template>

<script>
export default {
  props: ['message'],
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
    },
    hideContextMenu() {
      this.showMenu = false;
      document.removeEventListener('click', this.hideContextMenu);
    },
    handleAction(action) {
      this.$emit('message-action', action, this.message);
      this.hideContextMenu();
    }
  }
};
</script>

<style scoped>
.message-item {
  padding: 10px;
  border-bottom: 1px solid #eee;
  position: relative;
}
.message-header {
  display: flex;
  justify-content: space-between;
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