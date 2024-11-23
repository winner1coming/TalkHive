<template>
  <div class="main">
    <div class="contact-header">
        群聊通知
    </div>
    <div v-for="request in requests" :key="request.account_id" class="item">
      <img :src="request.avatar" alt="avatar" width="50" height="50" />
      <div class="left">
        <p class="name">{{ request.group_name }}</p>
        <div v-if="request.type === 'groupInvitation'" class="remark">
          {{(request.account_name+"邀请你加入群聊："+request.group_name)}}
        </div>
        <div v-else-if="request.type === 'groupApply'" class="remark">
          {{(request.account_name+"申请加入群聊："+request.group_name)}}
          <ToggleContent :previewText="'详情'">
            <p>备注: {{ request.reason }}</p>
          </ToggleContent>  
        </div>
        <div v-else-if="request.type === 'notification'" class="remark">
          <p >{{ ("您已被踢出群聊："+request.group_name) }}</p>
        </div>
      </div>
      <div class="right">
        <p class="remark">{{ request.time }}</p>
        <div v-if="request.status === 'pending' && request.type === 'groupInvitation'">
          <button @click="acceptInvitation(request.account_id,request.group_id)">同意</button>
          <button @click="rejectInvitation(request.account_id,request.group_id)">拒绝</button>
        </div>
        <div v-else-if="request.status === 'pending' && request.type === 'groupApply'">
          <button @click="acceptApply(request.account_id,request.group_id)">同意</button>
          <button @click="rejectApply(request.account_id,request.group_id)">拒绝</button>
        </div>
        <p v-else-if="request.status === 'accepted'">已同意</p>
        <p v-else-if="request.status === 'rejected'">已拒绝</p>
        <p v-else-if="request.status === 'waiting'">等待管理员处理</p>
      </div>
      
    </div>
  </div>
</template>

<script>
import { getGroupRequests, acceptGroupInvitationRequest, rejectGroupInvitationRequest, 
          acceptGroupApplyRequest,  rejectGroupApplyRequest} from '@/services/api.js';
import ToggleContent from '@/components/base/ToggleContent.vue';

export default {
  components: {
    ToggleContent
  },
  data() {
    return {
      requests: [
        {
          avatar: '',  // 群聊头像
          group_name: '相亲相爱一家人',    // 群名称
          account_name: 'alice',    // 好友备注或陌生人的名称（todo 后端应该有个逻辑判断是不是好友）
          account_id: '1',   // 申请者的id
          group_id: '1',   // 群聊id
          reason:"I want to be your friend",
          type: 'groupInvitation',  // （数据库中的apply_type）groupInvitation（邀请）,groupApply（申请）,notification
          status: 'pending',   // pending（需要我处理的）, accepted, rejected, waiting（对方处理中）,notification
          time: '2021-01-01 12:00:00',  // 待定
        },
        {
          avatar: '',
          group_name: '相亲相爱一家人',   
          account_name: 'Bob',    
          account_id: '1',  
          group_id: '1',   
          reason:"I want to be your friend",
          type: 'groupApply',
          status: 'pending',
          time: '2021-01-01 12:00:00',
        },
        {
          avatar: '',
          group_name: '相亲相爱一家人',
          account_name: 'Tom',
          account_id: '1',
          group_id: '1',
          reason:"I want to be your friend",
          type: 'notification',
          status: 'notification',
          time: '2021-01-01 12:00:00',
        }
      ],
    };
  },
  methods: {
    async fetchRequests() {
      const response = await getGroupRequests();
      this.requests = response.data;
    },
    async acceptInvitation(accountId,groupId) {
      await acceptGroupInvitationRequest(accountId,groupId);
      this.fetchRequests();
    },
    async rejectInvitation(accountId,groupId) {
      await rejectGroupInvitationRequest(accountId,groupId);
      this.fetchRequests();
    },
    async acceptApply(accountId,groupId) {
      await acceptGroupApplyRequest(accountId,groupId);
      this.fetchRequests();
    },
    async rejectApply(accountId,groupId) {
      await rejectGroupApplyRequest(accountId,groupId);
      this.fetchRequests();
    },
  },
  created() {
    this.fetchRequests();
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