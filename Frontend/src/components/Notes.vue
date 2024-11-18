<template>
  <div class="notes">
    <h2>笔记</h2>
    <ul>
      <li v-for="note in notes" :key="note.id">
        {{ note.title }}
      </li>
    </ul>
    <button @click="createNote">创建笔记</button>
  </div>
</template>

<script>
import { getNotes, createNote } from '../services/api';

export default {
  name: 'Notes',
  data() {
    return {
      notes: [],
    };
  },
  methods: {
    async fetchNotes() {
      const response = await getNotes();
      this.notes = response.data;
    },
    async createNote() {
      const title = prompt('请输入笔记标题');
      const content = prompt('请输入笔记内容');
      if (title && content) {
        await createNote(title, content);
        this.fetchNotes();
      }
    },
  },
  created() {
    this.fetchNotes();
  },
};
</script>

<style scoped>
.notes {
  padding: 20px;
}
</style>