<template>
  <div class="main">
    <div class="contact-header">
        好友列表
        <button style="float: right;">分组管理</button>
    </div>
    <itemList :items="items" :type="type" :tags="tags" @show-profile-card="showProfileCard"/>
    <ProfileCard ref="profileCard" />
  </div>
</template>

<script>
import { getFriends } from '@/services/contactList';
import itemList from './itemList.vue';
import ProfileCard from '@/components/base/ProfileCard.vue';
import { getProfileCard } from '@/services/api';
export default {
  components: {
    itemList,
    ProfileCard,
  },
  data() {
    return {
      type: 'friendList',  // friendList, groupList
      tags: ['家人', '朋友', '同事'],
      // items: [
      //   {
      //     avatar: '',
      //     account_id: '1',
      //     remark: 'John',   // 好友备注
      //     status: 'online',   // online, offline
      //     signature: '爱拼才会赢',    // 签名
      //     tag: '家人',   
      //   },
      // ],
      items: [],
      boundD: 0,
      boundR: 0,
    };
  },
  methods: {
    async fetchFriends() {
      const response = await getFriends();
      this.items = response.data;
    },
    async showProfileCard(event, send_account_id){
      const response = await getProfileCard(send_account_id); 
      const profile = response.data;
      this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
    },
  },
  mounted() {
    this.fetchFriends();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>