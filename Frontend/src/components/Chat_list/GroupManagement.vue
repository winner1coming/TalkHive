<template>
  <div v-if="visible" class="group-management">
    <h2>群聊管理</h2>
    <div class="group-members">
      <div v-for="member in group.members" :key="member.id" class="member">
        <img :src="member.avatar" alt="avatar" class="avatar">
        <p>{{ member.nickname }}</p>
      </div>
    </div>
    <div class="group-info">
      <p>群聊名称:</p>
      <EditableText :text="group.name" @update-text="group.name = $event" />  
      <p>群介绍:</p>
      <EditableText :text="group.description" @update-text="group.description = $event" />
      <p>群聊备注: </p>
      <EditableText :text="group.remark" @update-text="group.remark = $event" />
      <p>我的群昵称: </p>
      <EditableText :text="group.myNickname" @update-text="myNickname = $event" />
      <hr class="divider" />
      <p>是否显示群成员昵称: <SwitchButton v-model="group.showNicknames" /></p>
      <p>是否消息免打扰: <SwitchButton v-model="group.muteNotifications" /></p>
      <p>是否屏蔽: <SwitchButton v-model="group.blocked" /></p>
      <p>是否置顶: <SwitchButton v-model="group.pinned" /></p>
      <hr class="divider" />
      <p>聊天记录: <button @click="viewChatHistory">查看</button></p>
      <hr class="divider" />
      <p>管理员设置: <button @click="manageAdmins">设置</button></p>
      <p>是否全体禁言: <SwitchButton v-model="group.muteAll" /></p>
      <p>是否可以通过群成员邀请进入: <SwitchButton v-model="group.allowMemberInvite" /></p>
      <p>是否可以通过群号搜索进入: <SwitchButton v-model="group.allowSearch" /></p>
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
      group: {
        name: '111',
        id: '111',
        description: '111',
        remark: '111',
        myNickname: 'aa',
        members: [],
        muteAll: false,
        allowMemberInvite: false,
        allowSearch: false,
        showNicknames: false,
        muteNotifications: false,
        blocked: false,
        pinned: false,
      },
      
      newMemberId: ''
    };
  },
  methods: {
    show(){
      this.visible = true;
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide(){
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
    async addMember() {
      try {
        await addMemberToGroup(this.group.id, this.newMemberId);
        this.group.members.push({ id: this.newMemberId, name: '新成员' }); // 假设新成员的名字为 '新成员'
        this.newMemberId = '';
      } catch (error) {
        console.error('Failed to add member:', error);
      }
    },
    async removeMember(memberId) {
      try {
        await removeMemberFromGroup(this.group.id, memberId);
        this.group.members = this.group.members.filter(member => member.id !== memberId);
      } catch (error) {
        console.error('Failed to remove member:', error);
      }
    },
    async deleteGroup() {
      try {
        await deleteGroup(this.group.id);
        this.$emit('group-deleted', this.group.id);
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
    }
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