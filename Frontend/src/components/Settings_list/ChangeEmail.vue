<template>
    <div class="change-phone">
      <h3>更改邮箱</h3>
      <div class="text_group">
        <label for="oldPhone">原邮箱:</label>
        <input type="text" v-model="oldEmail" @input="validate_Email1" @blur="validate_oldEmail" :placeholder="oldEmail" />
      </div>
      <p v-if="errors.oldEmail" class="error">{{ errors.oldEmail }}</p>
      <div class="text_group">
        <label for="newPhone">新邮箱:</label>
        <input type="text" v-model="newEmail" @input="validate_Email2" :placeholder="newEmail" />
      </div>
      <p v-if="errors.newEmail" class="error">{{ errors.newEmail }}</p>
      <div class="verificate">
        <label for="nowCode">验证码:</label>
        <input type="text" v-model="nowCode" />
        <button  @click="sendCode">获取</button>
      </div>
    <p>注：验证发将发送到新邮箱中，请注意查收</p>
    <div class="button_container">
        <button @click="saveEmail">保存</button>
        <button @click="cancle" >取消</button>
    </div>
    </div>
  </template>
  
  <script>
  import {getCode,saveEmail} from '@/services/settingView.js';

  export default {
    props:{
      user: {
        type: Object,
        required: true,
      },
    },

    data() {
      return {
        oldEmail:'',
        newEmail: '',
        code:'',
        nowCode:'',
        errors: {
          oldEmail: '',
          newEmail: '',
          code:'',
        }
      };
    },

    methods: {
      //验证旧邮箱的格式
      validate_Email1(){
        if (!this.oldEmail) {
          this.errors.oldEmail = '邮箱不能为空';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.oldEmail)) {
          this.errors.oldEmail = '邮箱格式不正确';
        } else {
          this.errors.oldEmail = '';
        }
      },
      //验证新邮箱的格式
      validate_Email2(){
        if (!this.newEmail) {
          this.errors.newEmail = '邮箱不能为空';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.newEmail)) {
          this.errors.newEmail = '邮箱格式不正确';
        } else {
          this.errors.newEmail = '';
        }
      },

      //验证邮箱是否相等（验证与数据库的邮箱
      validate_oldEmail() {
        const email = this.user.email;
        if (this.oldEmail !== email) {
          this.errors.oldEmail = '原邮箱不正确';
        } else {
          this.errors.oldEmail = '';
        }
      },
      //验证验证码的正确性
      validateCode(){
          if(!this.nowCode){
            this.errors.code = '请输入验证码！';
          }
          else if(this.nowCode !== this.code){
            this.errors.code = '验证码不正确！';
          }else if(!this.code){
            this.code = '请先获取验证码~';
          }
          else{
            this.errors.code = '';
          }
      },
      //发送验证码——向后端发送获取验证码请求
      async sendCode(){
        try{
          const response = await getCode({
              new_email:this.newEmail,
          });
          if(response.success){
            alert('验证码已发送');
            this.code = response.code;
          }else{
            alert(response.message || '验证码发送失败');
          }
        }catch (error){
          console.error("验证码发送失败:",error);
        }
      },
      //按下保存的按钮
      async saveEmail() {
        this.validateCode();
        try{
          const response = await saveEmail({
            new_email:this.newEmail,
          });
          if(response.success){
            alert('邮箱更改成功');
            this.$emit('updateUser', { email: this.newEmail}); //修改父组件的email
            //更新缓存的内容
            let users = JSON.parse(localStorage.getItem('users')) || [];
            const userInfo = {
              email: this.newEmail,
            }

            const index = users.findIndex(user => user.email === this.oldEmail);
            if(index !== -1){
              users[index].email = userInfo.email;
            }
            
            localStorage.setItem('users',JSON.stringify(users));
          }else{
            alert(response.message || '邮箱更改失败');
          }
        }catch(error){
          console.error("更改邮箱失败:",error)
        }
      },

      async cancle(){
        this.$emit('updateUser', { email: this.oldEmail});
      }
    },
  };
  </script>
  
  <style scoped>
  .change-phone {
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap:20px;
  }

  h3{
    margin-top: 20px;
  }


p {
  font-size: 14px;
  margin: 0;
}

.text_group,
.verificate {
  display: grid;
  grid-template-columns: 80px 1fr auto;
  align-items: center;
  gap: 0px;
}

.text_group label,
.verificate label {
  text-align: right;
  margin-right: 15px;
}

.text_group input{
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.verificate input {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    width: 115px;
}

.verificate button {
  padding: 10px 10px;
  background-color: #42b983;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.verificate button:hover {
  background-color: #369f6d;
}

.button_container {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  width: 100%;
  max-width: 200px;
}

.button_container button {
  padding: 10px 20px;
  background-color: #42b983;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.button_container button:hover {
  background-color: #369f6d;
}
</style>