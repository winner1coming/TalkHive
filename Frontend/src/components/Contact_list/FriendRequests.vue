<template>
  <div class="main">
    <div class="contact-header">
        好友申请
    </div>
    <div v-for="request in requests" :key="request.account_id" class="item">
      <img :src="request.avatar" alt="avatar" width="50" height="50" />
      <div class="left">
        <p class="name">{{ request.name }}</p>
        <p class="remark">备注: {{ request.reason }}</p>
      </div>
      <div class="right">
        <p class="remark">{{ request.time }}</p>
        <div v-if="request.status === 'pending'">
          <button @click="acceptRequest(request.account_id)">同意</button>
          <button @click="rejectRequest(request.account_id)">拒绝</button>
        </div>
        <p v-else-if="request.status === 'accepted'">已同意</p>
        <p v-else-if="request.status === 'rejected'">已拒绝</p>
        <p v-else-if="request.status === 'waiting'">等待对方处理</p>
      </div>
      
    </div>
  </div>
</template>

<script>
import { getFriendRequests, acceptFriendRequest, rejectFriendRequest } from '@/services/api.js';

export default {
  name: 'FriendRequest',
  data() {
    return {
      requests: [
        {
          avatar: '',
          name: 'John Doe',
          account_id: '1',   // 申请者的id
          reason:"I want to be your friend",
          status: 'pending',   // pending（需要我处理的）, accepted, rejected, waiting（对方处理中）
          time: '2021-01-01 12:00:00',  // 待定
        },
        {
          avatar: '',
          name: 'Jane Doe',
          account_id: '2',
          reason:"I want to be your friend",
          status: 'accepted',
          time: '2021-01-01 12:00:00',
        },
      ],
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

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>
button {
  margin-right: 5px;
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
button:first-of-type {
  background-color: #28a745;
  color: white;
}
button:last-of-type {
  background-color: #dc3545;
  color: white;
}
</style>