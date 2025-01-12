<template>
  <div v-if="visible" class="profile-card" :style="{ top: `${y}px`, left: `${x}px` }">
    <div class="options-button" @click="handleOptionsClick($event)">···</div>
    <div class="avatar">
      <img :src="profile.avatar" alt="avatar" />
    </div>
    <div class="info">
      <div class="name">{{ profile.remark }}</div>    <!--我的备注-->
      <div class="remark">群名称: {{ profile.group_name }}</div>  
      <div class="remark">
        分组:{{ profile.tag }}
      </div>
      <div class="remark">群介绍: {{ profile.introduction }}</div>
    </div>
    <button @click="sendMessage">发信息</button>
    <!--菜单-->
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
    <!--添加好友-->
    <div class="add-modal" @click.self="hide" v-show="isAddVisible">
      <div class="add-content">
        <h2>申请理由</h2>
        <textarea 
          v-model="reason" 
          placeholder="输入申请理由.."
        />
        <button @click="add(reason)">确认</button>
      </div>
    </div>
    <!--更改分组-->
    <DivideMove
      :divides = "divides"
      v-show="isDivideMoveVisible"
      @divide-move="divideMove"
      @close="isDivideMoveVisible = false"
      ref="divideMove"
    />
  </div>
</template>

<script>
import { EventBus } from '@/components/base/EventBus';
import ContextMenu from '@/components/base/ContextMenu';
import DivideMove from '@/components/Contact_list/DivideMove.vue';
import * as chatListAPI from '@/services/chatList';
import * as contactListAPI from '@/services/contactList';
export default {
  components: {
    ContextMenu,
    DivideMove,
  },
  data() {
    return {
      visible: false,
      x: 0,
      y: 0,
      profile: null,
      type: 'friends',  // 'friends' or 'groups'
      isAddVisible: false,
      isDivideMoveVisible: false,
      divides: [],
      boundD: null,
      boundR: null,
    };
  },
  methods: {
    handleOptionsClick(event){
      let items = [];
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
      if(this.profile.is_mute){
        items.push('取消消息免打扰');
      }else{
        items.push('消息免打扰');
      }
      this.$refs.contextMenu.show(event, items , null, this.boundD, this.boundR);
    },
    async handleMenuSelect(item){
      if(item === '更改分组'){
        try{
          const response = await contactListAPI.getDivides('friends');
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.divides = response.data.divides;
        }catch(e){
          console.log(e);
        }
        this.isDivideMoveVisible = true;
        this.$refs.divideMove.selectedDivide = this.profile.tag;
        this.$refs.divideMove.multiple = false;
      }
      else if(item === '置顶'){
        try{
          const response = await chatListAPI.pinChat(this.profile.group_id, true, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-pin', this.profile.group_id, true);
            this.profile.is_pinned = true;
          }
        }catch(e){  
          console.log(e);
        }
      }
      else if(item === '取消置顶'){
        try{
          const response = await chatListAPI.pinChat(this.profile.group_id, false, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-pin', this.profile.group_id, false);
            this.profile.is_pinned = false;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '屏蔽'){
        try{
          const response = await chatListAPI.blockChat(this.profile.group_id, true, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-blocked', this.profile.group_id, true);
            this.profile.is_blocked = true;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '取消屏蔽'){
        try{
          const response = await chatListAPI.blockChat(this.profile.group_id, false, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-blocked', this.profile.group_id, false);
            this.profile.is_blocked = false;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '消息免打扰'){
        try{
          const response = await chatListAPI.setMute(this.profile.group_id, true, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-mute', this.profile.group_id, true);
            this.profile.is_mute = true;
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(item === '取消消息免打扰'){
        try{
          const response = await chatListAPI.setMute(this.profile.group_id, false, true);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          else{
            EventBus.emit('set-mute', this.profile.group_id, false);
            this.profile.is_mute = false;
          }
        }catch(e){
          console.log(e);
        }
      }
    },
    async divideMove(divide){
      try{
        const response = await contactListAPI.moveInDivide(this.type,this.profile.group_id, divide);
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
        const response = await chatListAPI.getChat(this.profile.group_id, true);
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
      this.isAddVisible = false;
      this.isDivideMoveVisible = false;
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
  background-color: var(--background-color);
  color: var(--text-color);
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
  font-size: var(--font-size-mlarge);
}

.remark {
  margin-top: 5px;
  font-size: var(--font-size-small);
  color: #666;
}

.add-modal {
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
.add-content {
  background-color: var(--background-color);
  color: var(--text-color);
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
textarea {
  margin: 10px;
  width: 80%;
  height: 70%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: none;
}

button {
  margin-top: 10px;
  padding: 5px 10px;
  border: none;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: var(--button-background-color1);
}
</style>