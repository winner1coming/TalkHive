<template>
  <div class="main">
    <div class="contact-header">
        黑名单
      <button style="float: right;" @click="showContextMenu">批量管理</button>
    </div>
    <div 
      v-for="person in blackList" 
      :key="person.account_id" 
      class="item"
    >
      <img :src="person.avatar" alt="avatar" width="50" height="50" @click="showProfileCard($event, person.account_id)"/>
      <div class="left" @click="showProfileCard($event, person.account_id)">
          <p class="name">{{ person.name?person.name:person.remark }}</p>
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
    <PersonProfileCard ref="profileCard" />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  </div>
</template>
  
<script>
import * as contactListAPI from '@/services/contactList';
import { getPersonProfileCard } from '@/services/api';


import BlackListManagement from './BlackListManagement.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import PersonProfileCard from '@/components/base/PersonProfileCard.vue';

export default {
  components: {
    BlackListManagement,
    PersonProfileCard,
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
      try{
        const response = await contactListAPI.getBlackList();
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.blackList = response.data.data;
      }catch(err){
        console.log(err);
      } 
    },
    async showProfileCard(event, send_account_id){
      try{
        const response = await getPersonProfileCard(send_account_id); 
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        const profile = response.data.data;
        this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
      }catch(err){
        console.log(err);
      }
      
    },
    async Remove(id) {
      try{
        const response = await contactListAPI.removeFromBlackList(id);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.blackList = this.blackList.filter(person => person.account_id !== id);
      }catch(err){
        console.log(err);
      }
    },
    showContextMenu(event){
      const items = [
        '批量移出',
        '批量移入',
      ];
      this.$refs.contextMenu.show(event, items, null, this.boundD, this.boundR);
    },
    async handleMenuSelect(item){
      if(item === '批量移出') {
        this.isBlackListManagementVisible = true;
        this.managementType = 'out';
        this.managePesons = this.blackList;  
      }
      else {
        this.isBlackListManagementVisible = true;
        this.managementType = 'in';
        try{
          const response = await contactListAPI.getFriends();
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.managePesons = response.data.data;
        }catch(err){
          console.log(err);
        }
      }
    },
    async confirmSelection(selectedPersons) {
      if(this.managementType === 'out') {
        for(const person of selectedPersons) {
          try{
            const response = await contactListAPI.removeFromBlackList(person.account_id);
            if(response.status !== 200){
              this.$root.notify(response.data.message, 'error');
            }
          }catch(err){
            console.log(err);
          }
        }
      }
      else {
        for(const person of selectedPersons) {
          try{
            const response = await contactListAPI.addToBlackList(person.account_id);
            if(response.status !== 200){
              this.$root.notify(response.data.message, 'error');
            }
          }catch(err){
            console.log(err);
          }
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