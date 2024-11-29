<template>
  <div class="message-item" 
    
  >
    <div v-if="this.$store.state.user.id === message.send_account_id" class="friend-message">
      <div class="avatar">
        <img :src="avatar" alt="avatar" />
      </div>
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ message.timestamp }}</span>
        </div>
        <div class="message-content" 
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event)">
        </div>
      </div>
    </div>

    <div v-else class="my-message">
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ message.timestamp }}</span>
        </div>
        <div class="message-content" 
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event)">
        </div>
      </div>
      <div class="avatar">
        <img :src="avatar" alt="avatar" />
      </div>
    </div>
    
    <div v-show="showMenu" 
      class="context-menu" 
      :style="{ top: `${axis.y}px`, left: `${axis.x}px` }"
    >
      <button @click="handleAction('reply')">回复</button>
      <button @click="handleAction('forward')">转发</button>
      <button @click="handleAction('delete')">删除</button>
    </div>
  </div>
</template>

<script>
export default {
  props: ['message', 'avatar'],
  data() {
    return {
      showMenu: false,
      axis: {
        x: 0,
        y: 0
      },
    };
  },
  methods: {
    showContextMenu(event) {
      var x = event.clientX;
      var y = event.clientY;
      this.axis = {
        x,
        y
      };
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
  min-height: 50px;
  padding: 5px;
  position: relative;
}
.friend-message {
  float: left;
  display: flex;
  align-items: flex-start;
}
.my-message {
  float: right;
  display: flex;
  align-items: flex-end;
}
.friend-message .avatar {
  align-self: flex-start;
}
.my-message .avatar {
  align-self: flex-end;
}
.avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.message-content-wrapper {
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