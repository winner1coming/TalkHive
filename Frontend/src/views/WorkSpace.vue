<template>
  <div class="workspace">
    <!-- 左侧工具栏 -->
    <aside class="toolbar">
      <ul>
        <li
          @click="navigateTo('/workspace/collabdocs')"
          :class="{ active: currentRoute === '/workspace/collabdocs' }"
        >
          <img :src="iconEdit" alt="在线协作" class="toolbar-icon" />
          在线协作
        </li>

        <li
          @click="navigateTo('/workspace/favorites')"
          :class="{ active: currentRoute === '/workspace/favorites' }"
        >
          <img :src="iconFavorite" alt="收藏" class="toolbar-icon" />
          收藏
        </li>

        <li
          @click="navigateTo('/workspace/notes')"
          :class="{ active: currentRoute === '/workspace/notes' }"
        >
          <img :src="iconNotes" alt="笔记" class="toolbar-icon" />
          笔记
        </li>

        <li
          @click="navigateTo('/workspace/code')"
          :class="{ active: currentRoute === '/workspace/code' }"
        >
          <img :src="iconCode" alt="代码" class="toolbar-icon" />
          代码
        </li>

        <li
          @click="navigateTo('/workspace/ddl')"
          :class="{ active: currentRoute === '/workspace/ddl' }"
        >
          <img :src="iconDDL" alt="DDL" class="toolbar-icon" />
          DDL
        </li>

        <li
          @click="navigateTo('/workspace/recycle')"
          :class="{ active: currentRoute === '/workspace/recycle' }"
        >
          <img :src="iconTrash" alt="回收站" class="toolbar-icon" />
          回收站
        </li>
      </ul>
    </aside>

    <!-- 右侧内容区域 -->
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script>
import iconEdit from '@/assets/icon/edit.png?w=80&format=webp&imagetools';
import iconFavorite from '@/assets/icon/favorite.png?w=80&format=webp&imagetools';
import iconNotes from '@/assets/icon/notes.png?w=80&format=webp&imagetools';
import iconCode from '@/assets/icon/code.png?w=80&format=webp&imagetools';
import iconDDL from '@/assets/icon/date-icon.png?w=80&format=webp&imagetools';
import iconTrash from '@/assets/icon/trash.png?w=80&format=webp&imagetools';

export default {
  name: 'WorkSpace',
  data() {
    return {
      currentRoute: '',
      iconEdit,
      iconFavorite,
      iconNotes,
      iconCode,
      iconDDL,
      iconTrash,
    };
  },
  watch: {
    $route(to) {
      this.currentRoute = to.path;
    },
  },
  created() {
    this.currentRoute = this.$route.path;
  },
  methods: {
    navigateTo(route) {
      this.$router.push(route);
    },
  },
};
</script>

<style scoped>
.workspace {
  display: flex;
  height: 100vh;
}

.toolbar {
  width: 150px;
  background-color: var(--background-color1);
  border-right: 1px solid var(--background-color2);
  padding: 10px;
  box-sizing: border-box;
}

.toolbar ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.toolbar li {
  padding: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: background-color 0.3s;
  color: var(--text-color);
}

.toolbar li:hover {
  background-color: var(--select-background-color);
}

.toolbar li.active {
  font-weight: bold;
  background-color: var(--select-background-color1);
  color: var(--select-text-color);
}

.toolbar-icon {
  width: 30px;
  height: 30px;
}

.content {
  flex: 1;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;
  background-color: var(--background-color);
}
</style>
