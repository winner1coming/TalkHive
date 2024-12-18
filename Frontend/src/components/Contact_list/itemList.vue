<template>
	<HeaderToggle v-for="tag in tags" :previewText="tag" @manage-devide="manageDevide($event, tag)">
		<div 
			v-for="item in filteredItemsByTag(tag)" 
			:key="item.account_id" 
			class="item" 
			@click="showProfileCard($event, item.account_id)"
			@contextmenu="showPersonContextMenu($event, item)"
		>
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
			return this.items.filter(item => item.tag === tag);
		},
		showProfileCard(event, tid){
			this.$emit('show-profile-card', event, tid);
		},
		showPersonContextMenu(event, item){
			this.$emit('show-person-context-menu', event, item);
		},
		manageDevide(event, tag){
			this.$emit('show-devide-context-menu', event, tag);
		}
	},
};
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>