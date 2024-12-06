<template>
    <ToggleContent v-for="tag in tags" :previewText="tag">
        <div v-for="item in filteredItemsByTag(tag)" :key="items.account_id" class="item" @click="showProfileCard($event, item.account_id)">
            <img :src="item.avatar" alt="avatar" width="50" height="50" />
            <div v-if="this.type === 'friendList'" class="left">
                <p class="name">{{ item.remark }}</p>
                <p class="remark">{{ (`[${item.status}]签名：${item.signature}`) }}</p>
            </div>
            <div v-else class="left">
                <p class="name">{{ item.remark }}</p>
            </div>
        </div>
    </ToggleContent>

</template>

<script>
    import ToggleContent from '@/components/base/ToggleContent.vue';
    export default {
        components: {
            ToggleContent
        },
        props:['items','type','tags'],
        data() {
            return {
            
            };
        },
        methods: {
            filteredItemsByTag(tag) {
                return this.items.filter(item => item.tag === tag);
            },
            async fetchFriends() {
                const response = await getFriends();
                this.friends = response.data;
            },
            showProfileCard(event, tid){
                this.$emit('show-profile-card', event, tid);
            }
        },
        created() {
            this.fetchFriends();
        },
    };
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>