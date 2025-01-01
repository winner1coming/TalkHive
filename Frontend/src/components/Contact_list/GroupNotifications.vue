<template>
  <div class="main">
    <div class="contact-header">
        群聊通知
    </div>
    <div v-for="request in requests" :key="request.apply_id" class="item">
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
          <p >{{ (request.account_name+"已退出群聊："+request.group_name) }}</p>
        </div>
      </div>
      <div class="right">
        <p class="remark">{{ request.time }}</p>
        <div v-if="request.status === 'pending' && request.type === 'groupInvitation' && request.receiver_id === this.$store.state.user.id">
          <button @click="acceptInvitation(request.account_id,request.group_id)">同意</button>
          <button @click="rejectInvitation(request.account_id,request.group_id)">拒绝</button>
        </div>
        <p v-else-if="request.status === 'pending' && request.type === 'groupInvitation'">等待对方处理</p>
        <div v-else-if="request.status === 'pending' && request.type === 'groupApply' && request.receiver_id === this.$store.state.user.id">
          <button @click="acceptApply(request.account_id,request.group_id)">同意</button>
          <button @click="rejectApply(request.account_id,request.group_id)">拒绝</button>
        </div>
        <p v-else-if="request.status === 'pending' && request.type === 'groupApply'">等待管理员处理</p>
        <p v-else-if="request.status === 'accepted'">已同意</p>
        <p v-else-if="request.status === 'rejected'">已拒绝</p>
      </div>
      
    </div>
  </div>
</template>

<script>
import ToggleContent from '@/components/base/ToggleContent.vue';
import { getGroupRequests, groupInvitationRequestPend, groupApplyRequestPend} from '@/services/contactList';
const ContactListAPI = {
  getGroupRequests,
  groupInvitationRequestPend,
  groupApplyRequestPend
};

export default {
  components: {
    ToggleContent
  },
  data() {
    return {
      // requests: [
      //   {
      //     apply_id: '1',  // 申请id
      //     avatar: '',  // 群聊头像
      //     group_name: '相亲相爱一家人',    // 群名称
      //     account_name: 'alice',    // 好友备注或陌生人的名称（todo 后端应该有个逻辑判断是不是好友）
      //     sender_id: '1',  // 申请者的tid
      //     receiver_id: '2',   // 接收者的tid
      //     group_id: '1',   // 群聊id
      //     reason:"I want to be your friend",
      //     type: 'groupInvitation',  // （数据库中的apply_type）groupInvitation（邀请）,groupApply（申请）,notification
      //     status: 'pending',   // pending（需要我处理的）, accepted, rejected, waiting（对方处理中）,notification
      //     time: '2021-01-01 12:00:00',  // 待定
      //   },
      //   {
      //     apply_id: '2',  // 申请id
      //     avatar: '',
      //     group_name: '相亲相爱一家人',   
      //     account_name: 'Bob',    
      //     sender_id: '2',  // 申请者的tid
      //     receiver_id: '1',   // 接收者的tid
      //     group_id: '1',   
      //     reason:"I want to be your friend",
      //     type: 'groupApply',
      //     status: 'pending',
      //     time: '2021-01-01 12:00:00',
      //   },
      //   {
      //     apply_id: '3',  // 申请id
      //     avatar: '',
      //     group_name: '相亲相爱一家人',
      //     account_name: 'Tom',
      //     sender_id: '1',  // 申请者的tid
      //     receiver_id: '2',   // 接收者的tid
      //     group_id: '1',
      //     reason:"I want to be your friend",
      //     type: 'notification',
      //     status: 'notification',
      //     time: '2021-01-01 12:00:00',
      //   }
      // ],
      requests: [],
    };
  },
  methods: {
    async fetchRequests() {
      try{
        const response = await ContactListAPI.getGroupRequests();
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.requests = response.data.data;
      }catch(err){
        console.error(err);
      }
    },
    async acceptInvitation(accountId,groupId) {
      try{
        const response = await ContactListAPI.groupInvitationRequestPend(accountId,groupId, true);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.fetchRequests();
      }catch(err){
        console.error(err);
      }
    },
    async rejectInvitation(accountId,groupId) {
      try{
        const response = await ContactListAPI.groupInvitationRequestPend(accountId,groupId, false);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.fetchRequests();
      }catch(err){
        console.error(err);
      }
    },
    async acceptApply(accountId,groupId) {
      try{
        const response = await ContactListAPI.groupApplyRequestPend(accountId, groupId, true);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.fetchRequests();
      }catch(err){
        console.error(err);
      }
    },
    async rejectApply(accountId,groupId) {
      try{
        const response = await ContactListAPI.groupApplyRequestPend(accountId, groupId, false);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.fetchRequests();
      }catch(err){
        console.error(err);
      }
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