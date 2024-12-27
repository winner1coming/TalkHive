<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <h2>加好友/群</h2>
      <SearchBar :isImmidiate="false" @search="search" @button-click="search"/>
      <ul class="items">
      <!-- 每个消息项 -->
      <li 
        v-for="result in results" 
        :key="result.tid"
      >
        <div class="avatar">   <!-- 头像-->
          <img :src="result.avatar" alt="avatar" />
        </div>
        <div class="info">   <!-- 信息-->
          <div class="name">{{ result.nickname }}</div>
          <div class="remark">{{ result.id }}</div>
        </div>
        <div >   
          <button @click="add(result.tid)">添加</button>
        </div>
      </li>
    </ul>
    </div>
  </div>
</template>

<script>
import { addStranger, searchStrangers } from '@/services/api';
import SearchBar from '@/components/base/SearchBar.vue';
export default {
  components: {
    SearchBar,
  },
  data() {
    return {
      // results:[
      //   {
      //     tid: '13872132',   // 若为群聊，则为群号
      //     id: '13872132',
      //     nickname: 'test',
      //     avatar: '',
      //   },
      // ],  // 搜索结果
      results: [],
    };
  },
  methods: {
    async search(query) {
      if(!query) return;
      try{
        const response = await searchStrangers(tid);
        if (response.status!==200) {
          console.error('Failed to add friend/group', response.data.message);
        }else{
          this.results = response.data;
        }
      }
      catch (error){
        console.error('Failed to search friend/group',error)
      }
      
    },
    async add(tid) {
      try{
        const response = await addStranger(tid);
        if (response.status!==200) {
          console.error('Failed to add friend/group', response.data.message);
        }
      }
      catch (error){
        console.error('Failed to add friend/group',error)
      }
    },
    close() {
      this.$emit('close');
    },
  },
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
  background-color: #fff;
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
  font-size: 1.2rem;
}
.remark {
  font-size: 0.8rem;
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>