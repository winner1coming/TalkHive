<template>
  <div 
    v-for="item in items" 
    :key="item.account_id" 
    class="item" 
    @click="showProfileCard($event, item.account_id)"
  >
    <img :src="item.avatar" alt="avatar" width="50" height="50" />
    <div class="left">
      <p class="name">{{ item.remark }}</p>
      <p class="remark" v-if="item.type === 'friends'">{{ (`[${item.status}]签名：${item.signature}`) }}</p>
      <p class="remark" v-else>{{ item.signature }}</p>
    </div>
  </div>

</template>

<script>
import * as contactListAPI from "@/services/contactList";
export default {
	props:['keyword'],
  data() {
    return {
      items: [],
    };
  },
  watch: {
    keyword: {
      async handler(val){
        if(!val){
          //this.$emit('stop-search');
          return;
        }
        try{
          const response = await contactListAPI.searchContacts(val);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          this.items = response.data.data;
        }catch(e){
          console.error(e);
        }
      },
      immediate: true,
    }
  },
	methods: {
		showProfileCard(event, tid){
			this.$emit('show-profile-card', event, tid);
		},
	},
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>