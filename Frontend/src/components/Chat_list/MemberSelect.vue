<template>
  <div v-if="visible" class="profile-card" :style="{ top: `${y}px`, left: `${x}px` }">
    <div class="avatar">
      <img :src="profile.avatar" alt="avatar" />
    </div>
    <div class="info">
      <div class="name">{{ profile.remark }}</div>    <!--我的备注-->
      <div class="remark">昵称: {{ profile.nickname }}</div>   
      <div class="remark" v-show="profile.groupNickname">群昵称: {{ profile.groupNickname }}</div>
      <div class="remark">
        分组:{{ profile.tag }}
      </div>
      <div class="remark">个性签名: {{ profile.signature }}</div>
    </div>
    <button @click="sendMessage">发信息</button>
  </div>
</template>

<script>
import { EventBus } from '@/components/base/EventBus';
import { getChat } from '@/services/chatList';
export default {
  props:['members'],
  data() {
    return {
      visible: false,
      x: 0,
      y: 0,
      profile: null,
      type: 'friend',  // 'friend' or 'group'
    };
  },
  methods: {
    show(event, profile, boundD, boundR) {  // boundD, boundR 为边界的坐标
      EventBus.emit('float-component-opened', this); // 通知其他组件
      const cardWidth = 200;
      const cardHeight = 400;
      const x = event.clientX + cardWidth > boundR ? event.clientX - cardWidth : event.clientX;
      const y = event.clientY + cardHeight > boundD ? boundD - cardHeight : event.clientY;
      this.x = x;
      this.y = y;
      this.profile = profile;
      this.visible = true;
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide() {
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
    async sendMessage() {
      this.hide();
      try{
        const response = await getChat(this.profile.tid, false);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.$store.dispatch('setChat', response.data.data);
        this.$router.push({name: 'chat'});
      }catch(e){
        console.log(e);
      }
      
      
    },
  },
  mounted() {
    EventBus.on('other-float-component', (component) => {
      if (this.visible && component !== this) {
        this.hide();
      }
    });
    EventBus.on('close-float-component', (clickedElement) => {
      if (this.visible && !this.$el.contains(clickedElement)) {
        this.hide();
      }
    });
  },
};
</script>

<style scoped>
.profile-card {
  position: absolute;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 10px;
  padding-left: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  width: 200px;
  background-color: #fff;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.avatar img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
}

.info {
  margin-top: 10px;
  text-align: center;
}

.name {
  font-weight: bold;
  font-size: 1.2rem;
}

.remark {
  margin-top: 5px;
  font-size: 0.9rem;
  color: #666;
}

button {
  margin-top: 10px;
  padding: 5px 10px;
  border: none;
  background-color: #007bff;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}
</style>