<template>
  <div class="main">
    <div class="contact-header">
        好友申请
    </div>
    <div v-for="request in requests" :key="request.apply_id" class="item">
      <img :src="request.avatar" alt="avatar" width="50" height="50" />
      <div class="left">
        <p class="name">{{ request.name }}</p>
        <p class="remark">备注: {{ request.reason }}</p>
      </div>
      <div class="right">
        <p class="remark">{{ request.time }}</p>
        <div v-if="request.status === 'pending' && request.receiver_id === this.$store.state.user.id">
          <button @click="acceptRequest(request.sender_id)">同意</button>
          <button @click="rejectRequest(request.sender_id)">拒绝</button>
        </div>
        <p v-else-if="request.status === 'pending'">等待对方处理</p>
        <p v-else-if="request.status === 'accepted'">已同意</p>
        <p v-else-if="request.status === 'rejected'">已拒绝</p>
      </div>
      
    </div>
  </div>
</template>

<script>
import { getFriendRequests, friendRequestPend } from '@/services/contactList';
const contactListAPI = {getFriendRequests, friendRequestPend};

export default {
  name: 'FriendRequest',
  data() {
    return {
      // requests: [
      //   {
      //     apply_id: '1',
      //     avatar: '',
      //     name: 'John Doe',
      //     sender_id: '1', // 申请者的tid
      //     receiver_id: '2',   // 接收者的tid
      //     reason:"I want to be your friend",
      //     status: 'pending',   // pending, accepted, rejected
      //     time: '2021-01-01 12:00:00',  // 待定
      //   },
      //   {
      //     apply_id: '2',
      //     avatar: '',
      //     name: 'Jane Doe',
      //     sender_id: '2', // 申请者的tid
      //     receiver_id: '1',   // 接收者的tid
      //     reason:"I want to be your friend",
      //     status: 'accepted',
      //     time: '2021-01-01 12:00:00',
      //   },
      // ],
      requests: [],
    };
  },
  methods: {
    async fetchRequests() {
      const response = await contactListAPI.getFriendRequests();
      this.requests = response.data;
    },
    async acceptRequest(requestId) {
      await contactListAPI.friendRequestPend(requestId, true);
      this.fetchRequests();
    },
    async rejectRequest(requestId) {
      await contactListAPI.friendRequestPend(requestId, false);
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