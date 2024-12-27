<template>
  <div v-if="visible" class="group-management">
    <div >
      <button
        @click="returnTo"
        class="arrow-button"
      >
        <
      </button>
    </div>
    <div v-show="componentStatus === 'main'">
      <div>
        <SearchBar 
          :isImmidiate="false" 
          :showButton="false"
          @search="searchMember" 
          @button-click="searchMember"/>
        <div class="group-members">
          <div 
            v-for="member in displayedMembers" 
            :key="member.account_id" 
            class="member"
            @click="showProfileCard($event, member.account_id)" 
            @contextmenu="showContextMenu($event, member.account_id)"
          >
            <img :src="member.avatar" alt="avatar" class="avatar">
            <p class="remark">{{ member.group_nickname.length > this.maxChars ? member.group_nickname.slice(0, this.maxChars)+'...' : member.group_nickname }}</p>
          </div>
          <div v-if="showMoreButton" class="member" @click="showAllMembers">
            <img src="" alt="plus" class="avatar">
            <p class="remark">显示更多</p>
          </div>
          <div class="member" @click="inviteMember">
            <div class="avatar add-member">
              <span>+</span>
            </div>
            <p class="remark">邀请新成员</p>
          </div>
        </div>
      </div>
      
      <div class="group-info">
        <p class="title">群聊名称:</p>
        <p class="detail">{{ groupInfo.group_name }}</p>
        <p class="title">群介绍:</p>
        <p class="detail">{{ groupInfo.introduction }}</p>
        <p class="title">群聊备注: </p>
        <EditableText class="detail" :text="group_remark" @update-text="changeGroupRemark" />
        <p class="title">我的群昵称: </p>
        <EditableText class="detail" :text="groupInfo.my_group_nickname" @update-text="changeGroupNickname" />
        <hr class="divider" />
        <!-- <p>是否显示群成员昵称: <SwitchButton v-model="groupInfo.showNicknames" /></p> -->
        <p class="title">是否消息免打扰: <SwitchButton v-model="isMute" @change-value="setMute"/></p>
        <p class="title">是否屏蔽: <SwitchButton v-model="isBlocked" @change-value="setBlock"/></p>
        <p class="title">是否置顶: <SwitchButton v-model="isPinned" @change-value="setPin"/></p>
        <hr class="divider" />
        <p class="flex-container" @click="viewChatHistory">
          <span class="title">聊天记录: </span>
          <span class="arrow-button" >></span>
        </p>
        <hr class="divider" />
        <p v-show="groupInfo.my_group_role==='group_owner'||groupInfo.my_group_role==='group_manager'" class="flex-container" @click="manageGroups">
          <span class="title">管理员设置: </span>
          <span class="arrow-button" >></span>
        </p>
        <hr v-show="groupInfo.my_group_role==='group_owner'||groupInfo.my_group_role==='group_manager'" class="divider" />
      </div>
      <div class="group-actions">
        <button @click="exitGroup">退出群聊</button>
        <button @click="hide">关闭</button>
      </div>
    </div>
    <div v-show="componentStatus === 'history'">
      <p>聊天记录</p>
      <SearchBar 
          :isImmidiate="false" 
          :showButton="false"
          @search="searchHistory" 
      />
      <div class="history-list">
        <div
          v-for="message in filteredHistory"
          :key="message.message_id"
          class="message-item"
        >
          <div class="message-header">
            <img :src="message.avatar" alt="avatar" />
            <div>
              <p class="message-sender">{{ message.sender }}</p>
              <p class="message-time">{{ message.create_time }}</p>
            </div>
          </div>
          <div class="message-content">
            <p class="message-text">{{ message.content }}</p>
          </div>
        </div>
      </div>
    </div>
    <div v-show="componentStatus === 'manage'">
      <p>管理员设置</p>
      <p class="title">全体禁言: <SwitchButton v-model="groupInfo.muteAll" @change-value=""/></p>
      <p class="detail"> 已禁言的成员：</p>
      <div class="muted-members-list">
        <div 
          v-for="member in groupInfo.members.filter(member => member.is_banned)" 
          :key="member.account_id" 
          class="muted-member"
        >
          <img :src="member.avatar" alt="avatar" class="avatar">
          <span class="remark">{{ member.group_nickname }}</span>
          <button @click="setBanned(member.account_id, false)">解禁</button>
        </div>
      </div>
      <hr class="divider" />
      <p class="title">申请入群的方式：</p>
      <p class="detail" style="margin-left: 15px;">成员邀请: <SwitchButton v-model="groupInfo.allow_invite" @change-value="changeInvitePermission"/></p>
      <p class="detail" style="margin-left: 15px;">群号搜索: <SwitchButton v-model="groupInfo.allow_id_search" @change-value="changeIdPermission"/></p>
      <p class="detail" style="margin-left: 15px;">群名称搜索: <SwitchButton v-model="groupInfo.allow_name_search" @change-value="changeNamePermission"/></p>
      <hr class="divider" />
      <p class="title">更改群头像<button>上传</button></p>
    </div>
    <ProfileCard ref="profileCard" />
    <ContextMenu ref="contextMenu" @select-item="handleMenuSelect"/>
  </div>
 
