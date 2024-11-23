<template>
    <ToggleContent  :previewText="'tag'">
        <div v-for="item in items" :key="items.account_id" class="item">
            <img :src="item.avatar" alt="avatar" width="50" height="50" />
            <div v-if="this.type === 'friendList'" class="left">
                <p class="name">{{ item.remark }}</p>
                <p class="remark">{{ (`[${item.status}]签名：${item.signature}`) }}</p>
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
            async fetchFriends() {
            const response = await getFriends();
            this.friends = response.data;
            },
        },
        created() {
            this.fetchFriends();
        },
    };
</script>

<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>

</style>