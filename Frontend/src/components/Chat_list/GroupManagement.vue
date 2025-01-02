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
      <div :class="{'search-bar':true, 'sticky-top':this.showAll}" >
        <input
          type="text"
          v-model="query"
          placeholder="搜索群成员..."
          @compositionstart="isComposing = true"
          @compositionend="isComposing = false;triggerSearch()"
          @input="triggerSearch"
        />
      </div>
      <!--群成员-->
      <div class="group-members">
        <div 
          v-for="member in displayedMembers" 
          :key="member.account_id" 
          class="member"
          @click="showProfileCard($event, member.account_id)" 
          @contextmenu="showContextMenu($event, member.account_id)"
        >
          <img :src="member.avatar" alt="avatar" class="avatar">
          <p class="remark">{{ member.remark? member.remark : (member.group_nickname?member.group_nickname:member.nickname)}}</p>
        </div>
        <!--邀请成员-->
        <div class="member" @click="inviteMember">
          <div>
            <img src="@/assets/images/plus.png" alt="plus" class="avatar">
          </div>
        </div>
      </div>
      <div v-if="showMoreButton" @click="this.showAll = true;">
        <p class="show-member-hint">显示更多</p>
      </div>
      <div v-else @click="this.showAll = false;" class="sticky-bottom">
        <p class="show-member-hint">收起</p>
      </div>
      <!--群聊信息-->
      <div class="group-info">
        <p class="title">群聊名称:</p>
        <p class="detail">{{ groupInfo.group_name }}</p>
        <p class="title">群介绍:</p>
        <p class="detail">{{ groupInfo.introduction }}</p>
        <p class="title">群聊备注: </p>
        <EditableText class="detail" :text="group_remark" @update-text="changeGroupRemark" />
        <p class="title">我的群昵称: </p>
        <EditableText class="detail" :text="groupInfo.my_group_nickname" @update-text="changeGroupNickname" />
        <p class="title">分组: </p>
        <p class="detail">
          {{ groupInfo.divide }} 
          <button @click="showDivideMove">更改</button>
        </p>
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
      <!--群聊设置-->
      <div class="group-actions">
        <button @click="exitGroup">退出群聊</button>
        <button @click="hide">关闭</button>
      </div>
    </div>
    <!--搜索群成员-->
    <div v-show="componentStatus === 'searchMembers'" style="width: 100%;">
      <div class="search-bar" >
        <input
          type="text"
          v-model="query"
          placeholder="搜索群成员..."
          @compositionstart="isComposing = true"
          @compositionend="isComposing = false;triggerSearch()"
          @input="triggerSearch"
          ref="searchBar"
        />
      </div>
      <!--群成员列表-->
      <div class="search-members">
        <div 
          v-if="!filteredMembers && filteredMembers.length !== 0"
          v-for="member in filteredMembers" 
          :key="member.account_id" 
          class="member"
          @click="showProfileCard($event, member.account_id)" 
          @contextmenu="showContextMenu($event, member.account_id)"
        >
          <img :src="member.avatar" alt="avatar" class="avatar">
          <p class="remark">{{ member.remark? member.remark : (member.group_nickname?member.group_nickname:member.nickname)}}</p>
        </div>
        <div v-else class="no-result">
          <p>无搜索结果</p>
        </div>
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
    <!--管理员设置-->
    <div v-show="componentStatus === 'manage'" style="width: 100%;">
      <p>管理员设置</p>
      <p class="title">全体禁言: <SwitchButton v-model="groupInfo.muteAll" @change-value="setAllBanned"/></p>
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
      <p class="title">更改群头像
        <button @click="this.$refs.fileInput.click();">上传</button>
        <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" accept="image/*" />
      </p>
    </div>
    <PersonProfileCard ref="profileCard" />
    <ContextMenu ref="contextMenu" @select-item="handleMenuSelect"/>
    <InviteMember v-show="inviteMemberVisible" @close="inviteMemberVisible=false"/>
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
import PersonProfileCard from '@/components/base/PersonProfileCard.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import InviteMember from '@/components/Chat_list/InviteMember.vue';
export default {
  components: {
    EditableText,
    SwitchButton,
    SearchBar,
    PersonProfileCard,
    ContextMenu,
    InviteMember,
  },
  components: {
    EditableText,
    SwitchButton,
    SearchBar,
    PersonProfileCard,
    ContextMenu,
    InviteMember,
  },
  data() {
    return {
      visible: false,

      isDivideMoveVisible: false,
      divides:[],

      query: "", // 搜索关键词
      isComposing: false, // 是否正在使用输入法输入，防止频繁触发搜索
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
        group_avatar: '',
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
      // 搜索历史方面
      searchHistoryKeyword:'',
      searchHistoryType:'all', // 'all', 'image', 'file', 'date','member'
      searchHistoryDate:'',
      searchHistoryMember:'',
      

      searchMembersKeyword:'',
      componentStatus: 'main',  // 'main', 'history', 'manage'
      inviteMemberVisible: false,
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
          this.groupInfo = response.data.data;
        }else{
          this.$root.notify(response.data.message, 'error');
        }
        const response2 = await contactListAPI.getDivides('groups');
        if(response2.status === 200){
          this.divides = response2.data.divides;
        }else{
          this.$root.notify(response2.data.message, 'error');
        }
      }
      catch(error){
        console.log('fetch group error:', error);
      }
      
    },
    initialize(){
      this.query = '';
      this.showAll = false;
      this.searchHistoryKeyword = '';
      this.searchHistoryType = 'all';
      this.searchHistoryDate = '';
      this.searchHistoryMember = '';
      this.searchMembersKeyword = '';
    },
    returnTo(){
      if(this.componentStatus === 'main'){
        this.hide();
      }else{
        this.initialize();
        this.componentStatus = 'main';
      }
    },
    // 搜索框
    triggerSearch() {
      if (this.isComposing) return; // 正在输入中，不触发搜索
      if(this.query === ''){
        this.componentStatus = 'main';
      }else if(this.componentStatus === 'main'){
        this.componentStatus = 'searchMembers';
        this.$nextTick(() => {
          this.$refs.searchBar.focus();
        });
      }
      this.searchMembersKeyword = this.query;
    },
    // searchMember(key){
    //   this.searchMembersKeyword = key;
    // },
    async showProfileCard(event, account_id){
      try{
        const response = await getPersonProfileCard(account_id, this.group_id);
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
    inviteMember(){
      this.inviteMemberVisible = true;

    },
    
    async changeGroupNickname(newNickname){
      try {
        const response = await contactListAPI.changeGroupNickname(this.group_id, newNickname);
        if (response.status === 200) {
          this.groupInfo.my_group_nickname = newNickname;
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      } catch (error) {
        console.log('change group nickname error:', error);
      }
    },
    // 对群的个人设置
    async changeGroupRemark(newRemark){
      try{
        const response = await contactListAPI.changeRemark(this.group_id, true,newRemark);
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
    // 更改分组
    showDivideMove(){
      this.isDivideMoveVisible = true;
      this.$refs.divideMove.selectedDivide = this.groupInfo.divide;
      this.$refs.divideMove.multiple = false;
    },
    async divideMove(divide){
      try{
        const response = await contactListAPI.moveInDivide('groups',this.account_id, divide);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        else{
          this.groupInfo.divide = divide;
        }
      }catch(e){
        console.log(e);
      }
    },
    async setMute(){
      try{
        const response = await chatListAPI.setMute(this.group_id, !this.isMute, true);
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
        const response = await chatListAPI.blockChat(this.group_id, !this.isBlocked, true);
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
        const response = await chatListAPI.pinChat(this.group_id, !this.isPinned, true);
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
        const response = await contactListAPI.setBanned(this.group_id,this.account_id, this.is_banned);
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
    async exitGroup(){
      // 退出群聊
      try{
        const response = await contactListAPI.exitGroup(this.group_id);
        if(response.status === 200){
          this.$emit('group-exited', this.group_id);
          this.hide();
        }
        else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('exit group error:', error);
      }
    },
    // // async addMember() {
    // //   try {
    // //     await addMemberToGroup(this.groupInfo.id, this.newMemberId);
    // //     this.groupInfo.members.push({ id: this.newMemberId, name: '新成员' }); // 假设新成员的名字为 '新成员'
    // //     this.newMemberId = '';
    // //   } catch (error) {
    // //     console.error('Failed to add member:', error);
    // //   }
    // // },
    // // async removeMember(memberId) {
    // //   try {
    // //     await removeMemberFromGroup(this.groupInfo.id, memberId);
    // //     this.groupInfo.members = this.groupInfo.members.filter(member => member.id !== memberId);
    // //   } catch (error) {
    // //     console.error('Failed to remove member:', error);
    // //   }
    // // },
    // async deleteGroup() {
    //   try {
    //     await deleteGroup(this.groupInfo.id);
    //     this.$emit('group-deleted', this.groupInfo.id);
    //     this.hide();
    //   } catch (error) {
    //     console.error('Failed to delete group:', error);
    //   }
    // },

    // 聊天记录
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

    manageGroups() {
      // 管理员设置
      this.componentStatus = 'manage';
    },
    // 处理文件选择
    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.group_avatar = e.target.result;
        };
        reader.readAsDataURL(file);
        this.changeGroupAvatar(this.group_avatar);
      }
    },
    async changeGroupAvatar(avatar){
      // 更改群头像
      try{
        const response = await contactListAPI.changeGroupAvatar(this.group_id, avatar);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('change group avatar error:', error);
      }
    },
    async setAllBanned(){
      try{
        const response = await contactListAPI.setAllBanned(this.group_id, this.groupInfo.muteAll);
        if(response.status === 200){
          this.groupInfo.muteAll = !this.groupInfo.muteAll;
        }
        else{
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('set all banned error:', error);
      }
    },
    async changeInvitePermission(){
      try{
        const response = await contactListAPI.setAllowInvite(this.group_id, this.groupInfo.allow_invite);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('change invite permission error:', error);
      }
    },
    async changeIdPermission(){
      try{
        const response = await contactListAPI.setAllowIdSearch(this.group_id, this.groupInfo.allow_id_search);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch(error){
        console.log('change id permission error:', error);
      }
    },
    async changeNamePermission(){
      try{
        const response = await contactListAPI.setAllowNameSearch(this.group_id, this.groupInfo.allow_name_search);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
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
          const response = await contactListAPI.setBanned(this.group_id, account_id, true);
          if(response.status === 200){
            let member = this.groupInfo.members.find(member => member.account_id === account_id);
            if(member){
              member.is_banned = true;
            }
          }
          else{
            this.$root.notify(response.data.message, 'error');
          }
        }
        catch(error){
          console.log('set banned error:', error);
        }
      }
      else if(option==='解禁'){
        try{
          const response = await contactListAPI.setBanned(this.group_id, account_id, false);
          if(response.status === 200){
            let member = this.groupInfo.members.find(member => member.account_id === account_id);
            if(member){
              member.is_banned = false;
            }
          }
          else{
            this.$root.notify(response.data.message, 'error');
          }
        }
        catch(error){
          console.log('set banned error:', error);
        }
      }
      else if(option==='移除'){
        try{
          const response = await contactListAPI.removeMember(this.group_id, account_id);
          if(response.status === 200){
            this.groupInfo.members = this.groupInfo.members.filter(member => member.account_id !== account_id);
          }
          else{
            this.$root.notify(response.data.message, 'error');
          }
        }
        catch(error){
          console.log('remove member error:', error);
        }
      }else if(option==='设为管理员'){
        try{
          const response = await contactListAPI.setAdmin(this.group_id, account_id, true);
          if(response.status === 200){
            let member = this.groupInfo.members.find(member => member.account_id === account_id);
            if(member){
              member.group_role = 'group_manager';
            }
          }
          else{
            this.$root.notify(response.data.message, 'error');
          }
        }
        catch(error){
          console.log('set manager error:', error);
        }
      }else if(option==='取消管理员'){
        try{
          const response = await contactListAPI.setAdmin(this.group_id, account_id, false);
          if(response.status === 200){
            let member = this.groupInfo.members.find(member => member.account_id === account_id);
            if(member){
              member.group_role = 'group_member';
            }
          }
          else{
            this.$root.notify(response.data.message, 'error');
          }
        }
        catch(error){
          console.log('set manager error:', error);
        }
      }else if(option==='转让群主'){
        try{
          const response = await contactListAPI.transferOwner(this.group_id, account_id);
          if(response.status === 200){
            this.groupInfo.group_owner = account_id;
          }
          else{
            this.$root.notify(response.data.message, 'error');
          }
        }
        catch(error){
          console.log('transfer owner error:', error);
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
    filteredMembers(){
      const keyword = this.searchMembersKeyword;
      if(!keyword) return this.groupInfo.members;
      return this.groupInfo.members.filter(member => {
        return member.group_nickname.includes(keyword) || member.id.includes(keyword) 
            || member.remark.includes(keyword) || member.nickname.includes(keyword);
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