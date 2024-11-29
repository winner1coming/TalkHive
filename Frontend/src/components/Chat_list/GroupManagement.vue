<template>
  <div class="group-management">
    <h2>群聊管理</h2>
    <ul>
      <li v-for="group in groups" :key="group.id">
        {{ group.name }}
        <button @click="editGroup(group.id)">编辑</button>
        <button @click="deleteGroup(group.id)">删除</button>
      </li>
    </ul>
    <button @click="createGroup">创建群聊</button>
  </div>
</template>

<script>
import { getGroups, createGroup, deleteGroup } from '@/services/contactList';

export default {
  name: 'GroupManagement',
  data() {
    return {
      groups: [],
    };
  },
  methods: {
    async fetchGroups() {
      const response = await getGroups();
      this.groups = response.data;
    },
    async createGroup() {
      const name = prompt('请输入群聊名称');
      if (name) {
        await createGroup(name);
        this.fetchGroups();
      }
    },
    async deleteGroup(groupId) {
      if (confirm('确定删除该群聊吗？')) {
        await deleteGroup(groupId);
        this.fetchGroups();
      }
    },
    editGroup(groupId) {
      // 编辑群聊逻辑
    },
  },
  created() {
    this.fetchGroups();
  },
};
</script>

<style scoped>
.group-management {
  padding: 20px;
}
</style>