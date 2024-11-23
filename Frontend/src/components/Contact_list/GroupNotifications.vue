<template>
  <div class="group-notification">
    <h2>群聊申请</h2>
    <ul>
      <li v-for="request in requests" :key="request.id">
        {{ request.name }}
        <button @click="acceptRequest(request.id)">接受</button>
        <button @click="rejectRequest(request.id)">拒绝</button>
      </li>
    </ul>
  </div>
</template>

<script>
import { getGroupRequests, acceptGroupRequest, rejectGroupRequest } from '@/services/api.js';

export default {
  name: 'GroupRequest',
  data() {
    return {
      requests: [],
    };
  },
  methods: {
    async fetchRequests() {
      const response = await getGroupRequests();
      this.requests = response.data;
    },
    async acceptRequest(requestId) {
      await acceptGroupRequest(requestId);
      this.fetchRequests();
    },
    async rejectRequest(requestId) {
      await rejectGroupRequest(requestId);
      this.fetchRequests();
    },
  },
  created() {
    this.fetchRequests();
  },
};
</script>

<style scoped>
.group-notification {
  padding: 20px;
}
</style>