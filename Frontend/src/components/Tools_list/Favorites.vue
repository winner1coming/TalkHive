<template>
  <div class="favorites">
    <h2>收藏</h2>
    <ul>
      <li v-for="favorite in favorites" :key="favorite.id">
        {{ favorite.title }}
      </li>
    </ul>
    <button @click="addFavorite">添加收藏</button>
  </div>
</template>

<script>
import { getFavorites, addFavorite } from '../services/api';

export default {
  name: 'Favorites',
  data() {
    return {
      favorites: [],
    };
  },
  methods: {
    async fetchFavorites() {
      const response = await getFavorites();
      this.favorites = response.data;
    },
    async addFavorite() {
      const itemId = prompt('请输入要收藏的项的 ID');
      if (itemId) {
        await addFavorite(itemId);
        this.fetchFavorites();
      }
    },
  },
  created() {
    this.fetchFavorites();
  },
};
</script>

<style scoped>
.favorites {
  padding: 20px;
}
</style>