<template>
  <div class="main">
    <div class="contact-header">
        黑名单
      <button style="float: right;" @click="showContextMenu">批量管理</button>
    </div>
    <div v-for="person in blackList" :key="person.account_id" class="item">
      <img :src="person.avatar" alt="avatar" width="50" height="50" />
      <div class="left">
          <p class="name">{{ person.name }}</p>
      </div>
      <div class="right">
          <button @click="Remove(person.account_id)">移出</button>
      </div>
    </div>
    <BlackListManagement
      :persons="managePesons"
      :type="managementType"
      v-show="isBlackListManagementVisible"
      @close="isBlackListManagementVisible = false"
      @confirm="confirmSelection"
      />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  </div>
</template>
  
<script>
import { removeFromBlackList, getBlackList } from '@/services/contactList';
import BlackListManagement from './BlackListManagement.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';

export default {
  components: {
    BlackListManagement,
    ContextMenu,
  },
  data() {
    return {
      // blackList: [
      //   {
      //     avatar: '',
      //     name: 'John Doe',
      //     account_id: '1',   // id
      //     signature:"爱拼才会赢",
      //   },
      //   {
      //     avatar: '',
      //     name: 'Jane Doe',
      //     account_id: '2',
      //     signature:"hi",
      //   },
      // ],
      blackList: [],
      isBlackListManagementVisible: false,
      managementType:'',
      managePesons: [
        {
          avatar: '',
          name: 'John Doe',
          account_id: '1',   // id
          signature:"爱拼才会赢",
        },
        {
          avatar: '',
          name: 'Jane Doe',
          account_id: '2',
          signature:"hi",
        },
      ],
      boundD: null, // 边界的坐标
			boundR: null, // 边界的坐标

    };
  },
  methods: {
    async fetchBlackList() {
      const response = await getBlackList();
      this.blackList = response.data;
    },
    async Remove(id) {
      const response = await removeFromBlackList(id);
      this.blackList = this.blackList.filter(person => person.account_id !== id);
    },
    showContextMenu(event){
      const items = [
        '批量移出',
        '批量移入',
      ];
      this.$refs.contextMenu.show(event, items, null, this.boundD, this.boundR);
    },
    handleMenuSelect(item){
      if(item === '批量移出') {
        this.isBlackListManagementVisible = true;
        this.managementType = 'out';
        // this.managePesons = this.blackList;  // todo api
      }
      else {
        this.isBlackListManagementVisible = true;
        this.managementType = 'in';
        // this.managePesons = this.blackList;  // todo api
      }
    },
    async confirmSelection(selectedPersons) {
      if(this.managementType === 'out') {
        for(const person of selectedPersons) {
          await removeFromBlackList(person.account_id);   // todo api
        }
      }
      else {
        for(const person of selectedPersons) {
          await addToBlackList(person.account_id);   // todo api
        }
      }
      this.fetchBlackList();
    },
  },
  created() {
    this.fetchBlackList();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>
button {
  margin-right: 5px;
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
button:first-of-type {
  background-color: #28a745;
  color: white;
}
button:last-of-type {
  background-color: #dc3545;
  color: white;
}
</style>