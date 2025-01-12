<template>
  <div class="main">
    <div class="contact-header">
        好友列表
        <button 
          class="header-button"
          @click="showContextMenu"
        >分组管理</button>
    </div>
    <itemList 
      :items="items" 
      :type="type" 
      :tags="tags" 
      @show-profile-card="showProfileCard"
      @show-person-context-menu="showPersonContextMenu"
      @show-divide-context-menu="showDivideContextMenu"
    />
    <PersonProfileCard ref="profileCard" />
    <DivideDelete
      :divides="tags.filter(tag => tag !== '全部')"
      v-show="isDivideDeleteVisible"
      @close="isDivideDeleteVisible = false"
      @delete-divides="deleteDivides"
    />
    <DivideAdd
      ref="divideAdd"
      v-show="isDivideAddVisible"
      @add-divide="addDivide"
      @rename-divide="renameDivide"
      @close="isDivideAddVisible = false"
    />
    <DivideManagement
      :type="managementType"
      :persons="persons"
      v-show="isDivideManagementVisible"
      @divide-in="divideIn"
      @divide-out="divideOut"
      @close="isDivideManagementVisible = false"
    />
    <DivideMove
      :divides="tags.filter(tag => tag !== '全部')"
      v-show="isDivideMoveVisible"
      @divide-move="divideMove"
      @divides-move="dividesMove"
      @close="isDivideMoveVisible = false"
      ref="divideMove"
    />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  </div>
</template>

<script>
import * as contactListAPI from '@/services/contactList';
import { EventBus } from '@/components/base/EventBus';
import { getPersonProfileCard } from '@/services/api';

import itemList from './itemList.vue';
import DivideDelete from './DivideDelete.vue';
import DivideAdd from './DivideAdd.vue';
import DivideManagement from './DivideManagement.vue';
import DivideMove from './DivideMove.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import PersonProfileCard from '@/components/base/PersonProfileCard.vue';


