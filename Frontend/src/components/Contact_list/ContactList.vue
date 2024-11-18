<template>
  <div class="contact-list">
    <h2>通讯录</h2>
    <ul>
      <li v-for="contact in contacts" :key="contact.id">
        {{ contact.name }}
      </li>
    </ul>
  </div>
</template>

<script>
import { getFriends } from '../services/api';

export default {
  name: 'ContactList',
  data() {
    return {
      contacts: [],
    };
  },
  methods: {
    async fetchContacts() {
      const response = await getFriends();
      this.contacts = response.data;
    },
  },
  created() {
    this.fetchContacts();
  },
};
</script>

<style scoped>
.contact-list {
  padding: 20px;
}
</style>