</template>

<script>
import * as contactListAPI from '@/services/contactList';
import * as chatListAPI from '@/services/chatList';
import {getProfileCard} from '@/services/api';
import { EventBus } from '@/components/base/EventBus';
import EditableText from '@/components/base/EditableText.vue';
import SwitchButton from '@/components/base/SwitchButton.vue';
import { changeGroupNickname } from '../../services/contactList';
import SearchBar from '@/components/base/SearchBar.vue';
import ProfileCard from '@/components/base/ProfileCard.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
export default {
  components: {
    EditableText,
    SwitchButton,
    SearchBar,
    ProfileCard,
    ContextMenu,
  },
  data() {
    return {
      visible: false,
      group_id:'',
      group_remark:'',
      isMute: false,
      isBlocked: false,
      isPinned: false,
      showAll:false,
      // chat的信息
      // chatInfo:{
      //   tags: ['mute'],
      //   name: '111',  // 我的群名备注
      // },
      // groupInfo: {
      //   group_owner: '111',  // 群主tid
      //   introduction: '111',
      //   // 入群权限 todo
        
      //   my_group_nickname: 'aa',   // 我在本群的群昵称
      //   members: [
      //     {
      //       account_id: '111',
      //       avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
      //       group_role: 'group_owner',
      //       group_nickname: 'aa',
      //       is_banned:false,
      //     },
      //     {
      //       account_id: '222',
      //       avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
      //       group_role: 'group_owner',
      //       group_nickname: 'bb',
      //     },
      //     {
      //       account_id: '333',
      //       avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
      //       group_role: 'group_owner',
      //       group_nickname: 'cc',
      //     },
      //   ],
      //   // showNicknames: false, 需求没做 todo
      // },
      groupInfo:{
        group_name: '',
        group_owner: '',  // 群主tid
        introduction: '',
        // 入群权限 
        allow_invite: true, 
        allow_id_search: true, 
        allow_name_search: true,
        my_group_nickname: '',   // 我在本群的群昵称
        my_group_role:'',
        members: [],
      },
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
      searchKeyword:'',
      componentStatus: 'main',  // 'main', 'history', 'manage'
      boundD: null, // 边界的坐标
      boundR: null, // 边界的坐标
    };
  },
  watch: {
    "$store.state.currentChat": {
      immediate: true,
      deep: true,
      handler(newVal) {
        if(!newVal) return;
        if(this.visible && newVal.id !== this.group_id){
          this.hide();
        }
        this.group_id = newVal.id;
        this.group_remark = newVal.name;
        this.isMute = newVal.tags.includes('mute');
        this.isBlocked = newVal.tags.includes('blocked');
        this.isPinned = newVal.tags.includes('pinned');
      }
    }
  },
  methods: {
    async fetchGroupInfo(){
      try{
        const response = await contactListAPI.getGroupInfo(this.group_id);
        if(response.status === 200){
          this.groupInfo = response.data;
        }else{
          // todo
          console.log(response.data.message);
        }
      }
      catch(error){
        console.log('fetch group error:', error);
      }
      
    },
    returnTo(){
      if(this.componentStatus === 'main'){
        this.hide();
      }else{
        this.componentStatus = 'main';
      }
    },
    async searchMember(key){
      if(key === ''){
        this.fetchGroupInfo();
        return;
      }
      try{
        const response = await contactListAPI.searchGroupMember(key);
        if(response.status === 200){
          this.groupInfo.members = response.data;
        }
        else{
          // todo
        }
      }
      catch(error){
        console.log('search member error:', error);
      }
    },
    async showProfileCard(event, account_id){
      try{
        const response = await getProfileCard(account_id);
        if(response.status === 200){
          const profile = response.data;
          this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
        }
        else{
          // todo
        }
      }
      catch(error){
        console.log('show profile card error:', error);
      }
    },
    async changeGroupRemark(newRemark){
      try{
        const response = await contactListAPI.changeRemark(this.group_id, newRemark);
        if (response.status === 200) {
          this.group_remark = newRemark;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.name = newRemark;
          this.$store.dispatch('setChat', chatInfo); // 更新chatList
        } else {
          console.log(response.data.message);
        }
      }
      catch(error){
        console.log('change group remark error:', error)
      }
    },
    async changeGroupNickname(newNickname){
      try {
        const response = await contactListAPI.changeGroupNickname(this.group_id, newNickname);
        if (response.status === 200) {
          this.groupInfo.my_group_nickname = newNickname;
        } else {
          console.log(response.data.message);
        }
      } catch (error) {
        console.log('change group nickname error:', error);
      }
    },
    async setMute(){
      try{
        const response = await chatListAPI.setMute(this.group_id, !this.isMute);
        if(response.status === 200){
          this.isMute = !this.isMute;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isMute ? [...chatInfo.tags, 'mute'] : chatInfo.tags.filter(tag => tag !== 'mute');
          this.$store.dispatch('setChat',chatInfo); // 更新chatList
        }else{
          // todo
        }
      }
      catch(error){
        // todo  滚动条需要设回去
        console.error('Failed to set mute:', error);
      }
    },
    async setBlock(){
      try{
        const response = await chatListAPI.blockChat(this.group_id, !this.isBlocked);
        if (response.status === 200) {
          this.isBlocked = !this.isBlocked;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isBlocked ? [...chatInfo.tags, 'blocked'] : chatInfo.tags.filter(tag => tag !== 'blocked');
          this.$store.dispatch('setChat', chatInfo);  // 更新chatList
        }else{
          // todo
        }
      }
      catch(error){
        // todo
        
      }
    },
    async setPin() {
      try {
        const response = await chatListAPI.pinChat(this.group_id, !this.isPinned);
        if (response.status === 200) {
          this.isPinned = !this.isPinned;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isPinned ? [...chatInfo.tags, 'pinned'] : chatInfo.tags.filter(tag => tag !== 'pinned');
          this.$store.dispatch('setChat', chatInfo);  // 更新chatList
        } else {
          // todo
        }
      } catch (error) {
        // todo
        console.error('Failed to set pin:', error);
      }
    },
    async setBanned(account_id, is_banned){
      try {
        const response = await chatListAPI.setBanned(this.group_id,this.account_id, this.is_banned);
        if (response.status === 200) {
          let member = this.groupInfo.members.find(member => member.account_id === account_id);
          if (member) {
            member.is_banned = is_banned;
          }
        } else {
          // todo
        }
      } catch (error) {
        // todo
        console.error('Failed to set pin:', error);
      }
    },
    async exitGroup(){
      // 退出群聊
      try{
        const response = await contactListAPI.exitGroup(this.group_id);
        if(response.status === 200){
          this.$emit('group-exited', this.group_id);
          this.hide();
        }
        else{
          // todo
          console.log(response.data.message);
        }
      }
      catch(error){
        console.log('exit group error:', error);
      }
    },
    async addMember() {
      try {
        await addMemberToGroup(this.groupInfo.id, this.newMemberId);
        this.groupInfo.members.push({ id: this.newMemberId, name: '新成员' }); // 假设新成员的名字为 '新成员'
        this.newMemberId = '';
      } catch (error) {
        console.error('Failed to add member:', error);
      }
    },
    async removeMember(memberId) {
      try {
        await removeMemberFromGroup(this.groupInfo.id, memberId);
        this.groupInfo.members = this.groupInfo.members.filter(member => member.id !== memberId);
      } catch (error) {
        console.error('Failed to remove member:', error);
      }
    },
    async deleteGroup() {
      try {
        await deleteGroup(this.groupInfo.id);
        this.$emit('group-deleted', this.groupInfo.id);
        this.hide();
      } catch (error) {
        console.error('Failed to delete group:', error);
      }
    },
    async viewChatHistory() {
      // 查看聊天记录
      this.componentStatus = 'history';
      chatListAPI.getHistory(this.group_id).then(response => {
        if (response.status === 200) {
          this.history = response.data;
        } else {
          // todo
          console.log(response.data.message); 
        }
      }).catch(error => {
        console.log('get chat history error:', error);
      });
    },
    searchHistory(keyword){
      this.searchKeyword = keyword;
    },
    manageGroups() {
      // 管理员设置
      this.componentStatus = 'manage';
    },
    async changeInvitePermission(){
      try{
        const response = await contactListAPI.changeInvitePermission(this.group_id, this.groupInfo.allow_invite);
        if(response.status !== 200){
          // todo
        }
      }
      catch(error){
        console.log('change invite permission error:', error);
      }
    },
    async changeIdPermission(){
      try{
        const response = await contactListAPI.changeIdPermission(this.group_id, this.groupInfo.allow_id_search);
        if(response.status !== 200){
          // todo
        }
      }
      catch(error){
        console.log('change id permission error:', error);
      }
    },
    async changeNamePermission(){
      try{
        const response = await contactListAPI.changeNamePermission(this.group_id, this.groupInfo.allow_name_search);
        if(response.status !== 200){
          // todo
        }
      }
      catch(error){
        console.log('change name permission error:', error);
      }
    },
    // 右键菜单
    showContextMenu(event, account_id) {
      if(account_id === this.groupInfo.group_owner || account_id === this.$store.state.user.id){
          return;
      }
      if(this.groupInfo.members.find(member => member.account_id === account_id).group_role === 'group_manager'){
        return;
      }
      if(this.groupInfo.my_group_role==='group_owner'){
        let items = ['移除'];
        if(this.groupInfo.members.find(member => member.account_id === account_id).is_banned){
          items.unshift('解禁');
        }else{
          items.unshift('禁言');
        }
        if(this.groupInfo.members.find(member => member.account_id === account_id).group_role === 'group_manager'){
          items.push('取消管理员');
        }else{
          items.push('设为管理员');
        }
        items.push('转让群主');
        this.$refs.contextMenu.show(event, items, account_id, this.boundD, this.boundR);
      }
      else if(this.groupInfo.my_group_role==='group_manager'){
        let items = ['移除'];
        if(this.groupInfo.members.find(member => member.account_id === account_id).is_banned){
          items.unshift('解禁');
        }else{
          items.unshift('禁言');
        }
        this.$refs.contextMenu.show(event, items, account_id, this.boundD, this.boundR);
      }
    },
    async handleMenuSelect(option, account_id){
      if(option==='禁言'){
        try{
          const response = await chatListAPI.setBanned(this.group_id, account_id, true);
          if(response.status === 200){
            let member = this.groupInfo.members.find(member => member.account_id === account_id);
            if(member){
              member.is_banned = true;
            }
          }
          else{
            // todo
            console.log('set banned error:', response.data.message);
          }
        }
        catch(error){
          console.log('set banned error:', error);
        }
      }
      else if(option==='解禁'){
        try{
          const response = await chatListAPI.setBanned(this.group_id, account_id, false);
          if(response.status === 200){
            let member = this.groupInfo.members.find(member => member.account_id === account_id);
            if(member){
              member.is_banned = false;
            }
          }
          else{
            console.log('set banned error:', response.data.message);
          }
        }
        catch(error){
          console.log('set banned error:', error);
        }
      }
      else if(option==='移除'){
        try{
          const response = await chatListAPI.removeMember(this.group_id, account_id);
          if(response.status === 200){
            this.groupInfo.members = this.groupInfo.members.filter(member => member.account_id !== account_id);
          }
          else{
            console.log('remove member error:', response.data.message);
          }
        }
        catch(error){
          console.log('remove member error:', error);
        }
      }
    },

    // 显示与隐藏
    show(){
      this.fetchGroupInfo();
      this.visible = true;
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide(){
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
  },
  computed:{
    maxChars(){  // 可以显示的字体个数
      return Math.floor(108.0 / parseInt(this.$store.state.settings.fontSize,10)* 0.6);
    },
    filteredHistory(){
      const keyword = this.searchKeyword;
      if(!keyword) return this.history;
      console.log(keyword);
      return this.history.filter(message => {
        return (
          message.sender.includes(keyword) ||
          message.content.includes(keyword)
        );
      });
    },
    displayedMembers() {
      return this.showAll ? this.groupInfo.members : this.groupInfo.members.slice(0, 20);
    },
    showMoreButton() {
      return !this.showAll && this.groupInfo.members.length > 20;
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
};
</script>

<style scoped src="@/assets/css/chatList.css"></style>
<style scoped>
.group-management {
  width: 200px;
  padding: 10px;
  background-color: #f6f1f1;
  border: 1px solid #ccc;
  border-radius: 5px;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  overflow-y: auto;
}

.arrow-button {
  background-color: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
  margin: 0;
  padding: 0;
}

.group-members {
  display: flex;
  flex-wrap: wrap;
}
.member {
  margin: 5px;
  text-align: center;
  width: 35px;
  height: 75px;
}
.remark {
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 0.7rem
}
.avatar {
  width: 35px;
  height: 35px;
  border-radius: 50%;
}


.group-info {
  margin-top: 20px;
  align-self: flex-start;
}
.title {
  color: black;
  text-align: left;
  font-weight: 500;
}
.detail {
  text-align: left;
  color: #888;
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
  background: #e0e0e0;
  margin: 10px 0;
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
  padding: 10px;
  border-bottom: 1px solid #e0e0e0;
}
.message-header img {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 10px;
}
.message-content {
  flex-grow: 1;
}
.message-sender {
  font-weight: bold;
  margin-bottom: 5px;
}
.message-text {
  margin-bottom: 5px;
}
.message-time {
  font-size: 0.8em;
  color: #888;
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