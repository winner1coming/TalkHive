<template>
  <div class="favorites">
    <h2>我的收藏</h2>
    <ul>
      <li v-for="item in favorites" :key="item.message_id" class="favorite-item">
        <span class="type">{{ getTypeLabel(item.type) }}</span>
        <span class="object-name">{{ item.object_name }}</span>
        <span class="sender"> - {{ item.sender_name }}</span>
        <span class="time"> - {{ item.time }}</span>
        <button @click="viewItem(item)">查看</button>
      </li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Favorites',
  data() {
    return {
      favorites: [],
      userId: 'your-user-id', // 用户的ID
    };
  },
  created() {
    this.fetchFavorites();
  },
  methods: {
    // 从后端获取收藏列表
    async fetchFavorites() {
      try {
        const response = await axios.get(`/workspace/favorites`, {
          params: { id: this.userId }, // 向后端传递用户ID
        });
        if (response.data.status === 200) {
          this.favorites = response.data.data;
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error("无法获取收藏列表:", error);
        alert("获取收藏列表失败！");
      }
    },
    // 获取收藏项的类型标签
    getTypeLabel(type) {
      switch (type) {
        case 'message': return '消息';
        case 'code': return '代码';
        case 'note': return '笔记';
        default: return '未知';
      }
    },
    // 点击查看某个收藏项
    viewItem(item) {
      if (item.type === 'message') {
        this.$router.push({
          path: `/workspace/favorites/${item.message_id}`,
          query: { table: item.message_list_name },
        });
      } else if (item.type === 'code' || item.type === 'note') {
        this.$router.push({
          path: `/workspace/${item.type}s/${item.message_id}`,
        });
      }
    },
  },
};
</script>

<style scoped>
.favorites {
  padding: 20px;
}

.favorites h2 {
  margin-bottom: 20px;
}

.favorites ul {
  list-style: none;
  padding: 0;
}

.favorite-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #ddd;
}

.favorite-item .type,
.favorite-item .object-name {
  font-weight: bold;
  margin-right: 10px;
}

.favorite-item .sender {
  color: #666;
}

.favorite-item button {
  margin-left: auto;
  padding: 5px 10px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

.favorite-item button:hover {
  background-color: #0056b3;
}
</style>
