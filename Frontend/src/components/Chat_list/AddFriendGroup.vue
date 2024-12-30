<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <h2>加好友/群</h2>
      <SearchBar :isImmidiate="false" @search="search" @button-click="search"/>
      <ul class="items" v-show="results.length">
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
            <button @click="tryAdd(result.tid, result.type)">添加</button>
          </div>
        </li>
    </ul>
    </div>
  </div>
  <div class="add-modal" @click.self="close" v-show="isAddVisible">
    <div class="add-content">
      <h2>申请理由</h2>
      <textarea 
        v-model="reason" 
        placeholder="输入申请理由.."
      />
      <button @click="add(reason)">确认</button>
    </div>
  </div>
</template>

<script>
import * as contactListAPI from '@/services/contactList';
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
      //     type:''
      //   },
      // ],  // 搜索结果
      results: [],
      isAddVisible: false,
      tid: '',
      type: '',
    };
  },
  methods: {
    async search(query) {
      if(!query) return;
      try{
        const response = await contactListAPI.searchStrangers(query);
        if (response.status!==200) {
          this.$root.notify(response.data.message, 'error');
        }else{
          this.results = response.data.data;
        }
      }
      catch (error){
        console.error('Failed to search friend/group',error)
      }
      
    },
    tryAdd(tid, type) {
      this.tid = tid;
      this.type = type;
      this.isAddVisible = true;
    },
    async add(reason) {
      try{
        let response;
        if(this.type==='group'){
          response = await contactListAPI.addGroup(this.tid, reason, );
        }else{
          response = await contactListAPI.addFriend(this.tid, reason);
        }
        if (response.status!==200) {
          this.$root.notify(response.data.message, 'error');
        }
      }
      catch (error){
        console.error('Failed to add friend/group',error)
      }
      this.close();
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
  max-height: 300px;
  overflow-y: auto;
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

.add-modal {
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
.add-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
textarea {
  margin: 10px;
  width: 80%;
  height: 70%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: none;
}
</style>