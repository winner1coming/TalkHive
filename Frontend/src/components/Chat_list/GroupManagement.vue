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
        <div class="group-members">
          <div v-for="member in groupInfo.members" :key="member.account_id" class="member">
            <img :src="member.avatar" alt="avatar" class="avatar">
            <p class="remark">{{ member.group_nickname.length > this.maxChars ? member.group_nickname.slice(0, this.maxChars)+'...' : member.group_nickname }}</p>
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
        <p>是否全体禁言: <SwitchButton v-model="groupInfo.muteAll" @change-value=""/></p>
        <p>是否可以通过群成员邀请进入: <SwitchButton v-model="groupInfo.allowMemberInvite" @change-value=""/></p>
        <p>是否可以通过群号搜索进入: <SwitchButton v-model="groupInfo.allowSearch" @change-value=""/></p>
      </div>
      <div class="group-actions">
        <div>
          <input v-model="newMemberId" placeholder="输入成员ID">
          <button @click="addMember">添加成员</button>
        </div>
        <button @click="deleteGroup">删除群聊</button>
        <button @click="hide">关闭</button>
      </div>
    </div>
    <div v-show="componentStatus === 'history'">
      <p>聊天记录</p>
    </div>
    <div v-show="componentStatus === 'manage'">
      <p>管理员设置</p>
    </div>
  </div>
</template>

<script>
import * as contactListAPI from '@/services/contactList';
import * as chatListAPI from '@/services/chatList';
import { EventBus } from '@/components/base/EventBus';
import EditableText from '@/components/base/EditableText.vue';
import SwitchButton from '@/components/base/SwitchButton.vue';
import { changeGroupNickname } from '../../services/contactList';
import SearchBar from '@/components/base/SearchBar.vue';
export default {
  components: {
    EditableText,
    SwitchButton,
    SearchBar,
  },
  data() {
    return {
      visible: false,
      group_id:'',
      group_remark:'',
      isMute: false,
      isBlocked: false,
      isPinned: false,
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
        // 入群权限 todo
        
        my_group_nickname: '',   // 我在本群的群昵称
        members: [],
        my_group_role:'',
      },
      newMemberId: '',
      componentStatus: 'main',  // 'main', 'history', 'manage'
    };
  },
  watch: {
    "$store.state.currentChat": {
      immediate: true,
      deep: true,
      handler(newVal) {
        if(!newVal) return;
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
    viewChatHistory() {
      // 查看聊天记录
      this.componentStatus = 'history';
    },
    manageGroups() {
      // 管理员设置
      this.componentStatus = 'manage';
    },
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
  },
  created(){
    //this.fetchGroupInfo();
  },
  mounted() {
    
    EventBus.on('other-float-component', (component) => {
      if (this.visible && component !== this) {
        this.hide();
      }
    });
    EventBus.on('close-float-component', (clickedElement) => {
      console.log(clickedElement);
      if (this.visible && !this.$el.contains(clickedElement)) {
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
  padding: 20px;
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
</style>