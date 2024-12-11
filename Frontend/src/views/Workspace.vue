<template>
    <div class="ddls">
      <h2>DDL</h2>
      <ul>
        <li v-for="ddl in ddls" :key="ddl.id">
          {{ ddl.title }} - {{ ddl.deadline }}
        </li>
      </ul>
      <button @click="createDDL">创建 DDL</button>
    </div>
  </template>
  
  <script>
  import { getDDLs, createDDL } from '../services/api';
  
  export default {
    name: 'DDLs',
    data() {
      return {
        ddls: [],
      };
    },
    methods: {
      async fetchDDLs() {
        const response = await getDDLs();
        this.ddls = response.data;
      },
      async createDDL() {
        const title = prompt('请输入 DDL 标题');
        const deadline = prompt('请输入 DDL 截止时间');
        if (title && deadline) {
          await createDDL(title, deadline);
          this.fetchDDLs();
        }
      },
    },
    created() {
      this.fetchDDLs();
    },
  };
  </script>
  
  <style scoped>
  .ddls {
    padding: 20px;
  }
  </style>