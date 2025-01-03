<template>
  <div class="message-item" >
    <div v-if="this.$store.state.user.id !== message.send_account_id" class="friend-message">
      <div class="avatar">
        <img :src="message.avatar" alt="avatar" @click="showProfileCard($event)"/>
      </div>
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ message.create_time }}</span>
        </div>
        <!--文本消息-->
        <div class="message-content" 
            v-if="message.type==='text'"
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event, message)">
        </div>
        <!--文件消息-->
        <div class="message-file" v-else>
          <div class="file-item">
            <img src="@/assets/images/default-file.png" alt="file"/>
            <div class="file-header">
              <div class="file-name">{{ message.content.name }}</div>
              <span class="file-size">{{ message.content.size }}</span>
            </div>
          </div>
          <span class="file-buttons">
            <button @click="downloadFile">下载</button>
            <button>预览</button>
            <!-- <a ref="link" style="visibility: hidden" :href="message.content" download>下载</a> -->
          </span>
        </div>
      </div>
    </div>

    <div v-else class="my-message">
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ message.create_time }}</span>
        </div>
        <!--文本消息-->
        <div class="message-content" 
             v-if="message.type==='text'"
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event, message)">
        </div>
        <!--文件消息-->
        <div class="message-file" v-else>
          <div class="file-item">
            <img src="@/assets/images/default-file.png" alt="file"/>
            <div class="file-name">{{ message.content.name }}</div>
          </div>
          <span class="file-size">{{ message.content.size }}</span>
          <span>
            <a :href="message.content" download>下载</a>
            <button>预览</button>
          </span>
          <!-- <a :href="message.content" download></a> -->
        </div>
      </div>
      <div class="avatar">
        <img :src="message.avatar" alt="avatar" @click="showProfileCard($event)"/>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  props: ['message'],
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
    downloadFile(){
      const blob = new Blob([this.message.content], { type: 'application/octet-stream' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', this.message.content.name);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url); // 释放 URL 对象
    },
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
.message-file{
  flex:5;
  background-color: #75baeb;
  padding: 10px;
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
}
.file-item{
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  flex-direction: row;
}
.message-file img{
  width: 50px;
  height: 50px;
}
.file-header{
  display: flex;
  flex-direction: column;
}
.file-name{
  margin-top: 5px;
  font-size: 0.8rem;
  color: #888;
}
.file-size{
  margin-top: 5px;
  font-size: 0.8rem;
  color: #888;
}
.file-buttons{
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  width: 100%;
  padding: 5px;
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