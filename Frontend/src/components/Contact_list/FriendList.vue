<template>
  <div class="main">
    <div class="contact-header">
        好友列表
    </div>
    <itemList :items="items" :type="type" :tags="tags" @show-profile-card="showProfileCard"/>
    <ProfileCard ref="profileCard" @go-to-chat="goToChat" />
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
        items: [
          {
            avatar: '',
            account_id: '1',
            remark: 'John',   // 好友备注
            status: 'online',   // online, offline
            signature: '爱拼才会赢',    // 签名
            tag: '家人',   
          },
        ],
    };
  },
  methods: {
    async fetchFriends() {
      const response = await getFriends();
      this.items = response.data;
    },
    async showProfileCard(event, send_account_id){
      // const profile = await this.getProfileCard(send_account_id); todo
      const profile = {
        tid: '0',  // tid
        avatar: new URL('@/assets/images/avatar.jpg', import.meta.url).href,  // 头像地址
        remark: '',  // 备注
        nickname: '', // 对方设置的昵称
        groupNickname: '',  // 对方的群昵称
        tag: '',  // 分组
        signature: '',  // 个性签名
      };
      this.$refs.profileCard.show(event, profile);
    },
    async goToChat(tid){
      this.$emit('go-to-chat', tid);   // todo 切换组件
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