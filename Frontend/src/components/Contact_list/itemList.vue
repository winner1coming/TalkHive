<template>
	<HeaderToggle v-for="tag in tags" :previewText="tag">
		<div v-for="item in filteredItemsByTag(tag)" :key="items.account_id" class="item" @click="showProfileCard($event, item.account_id)">
			<img :src="item.avatar" alt="avatar" width="50" height="50" />
			<div class="left">
				<p class="name">{{ item.remark }}</p>
				<p class="remark" v-if="this.type === 'friendList'">{{ (`[${item.status}]签名：${item.signature}`) }}</p>
				<p class="remark" v-else>{{ item.signature }}</p>
			</div>
		</div>
	</HeaderToggle>
</template>

<script>
import HeaderToggle from '@/components/Contact_list/HeaderToggle.vue';
export default {
	components: {
		HeaderToggle
	},
	props:['items','type','tags'],
	methods: {
		filteredItemsByTag(tag) {
			console.log(this.items);
			return this.items.filter(item => item.tag === tag);
		},
		showProfileCard(event, tid){
			this.$emit('show-profile-card', event, tid);
		}
	},
};

</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>