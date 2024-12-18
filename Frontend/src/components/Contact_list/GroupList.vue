<template>
  <div class="main">
    <div class="contact-header">
        群聊列表
    </div>
    <itemList :items="items" :type="type" :tags="tags" @show-profile-card="showProfileCard"/>
    <ProfileCard ref="profileCard"/>

  </div>
</template>

<script>
import { getGroups } from '@/services/contactList';
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
      type: 'groupList',  // friendList, groupList
      tags: ['家人', '朋友', '同事'],  // 从后端获取
      // items: [   // 从后端获取
      //   {
      //     avatar: '',
      //     account_id: '1',   // 群id
      //     signature: '这是一个群聊',  // 群介绍
      //     remark: 'John',   // 群聊备注或群名称
      //     tag: '家人',
      //   },
      // ],
      items: [],
      boundD: 0,
      boundR: 0,
    };
  },
  methods: {
    async fetchGroups() {
      const response = await getGroups();
      this.items = response.data;
    },
    async showProfileCard(event, send_account_id){
      const response = await getProfileCard(send_account_id); 
      const profile = response.data;
      this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
    },
  },
  mounted() {
    this.fetchGroups();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;

  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>