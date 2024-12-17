<template>
  <div v-if="visible" class="group-management">
    <h2>群聊管理</h2>
    <div class="group-info">
      <p><strong>群聊名称:</strong> {{ group.name }}</p>
      <p><strong>群聊ID:</strong> {{ group.id }}</p>
      <p><strong>成员列表:</strong></p>
      <ul>
        <li v-for="member in group.members" :key="member.id">
          {{ member.name }} ({{ member.id }})
          <button @click="removeMember(member.id)">移除</button>
        </li>
      </ul>
    </div>
    <div class="group-actions">
      <input v-model="newMemberId" placeholder="输入成员ID">
      <button @click="addMember">添加成员</button>
      <button @click="deleteGroup">删除群聊</button>
      <button @click="hide">关闭</button>
    </div>
  </div>
</template>

<script>
import { getGroups, createGroup, deleteGroup } from '@/services/contactList';
import { EventBus } from '@/components/base/EventBus';

export default {
  name: 'GroupManagement',
  data() {
    return {
      group:{
        id: 0,
        name: '',
        members: [],
      },
      visible: false,
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

<style scoped>
.group-management {
  padding: 20px;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 5px;
}
.group-info {
  margin-bottom: 20px;
}
.group-actions input {
  margin-right: 10px;
}
</style>