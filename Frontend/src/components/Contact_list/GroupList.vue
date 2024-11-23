<template>
  <div class="group-list">
    <h2>群聊列表</h2>
    <ul>
      <li v-for="group in groups" :key="group.id">
        {{ group.name }}
      </li>
    </ul>
  </div>
</template>

<script>
import { getGroups } from '@/services/api.js';

export default {
  name: 'GroupList',
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
  },
  created() {
    this.fetchGroups();
  },
};
</script>

<style scoped>
.group-list {
  padding: 20px;
}
</style>