<template>
  <div class="main">
    <div class="contact-header">
        好友列表
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
      @show-devide-context-menu="showDevideContextMenu"
    />
    <ProfileCard ref="profileCard" />
    <DevideDelete
      :devides="tags"
      v-show="isDevideDeleteVisible"
      @close="isDevideDeleteVisible = false"
      @delete-devides="deleteDevides"
    />
    <DevideAdd
      ref="devideAdd"
      v-show="isDevideAddVisible"
      @add-devide="addDevide"
      @rename-devide="renameDevide"
      @close="isDevideAddVisible = false"
    />
    <DevideManagement
      :type="managementType"
      :persons="persons"
      v-show="isDevideManagementVisible"
      @devide-in="devideIn"
      @devide-out="devideOut"
      @close="isDevideManagementVisible = false"
    />
    <DevideMove
      :devides="tags"
      v-show="isDevideMoveVisible"
      @devide-move="devideMove"
      @devides-move="devidesMove"
      @close="isDevideMoveVisible = false"
      ref="devideMove"
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
      type: 'friendList',  // friendList, groupList
      tags: ['家人', '朋友', '同事'],
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
      isDevideManagementVisible: false,
      isDevideDeleteVisible: false,
      isDevideManagementVisible: false,
      isDevideAddVisible:false,
      isDevideMoveVisible: false,
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
      const response = await contactListAPI.getFriends();
      this.items = response.data;
    },
    async fetchTags() {
      const response = await contactListAPI.getDevides('friends');
      this.tags = response.data.devides;
    },
    async showProfileCard(event, send_account_id){
      const response = await getProfileCard(send_account_id); 
      const profile = response.data;
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
        '拉黑',
      ];
      this.$refs.contextMenu.show(event, items, person, this.boundD, this.boundR, 60, 76);
    },
    async handleMenuSelect(item, obj=null){
      if(item === '新建分组') {
        this.$refs.devideAdd.type = 'add';
        this.isDevideAddVisible = true;
      }
      else if(item === '删除分组') this.isDevideDeleteVisible = true;
      else if(item === '移入'){
        this.isDevideManagementVisible = true;
        this.managementType = 'in';
        this.obj = obj;   // tag
        this.persons = this.items.filter(person => person.tag !== obj);
      }
      else if(item === '移出'){
        this.isDevideManagementVisible = true;
        this.managementType = 'out';
        this.obj = obj;   // tag
        this.persons = this.items.filter(person => person.tag === obj);
      }
      else if(item === '删除'){
        await contactListAPI.deleteDevide('friends', obj);   // obj为分组名
        this.tags = this.tags.filter(tag => tag !== obj);
        await this.fetchFriends();
      }
      else if(item === '更名'){
        this.$refs.devideAdd.type = 'rename';
        this.isDevideAddVisible = true;
        this.obj = obj;   // tag
      }
      else if(item === '移动'){
        this.isDevideMoveVisible = true;
        this.$refs.devideMove.selectedDevide = obj.tag;
        this.$refs.devideMove.multiple = false;
        this.obj = obj;
      }
      else if(item === '拉黑'){
        await contactListAPI.addToBlackList(obj.account_id);   // obj为好友对象
        this.items = this.items.filter(person => person !== obj);
      }
    },
    async addDevide(newDevide){
      if(this.tags.includes(newDevide)){
        alert('分组名已存在');
        return;
      }
      await contactListAPI.createDevide('friends', newDevide); 
      this.tags.push(newDevide);
    },
    async renameDevide(newDevide){
      if(this.obj===newDevide){
        //alert('分组名已存在');
        return;
      }
      await contactListAPI.renameDevide('friends', this.obj, newDevide); 
      this.tags = this.tags.map(tag => tag === this.obj ? newDevide : tag);
      await this.fetchFriends();
    },
    async deleteDevides(devides){
      devides.forEach(async (devide) => {
        try {
          await contactListAPI.deleteDevide('friends', devide);   
        } catch (error) {
          alert('删除分组失败！');
        }
      });
      await this.fetchTags();
      await this.fetchFriends();
    },
    async devideMove(devide){
      await contactListAPI.moveInDevide('friends', this.obj.account_id,devide);   // obj为好友对象
      this.obj.tag = devide;   // obj为好友对象
    },
    devideIn(selectedPersons){
      selectedPersons.forEach(async (person) => {
        try {
          await contactListAPI.moveInDevide('friends', person.account_id,this.obj);   
          person.tag = this.obj; 
        } catch (error) {
          alert('移入分组失败！');
        }
      });
    },
    devideOut(selectedPersons){
      this.isDevideMoveVisible = true;
      this.selectedPersons = selectedPersons;
      this.$refs.devideMove.selectedDevide = this.obj;
      this.$refs.devideMove.multiple = true;
    },
    devidesMove(devide){
      this.selectedPersons.forEach(async (person) => {
        try {
          await contactListAPI.moveInDevide('friends', person.account_id,devide);  
          person.tag = devide; 
        } catch (error) {
          alert('移入分组失败！');
        }
      });
    },
  },
  mounted() {
    this.fetchFriends();
    this.fetchTags();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>