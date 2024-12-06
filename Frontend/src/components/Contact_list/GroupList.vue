<template>
  <div class="main">
    <div class="contact-header">
        群聊列表
    </div>
    <itemList :items="items" :type="type" :tags="tags" />
  </div>
</template>

<script>
import { getGroups } from '@/services/contactList';
import itemList from './itemList.vue';

export default {
  components: {
    itemList,
  },
  data() {
    return {
      type: 'groupList',  // friendList, groupList
      tags: ['家人', '朋友', '同事'],  // 从后端获取
      items: [   // 从后端获取
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
    async fetchGroups() {
      const response = await getGroups();
      this.items = response.data;
    },
  },
  created() {
    this.fetchGroups();
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>