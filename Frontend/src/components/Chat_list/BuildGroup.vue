<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <div>
        <form @submit.prevent="createNewGroup">
          <!-- 上传头像 -->
          <div class="avatar-container">
            <img class="head-avatar" :src="group_avatar" @click="triggerFileInput" />
            <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" accept="image/*" />
          </div>
          <div>
            <p class="title">群聊名称：</p>
            <textarea 
              class="group-name" 
              v-model="group_name" 
              placeholder="不超过10个字" 
              maxlength="10"
              required 
            />
          </div>
          <div>
            <p class="title">群介绍：</p>
            <textarea 
              class="discription" 
              v-model="group_description" 
              placeholder="不超过100个字"
              maxlength="100"
            />
          </div>
          <hr class="divider" />
          <div>
              <p class="title">申请入群的方式：</p>
              <p class="detail" style="margin-left: 15px;">成员邀请: <SwitchButton v-model="allow_invite" @change-value="allow_invite=!allow_invite"/></p>
              <p class="detail" style="margin-left: 15px;">群号搜索: <SwitchButton v-model="allow_id_search" @change-value="allow_id_search=!allow_id_search"/></p>
              <p class="detail" style="margin-left: 15px;">群名称搜索: <SwitchButton v-model="allow_name_search" @change-value="allow_name_search!=allow_name_search"/></p>
          </div>
          <button type="submit" class="submit-button">创建</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { createGroup } from '@/services/contactList';
import SwitchButton from '@/components/base/SwitchButton.vue';
export default {
  components: { SwitchButton },
  data() {
    return {
      group_name: '',
      group_avatar:'',
      group_description: '',
      allow_invite: true,
      allow_id_search: true,
      allow_name_search: true,
    };
  },
  methods: {
    // 触发文件输入
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    // 处理文件选择
    handleFileChange(event) {
      const file = event.target.files[0];
      console.log(file);
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.group_avatar = e.target.result;
        };
        reader.readAsDataURL(file);
        console.log(this.group_avatar);
      }
    },
    async createNewGroup(){
      try {
        const response = await createGroup(this.group_name,this.group_avatar,this.group_description,this.allow_invite,this.allow_id_search,this.allow_name_search);
        if (response.status===200) {
          this.$emit('group-created');
          this.close();
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      } catch (error) {
        console.error('Error creating group:', error);
      }
    },
    close() {
      this.$emit('close');
    },
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000; /* 确保在最上层 */
}

.modal-content {
  background-color: var(--background-color);
  color: var(--text-color);
  padding: 20px;
  border-radius: 8px;
  width: 300px;
  height: 500px;
}

.avatar-container {
  position: relative;
  cursor: pointer;
}
.head-avatar {
  width: 80px;
  height: 80px;
  margin-top: 10px;
  margin-bottom: 20px;
  margin-left: 10px;
  border-radius: 100%;
  cursor: pointer;
}
.group-name {
  width: 100%;
  height: 25px;
  margin-top: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  padding: 5px;
  resize: none;
}
.discription {
  width: 100%;
  height: 100px;
  margin-top: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  padding: 5px;
  text-align: start;
  vertical-align: top;
  resize: none;
}

.divider {
  border: 0;
  height: 1px;
  background: #e0e0e0;
  margin: 10px 0;
}
.title {
  text-align: left;
  font-weight: 500;
}
.detail {
  text-align: left;
  color: #888;
}

.submit-button{
  width: 100%;
  padding: 10px;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border: none;
  border-radius: 5px;
  cursor: pointer;
  text-align: center;
  margin-top: 20px;
}
.submit-button:hover {
  background-color: var(--button-background-color1);
}

</style>