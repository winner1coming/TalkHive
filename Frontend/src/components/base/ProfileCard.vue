<template>
  <div v-if="visible" class="profile-card" :style="{ top: `${y}px`, left: `${x}px` }">
    <div class="options-button" @click="handleOptionsClick($event)">···</div>
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
    <button v-show="profile.is_friend" @click="sendMessage">发信息</button>
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
    <DevideMove
      :divides = "divides"
      v-show="isDevideMoveVisible"
      @divide-move="divideMove"
      @close="isDevideMoveVisible = false"
      ref="divideMove"
    />
  </div>
</template>

<script>
import { EventBus } from '@/components/base/EventBus';
import ContextMenu from '@/components/base/ContextMenu';
import DevideMove from '@/components/Contact_list/DevideMove.vue';
import * as chatListAPI from '@/services/chatList';
import * as contactListAPI from '@/services/contactList';
export default {
  components: {
    ContextMenu,
    DevideMove,
  },
  data() {
    return {
      visible: false,
      x: 0,
      y: 0,
      profile: null,
      type: 'friends',  // 'friends' or 'groups'
      isDevideMoveVisible: false,
      divides: [],
      boundD: null,
      boundR: null,
    };
  },
  methods: {
    handleOptionsClick(event){
      let items = [];
      if(!this.profile.is_friend){
        items.push('添加好友');
        this.$refs.contextMenu.show(event, items , null, this.boundD, this.boundR);
        return;
      }
      items.push('更改分组');
      if(this.profile.is_pinned){
        items.push('取消置顶');
      }else{
        items.push('置顶');
      }
      if(this.profile.is_blocked){
        items.push('取消屏蔽');
      }else{
        items.push('屏蔽');
      }
      if(this.profile.is_blacklist){
        items.push('取消拉黑');
      }else{
        items.push('拉黑');
      }
      if(this.profile.is_mute){
        items.push('取消消息免打扰');
      }else{
        items.push('消息免打扰');
      }
      items.push('删除好友');
      this.$refs.contextMenu.show(event, items , null, this.boundD, this.boundR);
    },
    async handleMenuSelect(item){
      if(item === '添加好友'){
        try{
          // 添加弹窗 todo
          const response = await contactListAPI.addFriend(this.profile.account_id);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '更改分组'){
        try{
          const response = await contactListAPI.getDivides();
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.divides = response.data.divides;
        }catch(e){
          console.log(e);
        }
        this.isDevideMoveVisible = true;
        this.$refs.divideMove.selectedDevide = this.profile.tag;
        this.$refs.divideMove.multiple = false;
      }
      else if(item === '置顶'){
        try{
          const response = await chatListAPI.pinChat(this.profile.account_id, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-pin', this.profile.account_id, true);
            this.profile.is_pinned = true;
          }
        }catch(e){  
          console.log(e);
        }
      }
      else if(item === '取消置顶'){
        try{
          const response = await chatListAPI.pinChat(this.profile.account_id, false);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-pin', this.profile.account_id, false);
            this.profile.is_pinned = false;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '屏蔽'){
        try{
          const response = await chatListAPI.blockChat(this.profile.account_id, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-blocked', this.profile.account_id, true);
            this.profile.is_blocked = true;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '取消屏蔽'){
        try{
          const response = await chatListAPI.blockChat(this.profile.account_id, false);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-blocked', this.profile.account_id, false);
            this.profile.is_blocked = false;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '消息免打扰'){
        try{
          const response = await chatListAPI.setMute(this.profile.account_id, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-mute', this.profile.account_id, true);
            this.profile.is_mute = true;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '取消消息免打扰'){
        try{
          const response = await chatListAPI.setMute(this.profile.account_id, false);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-mute', this.profile.account_id, false);
            this.profile.is_mute = false;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '拉黑'){
        try{
          const response = await contactListAPI.addToBlackList(this.profile.account_id);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-blacklist', this.profile.account_id, true);
            this.profile.is_blacklist = true;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '取消拉黑'){
        try{
          const response = await contactListAPI.removeFromBlackList(this.profile.account_id);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-blacklist', this.profile.account_id, false);
            this.profile.is_blacklist = false;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '删除好友'){
        try{
          const response = await contactListAPI.deleteFriend(this.profile.account_id);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            // todo
            this.profile
          }
        }catch(e){
          console.log(e);
        }
      }
    },
    async divideMove(divide){
      try{
        const response = await contactListAPI.moveInDevide(this.type,this.profile.account_id, divide);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        else{
          this.profile.tag = divide;
        }
      }catch(e){
        console.log(e);
      }
    },
    async sendMessage() {
      this.hide();
      try{
        const response = await chatListAPI.getChat(this.profile.tid);
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
      this.boundD = boundD;
      this.boundR = boundR;
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide() {
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
  },
  mounted() {
    // EventBus.on('other-float-component', (component) => {
    //   if (this.visible && component !== this) {
    //     this.hide();
    //   }
    // });
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
.options-button {
  position: absolute;
  top: 10px;
  right: 10px;
  cursor: pointer;
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