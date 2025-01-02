<template>
  <div v-if="visible" class="group-management">
    <div style="width: 100%;">
      <p
        @click="returnTo"
        class="arrow-button"
      >
        <
      </p>
    </div>
    <!--主页面-->
    <div v-show="componentStatus === 'main'">
      <!--好友信息-->
      <div class="group-info">
        <p class="title">好友名称:</p>
        <p class="detail">{{ friendInfo.group_name }}</p>
        <p class="title">签名:</p>
        <p class="detail">{{ friendInfo.introduction }}</p>
        <p class="title">备注: </p>
        <EditableText class="detail" :text="group_remark" @update-text="changeFriendRemark" />
        <p class="title">分组: </p>
        <EditableText class="detail" :text="friendInfo.my_group_nickname" @update-text="changeGroupNickname" />
        <hr class="divider" />
        <p class="title">是否消息免打扰: <SwitchButton v-model="isMute" @change-value="setMute"/></p>
        <p class="title">是否屏蔽: <SwitchButton v-model="isBlocked" @change-value="setBlock"/></p>
        <p class="title">是否置顶: <SwitchButton v-model="isPinned" @change-value="setPin"/></p>
        <hr class="divider" />
        <p class="flex-container" @click="viewChatHistory">
          <span class="title">聊天记录: </span>
          <span class="arrow-button" >></span>
        </p>
        <hr class="divider" />
      </div>
      <!--退出-->
      <div class="group-actions">
        <button @click="exitGroup">删除好友</button>
        <button @click="hide">关闭</button>
      </div>
    </div>
    <!--聊天记录-->
    <div v-show="componentStatus === 'history'" class="view-history">
      <p>聊天记录</p>
      <div class="search-bar" >
        <input
          type="text"
          v-model="searchHistoryKeyword"
          placeholder="搜索..."
          
        />
      </div>
      <!--搜索类型-->
      <div class="search-type">
        <button :class="{ 'type-button': true, 'active': searchHistoryType === 'all' }" @click="searchHistoryType='all'">全部</button>
        <button :class="{'type-button': true, active:searchHistoryType==='image'}" @click="searchHistoryType='image'">图片</button>
        <button :class="{'type-button': true, active:searchHistoryType==='file'}" @click="searchHistoryType='file'">文件</button>
        <button :class="{'type-button': true, active:searchHistoryType==='date'}" @click="searchHistoryType='date'">日期</button>
        <button :class="{'type-button': true, active:searchHistoryType==='member'}" @click="searchHistoryType='member'">成员</button>
      </div>
      <input
        type="date"
        v-show="searchHistoryType==='date'"
        v-model="searchHistoryDate"
        class="date-picker"
      />
      <!--筛选好的历史记录--> 
      <div v-if="filteredHistory" class="history-list">
        <div
          v-for="(message, index) in filteredHistory"
          :key="message.message_id"
          :ref="'message-' + index"
          class="message-item"
        >
          <div class="message-header">
            <img :src="message.avatar" alt="avatar" />
            <p class="message-sender">{{ message.sender }}</p>
            <p class="message-time">{{ message.create_time }}</p>
          </div>
          <div>

          </div>
          <div class="message-content">
            <p class="message-text">{{ message.content }}</p>
          </div>
        </div>
      </div>
      <div v-else-if="searchHistoryKeyword" class="no-result">
        <p>无搜索结果</p>
      </div>
      <!-- <div v-else-if="searchHistoryType==='Date'" class="no-result">
        <p>加载中...</p>
      </div> -->
      <div v-else class="no-result">
        <p>输入关键词或按类型查找</p>
      </div>
    </div>
    <ProfileCard ref="profileCard" />
    <ContextMenu ref="contextMenu" @select-item="handleMenuSelect"/>
  </div>
 
