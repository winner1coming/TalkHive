<template>
  <div v-if="visible" class="group-management">
    <h2>群聊管理</h2>
    <div class="group-members">
      <div v-for="member in groupInfo.members" :key="member.account_id" class="member">
        <img :src="member.avatar" alt="avatar" class="avatar">
        <p>{{ member.group_nickname }}</p>
      </div>
    </div>
    <div class="group-info">
      <p>群聊名称:</p>
      <EditableText v-if="chat" :text="groupInfo.group_name" @update-text="groupInfo.group_name = $event" />  
      <p>群介绍:</p>
      <EditableText :text="groupInfo.introduction" @update-text="groupInfo.introduction = $event" />
      <p>群聊备注: </p>
      <EditableText :text="group_remark" @update-text="groupInfo.remark = $event" />
      <p>我的群昵称: </p>
      <EditableText :text="groupInfo.my_group_nickname" @update-text="groupInfo.my_group_nickname = $event" />
      <hr class="divider" />
      <!-- <p>是否显示群成员昵称: <SwitchButton v-model="groupInfo.showNicknames" /></p> -->
      <p>是否消息免打扰: <SwitchButton v-model="isMute" @change-value="setMute"/></p>
      <p>是否屏蔽: <SwitchButton v-model="isBlocked" @change-value="setBlock"/></p>
      <p>是否置顶: <SwitchButton v-model="isPinned" @change-value="setPin"/></p>
      <hr class="divider" />
      <p>聊天记录: <button @click="viewChatHistory">查看</button></p>
      <hr class="divider" />
      <p>管理员设置: <button @click="manageAdmins">设置</button></p>
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
</template>

<script>
import { getGroups, createGroup, deleteGroup } from '@/services/contactList';
import * as chatListAPI from '@/services/chatList';
import { EventBus } from '@/components/base/EventBus';
import EditableText from '@/components/base/EditableText.vue';
import SwitchButton from '@/components/base/SwitchButton.vue';
export default {
  components: {
    EditableText,
    SwitchButton
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
      groupInfo: {
        group_owner: '111',  // 群主tid
        introduction: '111',
        // 入群权限 todo
        
        my_group_nickname: 'aa',   // 我在本群的群昵称
        members: [
          {
            account_id: '111',
            avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
            group_role: 'group_owner',
            group_nickname: 'aa',
          },
          {
            account_id: '222',
            avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
            group_role: 'group_owner',
            group_nickname: 'bb',
          },
          {
            account_id: '333',
            avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
            group_role: 'group_owner',
            group_nickname: 'cc',
          },
        ],
        // showNicknames: false, 需求没做 todo
      },
      
      newMemberId: ''
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
    async setMute(){
      try{
        const response = await chatListAPI.setMute(this.group_id, !this.isMute);
        if(response.status === 200){
          this.isMute = !this.isMute;
          let chatInfo = { ...this.$store.state.currentChat };
          chatInfo.tags = this.isMute ? [...chatInfo.tags, 'mute'] : chatInfo.tags.filter(tag => tag !== 'mute');
          this.$store.dispatch('setChat',chatInfo);  todo
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
      // 查看聊天记录逻辑
    },
    manageAdmins() {
      // 管理员设置逻辑
    },
    show(){
      this.visible = true;
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide(){
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
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

<style scoped src="@/assets/css/chatList.css"></style>
<style scoped>
.group-management {
  padding: 20px;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 5px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.group-members {
  display: flex;
  flex-wrap: wrap;
}

.member {
  margin: 10px;
  text-align: center;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}

.group-info {
  margin-top: 20px;
  align-self: flex-start;
}
.group-info p {
  text-align: left;
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