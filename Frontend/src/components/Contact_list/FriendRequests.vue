<template>
  <div class="friend-request">
    <h2>好友申请</h2>
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
import { getFriendRequests, acceptFriendRequest, rejectFriendRequest } from '@/services/api.js';

export default {
  name: 'FriendRequest',
  data() {
    return {
      requests: [],
    };
  },
  methods: {
    async fetchRequests() {
      const response = await getFriendRequests();
      this.requests = response.data;
    },
    async acceptRequest(requestId) {
      await acceptFriendRequest(requestId);
      this.fetchRequests();
    },
    async rejectRequest(requestId) {
      await rejectFriendRequest(requestId);
      this.fetchRequests();
    },
  },
  created() {
    this.fetchRequests();
  },
};
</script>

<style scoped>
.friend-request {
  padding: 20px;
}
</style>