</template>
<script>
import * as contactListAPI from '@/services/contactList';
import * as chatListAPI from '@/services/chatList';
import {getPersonProfileCard} from '@/services/api';
import { EventBus } from '@/components/base/EventBus';
import EditableText from '@/components/base/EditableText.vue';
import SwitchButton from '@/components/base/SwitchButton.vue';
import SearchBar from '@/components/base/SearchBar.vue';
import ProfileCard from '@/components/base/ProfileCard.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import InviteMember from '@/components/Chat_list/InviteMember.vue';
export default {
  components: {
    EditableText,
    SwitchButton,
    SearchBar,
    ProfileCard,
    ContextMenu,
    InviteMember,
  },
  components: {
    EditableText,
    SwitchButton,
    SearchBar,
    ProfileCard,
    ContextMenu,
    InviteMember,
  },
  data() {
    return {
      visible: false,
      query: "", // 搜索关键词
      isComposing: false, // 是否正在使用输入法输入，防止频繁触发搜索
      account_id:'',
      group_remark:'',
      isMute: false,
      isBlocked: false,
      isPinned: false,
      showAll:false,
      history:[
        {
          message_id:'1',
          create_time:'12:20',
          send_account_id:'1',
          sender:'Alice',
          content:'fadaf',
          type:'text',
          avatar:'',
        },
      ],
      // 搜索历史方面
      searchHistoryKeyword:'',
      searchHistoryType:'all', // 'all', 'image', 'file', 'date','member'
      searchHistoryDate:'',
      searchHistoryMember:'',
      componentStatus: 'main',  // 'main', 'history', 'manage'
      boundD: null, // 边界的坐标
      boundR: null, // 边界的坐标
    }
  },
  watch: {
    "$store.state.currentChat": {
      immediate: true,
      deep: true,
      handler(newVal) {
        if(!newVal) return;
        if(this.visible && newVal.id !== this.account_id){
          this.hide();
        }
        this.account_id = newVal.id;
        this.group_remark = newVal.name;
        this.isMute = newVal.tags.includes('mute');
        this.isBlocked = newVal.tags.includes('blocked');
        this.isPinned = newVal.tags.includes('pinned');
      }
    }
  },
  methods: {
    async fetchFriendInfo(){
      try{
        const response = await contactListAPI.getPersonProfileCard(this.account_id);
        if(response.status === 200){
          this.friendInfo = response.data.data;
        }else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('fetch group error:', error);
      }
      
    },
    initialize(){
      
    },
    returnTo(){
      if(this.componentStatus === 'main'){
        this.hide();
      }else{
        this.initialize();
        this.componentStatus = 'main';
      }
    },
    async showProfileCard(event, account_id){
      try{
        const response = await getPersonProfileCard(account_id, this.account_id);
        if(response.status === 200){
          const profile = response.data.data;
          this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
        }
        else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('show profile card error:', error);
      }
    },
    // 对好友的个人设置
    async changeFriendRemark(newRemark){
      try{
        const response = await contactListAPI.changeRemark(this.account_id, newRemark);
        if (response.status === 200) {
          this.group_remark = newRemark;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.name = newRemark;
          this.$store.dispatch('setChat', chatInfo); // 更新chatList
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('change group remark error:', error)
      }
    },
    async setMute(){
      try{
        const response = await chatListAPI.setMute(this.account_id, !this.isMute, true);
        if(response.status === 200){
          this.isMute = !this.isMute;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isMute ? [...chatInfo.tags, 'mute'] : chatInfo.tags.filter(tag => tag !== 'mute');
          this.$store.dispatch('setChat',chatInfo); // 更新chatList
        }else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        // todo  滚动条需要设回去
        console.error('Failed to set mute:', error);
      }
    },
    async setBlock(){
      try{
        const response = await chatListAPI.blockChat(this.account_id, !this.isBlocked, true);
        if (response.status === 200) {
          this.isBlocked = !this.isBlocked;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isBlocked ? [...chatInfo.tags, 'blocked'] : chatInfo.tags.filter(tag => tag !== 'blocked');
          this.$store.dispatch('setChat', chatInfo);  // 更新chatList
        }else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        // todo
        
      }
    },
    async setPin() {
      try {
        const response = await chatListAPI.pinChat(this.account_id, !this.isPinned, true);
        if (response.status === 200) {
          this.isPinned = !this.isPinned;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isPinned ? [...chatInfo.tags, 'pinned'] : chatInfo.tags.filter(tag => tag !== 'pinned');
          this.$store.dispatch('setChat', chatInfo);  // 更新chatList
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      } catch (error) {
        // todo
        console.error('Failed to set pin:', error);
      }
    },
    async setBanned(account_id, is_banned){
      try {
        const response = await contactListAPI.setBanned(this.account_id,this.account_id, this.is_banned);
        if (response.status === 200) {
          let member = this.groupInfo.members.find(member => member.account_id === account_id);
          if (member) {
            member.is_banned = is_banned;
          }
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      } catch (error) {
        // todo
        console.error('Failed to set pin:', error);
      }
    },
    async deleteFriend(){
      try{
        const response = await contactListAPI.deleteFriend(this.account_id);
        if(response.status === 200){
          this.$root.notify('删除成功', 'success');
          this.hide();
          this.$store.dispatch('deleteChat', this.account_id);
        }else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('delete friend error:', error);
      }
    },
    async viewChatHistory() {
      // 查看聊天记录
      this.componentStatus = 'history';
      chatListAPI.getHistory(this.group_id).then(response => {
        if (response.status === 200) {
          this.history = response.data.data;
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      }).catch(error => {
        console.log('get chat history error:', error);
      });
    },
    searchHistory(keyword){
      this.searchHistoryKeyword = keyword;
    },
    scrollToMessage(index) {
      this.$nextTick(() => {
        const messageElement = this.$refs['message-' + index];
        if (messageElement && messageElement[0]) {
          messageElement[0].scrollIntoView({ behavior: 'smooth' });
        }
      });
    },

    // 显示与隐藏
    show(){
      this.fetchGroupInfo();
      this.visible = true;
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide(){
      this.visible = false;
      this.componentStatus = 'main';
      initialize();
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
  },
  computed:{
    // maxChars(){  // 可以显示的字体个数
    //   return Math.floor(40.0 / parseInt(this.$store.state.settings.fontSize,10) / 0.6);
    // },
    filteredHistory(){
      const keyword = this.searchHistoryKeyword;
      if(this.searchHistoryType === 'all'){
        if(!keyword) return null;
        return this.history.filter(message => {
          return (
            message.type === 'text' &&
            message.content.includes(keyword)
          );
        });
      }
      else if(this.searchHistoryType === 'image'){  // todo
        return this.history.filter(message => {
          return message.type === 'image' && message.content.includes(keyword);
        });
      }
      else if(this.searchHistoryType === 'file'){  // todo
        return this.history.filter(message => {
          return message.type === 'file' && message.content.includes(keyword);
        });
      }
      else if(this.searchHistoryType === 'date'){ 
        if(!this.searchHistoryDate) return null;
        if(!this.searchHistoryKeyword){
          this.scrollToMessage(this.history.findIndex(message => message.create_time.includes(this.searchHistoryDate)));
          return this.history;
        }
        return this.history.filter(message => {
          return message.create_time.includes(this.searchHistoryDate) &&
            message.content.includes(keyword)
        });
      }
      else if(this.searchHistoryType === 'member'){
        return this.history.filter(message => {
          return message.send_account_id===this.searchHistoryMember 
            && message.content.includes(keyword);
        });
      }
    },
  },
  created(){
    //this.fetchGroupInfo();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
  mounted() {
    EventBus.on('close-float-component', (clickedElement) => {
      if (this.visible && !this.$el.contains(clickedElement)) {
        console.log(this.$el);
        this.hide();
      }
    });
  },
  beforeDestroy() {
    EventBus.off('close-float-component');
  }
}
</script>

<style scoped src="@/assets/css/chatList.css"></style>
<style scoped>
.group-management {
  width: 300px;
  background-color: #f6f1f1;
  border: 1px solid #ccc;
  border-radius: 5px;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  overflow-y: auto;
  overflow-x: hidden;
}

.arrow-button {
  background-color: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
  margin: 0;
  padding: 0;
  text-align: left;
}

/* 固定在首尾 */
.sticky-top {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: #f6f1f1;
}
.sticky-bottom {
  background-color: #f6f1f1;
  position: sticky;
  bottom: 0px;
  z-index: 10;
}

.view-history {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}
.search-bar {
  display: flex;
  padding: 10px;
  width: 90%;
}
.search-bar input {
  flex: 1;
  padding: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.group-members {
  display: flex;
  flex-wrap: wrap;
}
.member {
  margin: 5px;
  text-align: center;
  width: 40px;
  height: 60px;
}
.remark {
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 0.6rem
}
.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.show-member-hint {
  text-align: center;
  color: #888;
  cursor: pointer;
  height: 30px;
}

.group-info {
  margin-top: 20px;
  align-self: flex-start;
}
.title {
  color: black;
  text-align: left;
  font-weight: 500;
  padding: 5px;
}
.detail {
  text-align: left;
  color: #888;
  padding: 5px;
}

.flex-container{
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-actions {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
}

.group-actions input {
  margin-right: 10px;
}
.group-actions button {
  margin-top: 10px;
}

.divider {
  border: 0;
  height: 1px;
  width: 100%;
  background: #e0e0e0;
  margin: 10px 0;

}

.search-members{
  display: flex;
  flex-wrap: wrap;
  align-content: flex-start;
  height: 600px;
  width: 100%;
  overflow-y: auto;
}
.no-result {
  text-align: center;
  color: #888;
  margin-top: 20px;
  width: 100%;
}

.search-type{
  display: flex;
  justify-content: space-around;
  padding: 5px;
}
.type-button{
  padding: 5px;
  border: none;
  cursor: pointer;
  background-color: transparent;
}
.type-button:hover{
  background-color: transparent;
}
.type-button.active{
  color: #7184da;
  background-color: transparent;
}
.date-picker {
  display: block;
  padding: 5px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.date-picker:focus {
  border-color: #007bff;
  box-shadow: 0 0 5px rgba(0, 123, 255, 0.5);
}
.history-list {
  max-height: 500px; 
  overflow-y: auto;
  border: 1px solid #e0e0e0; 
  border-radius: 5px;
  background-color: #f9f9f9; 
}
.message-item {
  display: flex;
  align-items: flex-start;
  border: 1px solid #ccc; 
  border-radius: 5px; 
  background-color: #fff; 
  flex-direction: column;
}
.message-header {
  display: flex;
  align-items: center;
  border-bottom: 1px solid #e0e0e0;
  width: 100%;
  font-size: 0.8em;
  color: #888;
  justify-content: space-between;
}
.message-header img {
  width: 35px;
  height: 35px;
  border-radius: 50%;
  margin-right: 10px;
  padding: 3px;
}
.message-content {
  flex-grow: 1;
}
.message-text {
  margin-bottom: 5px;
  text-align: left;
  padding: 3px;
}

.muted-members-list {
  max-height: 300px; 
  overflow-y: auto; 
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  background-color: #f9f9f9; 
}
.muted-member {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  border: 1px solid #ccc; 
  border-radius: 5px; 
  background-color: #fff; 
}
.muted-member .avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 10px;
}

</style>
