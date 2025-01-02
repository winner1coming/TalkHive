<template>
  <div class="main">
    <div class="contact-header">
        群聊列表
        <button 
          style="float: right;"
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
      :persons="groups"
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
import { getGroupProfileCard } from '@/services/api';

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
      type: 'groupList',  // friendList, groupList
      tags: [],  // 从后端获取
      items: [],   // 群组列表
      boundD: 0,
      boundR: 0,
      isDivideDeleteVisible: false,
      isDivideManagementVisible: false,
      isDivideAddVisible:false,
      isDivideMoveVisible: false,
     
      groups: [],    // 用于移入移出分组
      managementType: '',
      obj: null,
      selectedPersons: [],
    };
  },
  methods: {
    async fetchGroups() {
      try {
        const response = await contactListAPI.getGroups();
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.items = response.data.data;
      } catch (error) {
        console.log(error);
      }
    },
    async fetchTags() {
      try {
        const response = await contactListAPI.getDivides('groups');
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.tags = response.data.divides;
        this.tags.unshift('我创建的');
        this.tags.unshift('全部');
      } catch (error) {
        console.log(error);
      }
    },
    async showProfileCard(event, send_account_id){
      try{
        const response = await getGroupProfileCard(send_account_id); 
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        const profile = response.data.data;
        this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
      }catch(error){
        console.log(error);
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
      let items = [
        '移动',
      ];
      if(person.group_owner === this.$store.state.account_id){
        items.push('解散群聊');
      }else{
        items.push('退出群聊');
      }
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
        this.groups = this.items.filter(person => person.tag !== obj);
      }
      else if(item === '移出'){
        this.isDivideManagementVisible = true;
        this.managementType = 'out';
        this.obj = obj;   // tag
        this.groups = this.items.filter(person => person.tag === obj);
      }
      else if(item === '删除'){
        try{
          const response = await contactListAPI.deleteDivide('groups', obj);   // obj为分组名
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          } 
          this.tags = this.tags.filter(tag => tag !== obj);
          await this.fetchGroups();
        }catch(err){
          console.error(err);
        }
      }
      else if(item === '更名'){
        this.$refs.divideAdd.type = 'rename';
        this.isDivideAddVisible = true;
        this.obj = obj;
      }
      else if(item === '移动'){
        this.isDivideMoveVisible = true;
        this.$refs.divideMove.selectedDivide = obj.tag;
        this.$refs.divideMove.multiple = false;
        this.obj = obj;
      }else if(item === '解散群聊'){
        try{
          const response = await contactListAPI.dismissGroup(obj.account_id);   // obj为好友对象
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.items = this.items.filter(person => person.account_id !== obj.account_id);
        }catch(err){
          console.error(err);
        }
      }else if(item === '退出群聊'){
        try{
          const response = await contactListAPI.exitGroup(obj.account_id);   // obj为好友对象
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.items = this.items.filter(person => person.account_id !== obj.account_id);
        }catch(err){
          console.error(err);
        }
      }
    },
    async addDivide(newDivide){
      if(this.tags.includes(newDivide)){
        alert('分组名已存在');
        return;
      }
      try {
        const response = await contactListAPI.createDivide('groups', newDivide); 
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        } 
        this.tags.push(newDivide);
      } catch (error) {
        console.log(error);
      }
    },
    async renameDivide(newDivide){
      if(this.obj===newDivide){
        //alert('分组名已存在');
        return;
      }
      try {
        const response = await contactListAPI.renameDivide('groups', this.obj, newDivide);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        } 
        this.tags = this.tags.map(tag => tag === this.obj ? newDivide : tag); 
      } catch (error) {
        console.log(error);
      }
      await this.fetchGroups();
    },
    async deleteDivides(divides){
      divides.forEach(async (divide) => {
        try {
          const response = await contactListAPI.deleteDivide('groups', divide);   
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
        } catch (error) {
          console.log(error);
        }
      });
      await this.fetchTags();
      await this.fetchGroups();
    },
    async divideMove(divide){
      try {
        const response = await contactListAPI.moveInDivide('groups', this.obj.account_id,divide);   // obj为好友对象
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.obj.tag = divide;   // obj为好友对象
      } catch (error) {
        console.log(error);
      }
    },
    divideIn(selectedPersons){
      selectedPersons.forEach(async (person) => {
        try {
          const response = await contactListAPI.moveInDivide('groups', person.account_id,this.obj);   
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
          const response = await contactListAPI.moveInDivide('groups', person.account_id, divide);  
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
    this.fetchGroups();
    this.fetchTags();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
    this.fetchTags();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>


</style>