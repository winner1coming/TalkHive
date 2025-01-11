<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <div>
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
              <div class="name">{{ result.remark }}</div>
            </div>
            <div >   
              <button @click="">+</button>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <!-- 已选中部分-->
    <div>
      
    </div>
  </div>
</template>

<script>
import { addFriendGroup, searchFriendGroup } from '@/services/api';
import SearchBar from '@/components/base/SearchBar.vue';
export default {
  components: {
    SearchBar,
  },
  data() {
    return {
      results:[ 
        {
          tid: '13872132',   // 若为群聊，则为群号
          id: '13872132',
          nickname: 'test',
          avatar: '',
        },
      ],  // 搜索结果
    };
  },
  methods: {
    async search(query) {
      this.results = await searchFriendGroup(query);
    },
    async add(tid) {
      await addFriendGroup(tid);
    },
    close() {
      this.$emit('close');
    },
  },
};
</script>

<style src="@/assets/css/contacList.css"></style>
<style scoped>
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