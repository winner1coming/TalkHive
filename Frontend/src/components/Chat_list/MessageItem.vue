<template>
  <div class="message-item" >
    <div v-if="this.$store.state.user.id !== message.send_account_id" class="friend-message">
      <div class="avatar">
        <img :src="avatar" alt="avatar" @click="showProfileCard($event)"/>

      </div>
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ message.create_time }}</span>
        </div>
        <div class="message-content" 
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event, message)">

        </div>
      </div>
    </div>

    <div v-else class="my-message">
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ message.create_time }}</span>
        </div>
        <div class="message-content" 
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event, message)">
        </div>
      </div>
      <div class="avatar">
        <img :src="avatar" alt="avatar" @click="showProfileCard($event)"/>
      </div>
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
    showContextMenu(event, message) {
      this.$emit('show-context-menu',event, message);
    },
    showProfileCard(event){
      this.$emit('show-profile-card', event, this.message.send_account_id);
    }
    
  },

};
</script>

<style scoped>
.message-item {
  display: flex;
  padding: 5px;
  position: relative;
  width: 100%;
}
.friend-message {
  align-self: flex-start; 
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  width: 100%;
}
.my-message {
  align-self: flex-end;
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  width: 100%;

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