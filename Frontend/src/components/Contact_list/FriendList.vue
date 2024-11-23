<template>
  <div class="friend-list">
    <h2>好友列表</h2>
    <ul>
      <li v-for="friend in friends" :key="friend.id">
        {{ friend.name }}
      </li>
    </ul>
  </div>
</template>

<script>
import { getFriends } from '@/services/api.js';

export default {
  name: 'FriendList',
  data() {
    return {
      friends: [],
    };
  },
  methods: {
    async fetchFriends() {
      const response = await getFriends();
      this.friends = response.data;
    },
  },
  created() {
    this.fetchFriends();
  },
};
</script>

<style scoped>
.friend-list {
  padding: 20px;
}
</style>