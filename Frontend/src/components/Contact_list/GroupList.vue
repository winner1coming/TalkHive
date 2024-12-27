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
      @show-divide-context-menu="showDevideContextMenu"
    />
    <ProfileCard ref="profileCard" />
    <DevideDelete
      :divides="tags"
      v-show="isDevideDeleteVisible"
      @close="isDevideDeleteVisible = false"
      @delete-divides="deleteDevides"
    />
    <DevideAdd
      ref="divideAdd"
      v-show="isDevideAddVisible"
      @add-divide="addDevide"
      @rename-divide="renameDevide"
      @close="isDevideAddVisible = false"
    />
    <DevideManagement
      :type="managementType"
      :persons="groups"
      v-show="isDevideManagementVisible"
      @divide-in="divideIn"
      @divide-out="divideOut"
      @close="isDevideManagementVisible = false"
    />
    <DevideMove
      :divides="tags"
      v-show="isDevideMoveVisible"
      @divide-move="divideMove"
      @divides-move="dividesMove"
      @close="isDevideMoveVisible = false"
      ref="divideMove"
    />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  </div>
</template>

<script>
import * as contactListAPI from '@/services/contactList';
import { getProfileCard } from '@/services/api';

import itemList from './itemList.vue';
import DevideDelete from './DevideDelete.vue';
import DevideAdd from './DevideAdd.vue';
import DevideManagement from './DevideManagement.vue';
import DevideMove from './DevideMove.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import ProfileCard from '@/components/base/ProfileCard.vue';

export default {
  components: {
    itemList,
    ProfileCard,
    DevideDelete,
    DevideAdd,
    DevideManagement,
    DevideMove,
    ContextMenu,
  },
  data() {
    return {
      type: 'groupList',  // friendList, groupList
      tags: [],  // 从后端获取
      // items: [   // 从后端获取
      //   {
      //     avatar: '',
      //     account_id: '1',   // 群id
      //     signature: '这是一个群聊',  // 群介绍
      //     remark: 'John',   // 群聊备注或群名称
      //     tag: '家人',
      //   },
      // ],
      items: [],   // 群组列表
      boundD: 0,
      boundR: 0,
      isDevideManagementVisible: false,
      isDevideDeleteVisible: false,
      isDevideManagementVisible: false,
      isDevideAddVisible:false,
      isDevideMoveVisible: false,
      // groups:[
			// 	{
			// 		acount_id: '13872131',   
			// 		name: 'test',   // 备注或名称
			// 		avatar: '',
			// 		signature: '爱拼才会赢',    // 签名
			// 	},
			// 	{
			// 		acount_id: '13872132',  
			// 		name: 'test',
			// 		avatar: '',
			// 		signature: '爱拼才会赢',    // 签名
			// 	},
			// 	{
			// 		acount_id: '13872133',   
			// 		name: 'test',
			// 		avatar: '',
			// 		signature: '爱拼才会赢',    // 签名
			// 	},
			// ],// （type为in时是除该分组外的所有好友，out时为当前分组内的好友）
      groups: [],    // 用于移入移出分组
      managementType: '',
      obj: null,
      selectedPersons: [],
    };
  },
  methods: {
    async fetchGroups() {
      const response = await contactListAPI.getGroups();
      this.items = response.data.data;
    },
    async fetchTags() {
      const response = await contactListAPI.getDevides('friends');
      this.tags = response.data.data.divides;
      this.tags.unshift('全部');
    },
    async showProfileCard(event, send_account_id){
      const response = await getProfileCard(send_account_id); 
      const profile = response.data.data;
      this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
    },
    showContextMenu(event){
      const items = [
        '新建分组',
        '删除分组',
      ];
      this.$refs.contextMenu.show(event, items, null, this.boundD, this.boundR);
    },
    showDevideContextMenu(event, tag){
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
      ];
      this.$refs.contextMenu.show(event, items, person, this.boundD, this.boundR, 60, 76);
    },
    async handleMenuSelect(item, obj=null){
      if(item === '新建分组') {
        this.$refs.divideAdd.type = 'add';
        this.isDevideAddVisible = true;
      }
      else if(item === '删除分组') this.isDevideDeleteVisible = true;
      else if(item === '移入'){
        this.isDevideManagementVisible = true;
        this.managementType = 'in';
        this.obj = obj;   // tag
        this.groups = this.items.filter(person => person.tag !== obj);
      }
      else if(item === '移出'){
        this.isDevideManagementVisible = true;
        this.managementType = 'out';
        this.obj = obj;   // tag
        this.groups = this.items.filter(person => person.tag === obj);
      }
      else if(item === '删除'){
        await contactListAPI.deleteDevide('groups', obj);   // obj为分组名
        this.tags = this.tags.filter(tag => tag !== obj);
        await this.fetchGroups();
      }
      else if(item === '更名'){
        this.$refs.divideAdd.type = 'rename';
        this.isDevideAddVisible = true;
        this.obj = obj;
      }
      else if(item === '移动'){
        this.isDevideMoveVisible = true;
        this.$refs.divideMove.selectedDevide = obj.tag;
        this.$refs.divideMove.multiple = false;
        this.obj = obj;
      }
    },
    async addDevide(newDevide){
      if(this.tags.includes(newDevide)){
        alert('分组名已存在');
        return;
      }
      await contactListAPI.createDevide('groups', newDevide); 
      this.tags.push(newDevide);
    },
    async renameDevide(newDevide){
      if(this.obj===newDevide){
        //alert('分组名已存在');
        return;
      }
      await contactListAPI.renameDevide('groups', this.obj, newDevide);
      this.tags = this.tags.map(tag => tag === this.obj ? newDevide : tag); 
      await this.fetchGroups();
    },
    async deleteDevides(divides){
      divides.forEach(async (divide) => {
        try {
          await contactListAPI.deleteDevide('groups', divide);   
        } catch (error) {
          alert('删除分组失败！');
        }
      });
      await this.fetchTags();
      await this.fetchGroups();
    },
    async divideMove(divide){
      await contactListAPI.moveInDevide('groups', this.obj.account_id,divide);   // obj为好友对象
      this.obj.tag = divide;   // obj为好友对象
    },
    divideIn(selectedPersons){
      selectedPersons.forEach(async (person) => {
        try {
          await contactListAPI.moveInDevide('groups', person.account_id,this.obj);   
          person.tag = this.obj; 
        } catch (error) {
          alert('移入分组失败！');
        }
      });
    },
    divideOut(selectedPersons){
      this.isDevideMoveVisible = true;
      this.selectedPersons = selectedPersons;
      this.$refs.divideMove.selectedDevide = this.obj;
      this.$refs.divideMove.multiple = true;
    },
    dividesMove(divide){
      this.selectedPersons.forEach(async (person) => {
        try {
          await contactListAPI.moveInDevide('groups', person.account_id, divide);  
          person.tag = divide;
        } catch (error) {
          alert('移入分组失败！');
      }
      });
    },
  },
  mounted() {
    this.fetchGroups();
    this.fetchTags();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>