export default {
  components: {
    itemList,
    PersonProfileCard,
    DivideDelete,
    DivideAdd,
    DivideManagement,
    DivideMove,
    ContextMenu,
  },
  data() {
    return {
      type: 'friendList',  // friendList, groupList
      tags: [],  // 从后端获取
      // items: [
      //   {
      //     avatar: '',
      //     account_id: '1',
      //     remark: 'John',   // 好友备注
      //     status: 'online',   // online, offline
      //     signature: '爱拼才会赢',    // 签名
      //     tag: '家人',   
      //   },
      // ],
      items: [],   // 好友列表
      boundD: 0,
      boundR: 0,
      isDivideDeleteVisible: false,
      isDivideManagementVisible: false,
      isDivideAddVisible:false,
      isDivideMoveVisible: false,
      // persons:[
			// ],//  （type为in时是除该分组外的所有好友，out时为当前分组内的好友）
      persons: [],  // 用于移入移出分组
      managementType: '',
      obj: null,
      selectedPersons: [],
    };
  },
  methods: {
    async fetchFriends() {
      try{
        const response = await contactListAPI.getFriends();
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.items = response.data.data;
      }
      catch(err){
        console.log(err);
      }
    },
    async fetchTags() {
      try{
        const response = await contactListAPI.getDivides('friends');
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.tags = response.data.divides;
        this.tags.unshift('全部');
      }
      catch(err){
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
    showContextMenu(event){
      const items = [
        '新建分组',
        '删除分组',
      ];
      this.$refs.contextMenu.show(event, items, null, this.boundD, this.boundR);
    },
    showDivideContextMenu(event, tag){
      const items = [
        '移入',
        '移出',
        '删除',
        '更名',
      ];
      this.$refs.contextMenu.show(event, items, tag, this.boundD, this.boundR, 60, 76);
    },
    showPersonContextMenu(event, person){  
      const items = [
        '移动',
        '拉黑',
        '删除好友',
      ];
      this.$refs.contextMenu.show(event, items, person, this.boundD, this.boundR, 60, 76);
    },
    async handleMenuSelect(item, obj=null){
      if(item === '新建分组') {
        this.$refs.divideAdd.type = 'add';
        this.isDivideAddVisible = true;
      }
      else if(item === '删除分组') this.isDivideDeleteVisible = true;
      else if(item === '移入'){
        this.isDivideManagementVisible = true;
        this.managementType = 'in';
        this.obj = obj;   // tag
        this.persons = this.items.filter(person => person.tag !== obj);
      }
      else if(item === '移出'){
        this.isDivideManagementVisible = true;
        this.managementType = 'out';
        this.obj = obj;   // tag
        this.persons = this.items.filter(person => person.tag === obj);
      }
      else if(item === '删除'){
        try {
          const response = await contactListAPI.deleteDivide('friends', obj);   // obj为分组名
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.tags = this.tags.filter(tag => tag !== obj);
          await this.fetchFriends();
        } catch (error) {
          console.log(error);
        }
      }
      else if(item === '更名'){
        this.$refs.divideAdd.type = 'rename';
        this.isDivideAddVisible = true;
        this.obj = obj;   // tag
      }
      else if(item === '移动'){
        this.isDivideMoveVisible = true;
        this.$refs.divideMove.selectedDivide = obj.tag;
        this.$refs.divideMove.multiple = false;
        this.obj = obj;
      }
      else if(item === '拉黑'){
        try{
          const response = await contactListAPI.addToBlackList(obj.account_id);   // obj为好友对象
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.items = this.items.filter(person => person !== obj);
        }catch(err){
          console.log(err);
        }
      }else if(item === '删除好友'){
        try{
          const response = await contactListAPI.deleteFriend(obj.account_id);  
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.items = this.items.filter(person => person !== obj);
        }catch(err){
          console.log(err);
        }
      }
    },
    async addDivide(newDivide){
      if(this.tags.includes(newDivide)){
        alert('分组名已存在');
        return;
      }
      try{
        const response = await contactListAPI.createDivide('friends', newDivide); 
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.tags.push(newDivide);
      }catch(err){
        console.log(err);
      }
    },
    async renameDivide(newDivide){
      if(this.obj===newDivide){
        //alert('分组名已存在');
        return;
      }
      const response = await contactListAPI.renameDivide('friends', this.obj, newDivide); 
      if(response.status !== 200){
        this.$root.notify(response.data.message, 'error');
        return;
      }
      this.tags = this.tags.map(tag => tag === this.obj ? newDivide : tag);
      await this.fetchFriends();
    },
    async deleteDivides(divides){
      divides.forEach(async (divide) => {
        try {
          const response = await contactListAPI.deleteDivide('friends', divide);   
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
        } catch (error) {
          console.log(error);
        }
      });
      await this.fetchTags();
      await this.fetchFriends();
    },
    async divideMove(divide){
      try{
        const response = await contactListAPI.moveInDivide('friends', this.obj.account_id,divide);   // obj为好友对象
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.obj.tag = divide;   // obj为好友对象
      }
      catch(err){
        console.log(err);
      }
    },
    divideIn(selectedPersons){
      selectedPersons.forEach(async (person) => {
        try {
          const response = await contactListAPI.moveInDivide('friends', person.account_id,this.obj);   
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          person.tag = this.obj; 
        } catch (error) {
          console.log(error);
        }
      });
    },
    divideOut(selectedPersons){
      this.isDivideMoveVisible = true;
      this.selectedPersons = selectedPersons;
      this.$refs.divideMove.selectedDivide = this.obj;
      this.$refs.divideMove.multiple = true;
    },
    dividesMove(divide){
      this.selectedPersons.forEach(async (person) => {
        try {
          const response = await contactListAPI.moveInDivide('friends', person.account_id,divide);  
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          person.tag = divide; 
        } catch (error) {
          console.log(error);
        }
      });
    },
  },
  mounted() {
    this.fetchFriends();
    this.fetchTags();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
    EventBus.on('updateFriendList', (newFriend) => {
      this.items.unshift(newFriend);
    });
  },
  beforeUnmount() {
    EventBus.off('updateFriendList');
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>
.header-button{
  float: right;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: var(--font-size-small);
}
</style>