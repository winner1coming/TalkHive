<template>
    <div class="permission">
      <h2>添加我的方式</h2>
      <div class="user-info">
        <span class="user-id">ID: {{ ID }}</span>
        <button class="toggle-button" :class="{ 'status-on': idStatus === 'on' }" @click="toggleIdStatus">
          <span class="toggle-circle"></span>
        </button>
        <p>{{ msg1 }}</p>
      </div>
      <div class="user-info">
        <span class="user-phone">昵称: {{ nickname }}</span>
        <button class="toggle-button" :class="{ 'status-on': nicknameStatus === 'on' }" @click="toggleNicknameStatus">
          <span class="toggle-circle"></span>
        </button>
        <p>{{ msg2 }}</p>
      </div>
      <Windows :visible="modalVisible" :message="modalMessage" @close="closeModal" />
    </div>
  </template>
  
  <script>
  import {isIDAdd, isNicknameAdd} from '@/services/api.js';
  import Windows from '@/components/base/Windows.vue'

  export default {
    components:{
      Windows,
    },
    //从父组件获取信息
    props:{
      users: {
        type: Object,
        required: true,
        default: () => ({
            ID:'',
            friend_permissionID: 'off',
            friend_permissionNickname: 'off',
          }),
        validator: (value) => {
          return ['ID','friend_permissionID', 'friend_permissionNickname'].every(key => value.hasOwnProperty(key));
        },
      },
    },

    data(){
      return {
        ID:this.$store.state.user.id,
        nickname:this.$store.state.user.username,
        idStatus: this.users.friend_permissionID, // 'on' or 'off'
        nicknameStatus: this.users.friend_permissionNickname, // 'on' or 'off'
        modalVisible: false,
        modalMessage: '',
        msg1:'',
        msg2:'',
      };
    },
    methods: {
      sleep(ms) {
      return new Promise(resolve => setTimeout(resolve, ms));
    },
      //设置ID的权限
      async toggleIdStatus() {
        const jud = this.idStatus === 'on' ? true : false;
        try{
          const response = await isIDAdd({id:this.ID,friend_permissionID :jud});
          if(response.success){
            this.msg1 = 'ID权限设置成功！';
            this.idStatus = this.idStatus === 'on' ? 'off' : 'on';
          }
          else{
            this.showModal(response.message);
          }
        }catch(error){
          this.showModal("服务器崩掉啦？");
          console.error("设置ID权限失败");
        }
      },
      //设置昵称的权限
      async toggleNicknameStatus() {
        const jud = this.nicknameStatus === 'on' ? true : false;
        try{
          const response = await isNicknameAdd({id:this.ID,friend_permissionNickname :jud});
          if(response.success){
            this.msg2 = '昵称权限设置成功！';
            this.nicknameStatus = this.nicknameStatus === 'on' ? 'off' : 'on';
          }
          else{
            this.showModal(response.message||"设置昵称权限失败");
          }
        }catch(error){
          console.error("设置昵称权限失败");
          this.showModal("服务器崩掉啦？");
        }
      },
      //展示弹窗
      showModal(message) {
        this.modalMessage = message;
        this.modalVisible = true;
      },
      //关闭弹窗
      closeModal() {
        this.modalVisible = false;
      },
    },
  };
  </script>
  
  <style scoped>
  .permission {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 20px;
  }
  
  .user-id, .user-phone {
    font-size: 16px;
  }
  
  .toggle-button {
    width: 60px;
    height: 30px;
    border-radius: 15px; /* 调整为圆角矩形 */
    background-color: gray; /* 默认关闭状态 */
    border: none;
    cursor: pointer;
    position: relative;
    transition: background-color 0.3s;
  }
  
  .toggle-button.status-on {
    background-color: green; /* 开启状态 */
  }
  
  .toggle-circle {
    width: 26px;
    height: 26px;
    border-radius: 50%;
    background-color: white;
    position: absolute;
    top: 2px;
    left: 2px;
    transition: transform 0.3s;
  }
  
  .toggle-button.status-on .toggle-circle {
    transform: translateX(30px); /* 滑动到右边 */
  }
  </style>