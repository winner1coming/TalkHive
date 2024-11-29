<template>
  <div class="main">
    <div class="contact-header">
        好友列表
    </div>
    <itemList :items="items" :type="type" :tags="tags" />
  </div>
</template>

<script>
import { getFriends } from '@/services/contactList';
import itemList from './itemList.vue';

export default {
  components: {
    itemList,
  },
  data() {
    return {
      type: 'friendList',  // friendList, groupList
        tags: ['家人', '朋友', '同事'],
        items: [
          {
            avatar: '',
            account_id: '1',
            remark: 'John',   // 好友备注
            status: 'online',   // online, offline
            signature: '爱拼才会赢',
            tag: '家人',
          },
        ],
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

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>