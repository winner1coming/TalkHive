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
      @show-devide-context-menu="showDevideContextMenu"
    />
    <ProfileCard ref="profileCard" />
    <DevideDelete
      :devides="tags"
      v-show="isDevideDeleteVisible"
      @close="isDevideDeleteVisible = false"
    />
    <DevideAdd
      v-show="isDevideAddVisible"
      @close="isDevideAddVisible = false"
    />
    <DevideManagement
      :type="managementType"
      :persons="groups"
      v-show="isDevideManagementVisible"
      @close="isDevideManagementVisible = false"
    />
    <DevideMove
      :devides="tags"
      :selectedTag="selectedTag"
      v-show="isDevideMoveVisible"
      @close="isDevideMoveVisible = false"
    />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  </div>
</template>

<script>
import { getGroups } from '@/services/contactList';
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
      tags: ['家人', '朋友', '同事'],  // 从后端获取
      // items: [   // 从后端获取
      //   {
      //     avatar: '',
      //     account_id: '1',   // 群id
      //     signature: '这是一个群聊',  // 群介绍
      //     remark: 'John',   // 群聊备注或群名称
      //     tag: '家人',
      //   },
      // ],
      items: [],
      boundD: 0,
      boundR: 0,
      isDevideManagementVisible: false,
      isDevideDeleteVisible: false,
      isDevideManagementVisible: false,
      isDevideAddVisible:false,
      isDevideMoveVisible: false,
      groups:[
				{
					acount_id: '13872131',   
					name: 'test',   // 备注或名称
					avatar: '',
					signature: '爱拼才会赢',    // 签名
				},
				{
					acount_id: '13872132',  
					name: 'test',
					avatar: '',
					signature: '爱拼才会赢',    // 签名
				},
				{
					acount_id: '13872133',   
					name: 'test',
					avatar: '',
					signature: '爱拼才会赢',    // 签名
				},
			],// todo 从后端获取（type为in时是除该分组外的所有好友，out时为当前分组内的好友）
      managementType: '',
      selectedDevide: null,
    };
  },
  methods: {
    async fetchGroups() {
      const response = await getGroups();
      this.items = response.data;
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
      if(item === '新建分组') this.isDevideAddVisible = true;
      else if(item === '删除分组') this.isDevideDeleteVisible = true;
      else if(item === '移入'){
        this.isDevideManagementVisible = true;
        this.managementType = 'in';
        // this.groups  // todo api
      }
      else if(item === '移出'){
        this.isDevideManagementVisible = true;
        this.managementType = 'out';
        // this.groups  // todo api
      }
      else if(item === '移动'){
        this.isDevideMoveVisible = true;
        this.selectedTag = obj.tag;
      }
    },
  },
  mounted() {
    this.fetchGroups();
    this.boundD = document.documentElement.clientHeight;
    this.boundR = document.documentElement.clientWidth;
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>