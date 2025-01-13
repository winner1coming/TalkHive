<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <h2>邀请好友</h2>
      <SearchBar :isImmidiate="false" @search="search" @button-click="search"/>
      <ul class="items">
        <li 
          v-for="friend in filteredFriends" 
          :key="friend.accound_id"
        >
          <div class="avatar">   <!-- 头像-->
            <img :src="friend.avatar" alt="avatar" />
          </div>
          <div class="info">   <!-- 信息-->
            <div class="name">{{ friend.remark? friend.remark : friend.name }}</div>
            <div class="remark">{{ friend.id }}</div>
          </div>
          <div >   
            <button @click="inviteMember(friend.account_id)">邀请</button>
          </div>
        </li>
    </ul>
    </div>
  </div>
</template>

<script>
import SearchBar from '@/components/base/SearchBar.vue';
import * as contactListAPI from '@/services/contactList';
export default {
  props:['group_id'],
  components: {
    SearchBar,
  },
  data() {
    return {
      friends: [
        // {
        //   accound_id: 1,
        //   remark: '张三',  // 备注
        //   nickname: '张三',  // 昵称
        //   id: 'zhangsan',
        //   avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
        //   signature: '这个人很懒，什么都没留下',
        //   divide: '家人',
        // },
        // {
        //   accound_id: 2,
        //   remark: '李四',
        //   nickname: '李四',
        //   id: 'lisi',
        //   avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
        //   signature: '这个人很懒，什么都没留下',
        //   divide: '家人',
        // },
        // {
        //   accound_id: 3,
        //   remark: '王五',
        //   nickname: '王五',
        //   id: 'wangwu',
        //   avatar: 'https://cdn.jsdelivr.net/gh/lin09/dist/img/avatar.jpg',
        //   signature: '这个人很懒，什么都没留下',
        //   divide: '家人',
        // },
      ],
      filteredFriends: [],
    };
  },
  methods: {
    search(query) {
      if (!query) {
        this.filteredFriends = this.friends;
        return;
      }
      // 根据搜索条件过滤好友列表
      this.filteredFriends = this.friends.filter(friend => {
        if(!friend.remark){
          return friend.nickname.includes(query) || friend.id.includes(query) || friend.remark.includes(query);
        }else{
          return friend.nickname.includes(query) || friend.id.includes(query);
        }
      });
      
    },
    async fetchFriendsNotInGroup() {
      try {
        const response = await contactListAPI.fetchFriendsNotInGroup(this.group_id);
        if (response.status !== 200) {
          this.$root.notify(response.data.message, 'error');
        } else {
          this.friends = response.data.data;
          this.filteredFriends = this.friends;
        }
      } catch (error) {
        console.error('Failed to fetch friends not in group', error);
      }
    },
    async inviteMember(friend_id) {
      this.close();
      try {
        const response = await contactListAPI.inviteMember(this.group_id, friend_id);
        if (response.status !== 200) {
          this.$root.notify(response.data.message, 'error');
        } 
      } catch (error) {
        console.error('Failed to invite member', error);
      }
    },
    close() {
      this.$emit('close');
    },

  },
  created(){
    this.fetchFriendsNotInGroup();
  }
};
</script>

<style scoped src="@/assets/css/chatList.css"></style>
<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000; /* 确保在最上层 */
}

.modal-content {
  background-color: var(--background-color);
  color: var(--text-color);
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  height: 400px;
}
.items {
  list-style: none;
  padding: 0;
}
.items li {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}
.avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.info {
  flex: 5;
  margin-left: 10px;
  text-align: left;
}
.name{
  font-weight: bold;
  font-size: var(--font-size-mlarge);
}
.remark {
  font-size: var(--font-size-small);
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>