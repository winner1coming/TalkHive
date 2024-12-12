<template>
  <!-- 注册页面容器 -->
  <div class="register">
    <div class="container">
    <!-- 页面标题 -->
    <h2>注册</h2>

    <div class="avatar-container">
      <img class="headavatar" :src="avatar" @click="triggerFileInput" />
      <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" accept="image/*" />
    </div>

    <!--性别填写（同样设置为必填项）-->
    <div class="input-group">
      <label>性别*:</label>
      <div class="radio-group">
        <label> <input type="radio" v-model="gender" value="male" /> 男</label>
        <label> <input type="radio" v-model="gender" value="female" /> 女</label>
      </div>
    </div>
    <p v-if="errors.gender" class="error">{{ errors.gender }}</p>

    <!-- ID 输入框 -->
    <div class="input-group">
      <label for="id">账号*:</label>
      <input id="id" type="text" v-model="id" placeholder="" @blur="validateId" />
    </div>
    <p v-if="errors.id" class="error">{{ errors.id }}</p>
    
    
    <!-- 昵称输入框 -->
    <div class="input-group">
      <label for="nickname">昵称*:</label>
      <input id="nickname" type="text" v-model="nickname" placeholder="" @blur="validateNickname" />
    </div>
    <p v-if="errors.nickname" class="error">{{ errors.nickname }}</p>
    
    <!--生日输入-->
    <div class="input-group">
      <label for="birthday">生日*:</label>
      <div class="birthday-input-group">
        <!-- 日期选择器 -->
        <input
          type="date"
          v-model="birthday"
          id="birthday"
          placeholder=" "
          @change="validateBirthday"
          class="date-picker"
        />
      </div>
    </div>
    <p v-if="errors.birthday" class="error">{{ errors.birthday }}</p>

    <!-- 手机号输入框 -->
    <div class="input-group">
    <label for="phoneNumber">手机号:</label>
    <input id="phoneNumber" type="text" v-model="phoneNumber" placeholder="" @blur="validatePhoneNumber" />
    </div>
    <p v-if="errors.phoneNumber" class="error">{{ errors.phoneNumber }}</p>
 

    <!-- 邮箱输入框 -->
    <div class="input-group">
      <label for="email">邮箱*:</label>
      <input id="email" type="text" v-model="email" placeholder="" @blur="validateEmail" />
    </div>
    <p v-if="errors.email" class="error">{{ errors.email }}</p>
    
    <!-- 密码输入框 -->
    <div class="input-group">
      <label for="password">密码*:</label>
      <input id="password" type="password" v-model="password" placeholder="" @blur="validatePassword" />
    </div>
    <p v-if="errors.password" class="error">{{ errors.password }}</p>
    
    <!-- 确认密码输入框 -->
    <div class="input-group">
      <label for="confirmPassword">确认密码*:</label>
      <input id="confirmPassword" type="password" v-model="confirmPassword" placeholder="" @blur="validateConfirmPassword" />
    </div>
    <p v-if="errors.confirmPassword" class="error">{{ errors.confirmPassword }}</p>
    
    <!-- 验证码输入框 -->
    <div class="input-group">
      <label for="verificationCode">验证码*:</label>
      <input id="verificationCode" type="text" v-model="verificationCode" placeholder="" @blur="validateVerificationCode" />
      <button class="send-verification-code" @click="sendSmsCode" :disabled="isCountingDown" :class="{ 'counting-down': isCountingDown }">
          {{ isCountingDown ? `${countdown}秒后重试` : '获取' }}</button>
    </div>
    <p v-if="errors.verificationCode" class="error">{{ errors.verificationCode }}</p>
    
    <!-- 注册按钮 -->
    <button class="register-button" @click="register">注册</button>
    </div>
 </div>
</template>

<script>
import { Register, sendSmsCode } from '@/services/api'; // 导入注册和发送验证码 API

export default {


  data() {
    return {
      avatar:'@/assets/images/avatar.jpg',
      gender:'',
      id: '',
      nickname: '',
      birthday:'',
      phoneNumber: '',
      email:'',
      password: '',
      confirmPassword: '',
      verificationCode: '',
      errors: {
        avatar:'',
        gender:'',
        id: '',
        nickname: '',
        birthday:'',
        phoneNumber: '',
        email:'',
        password: '',
        confirmPassword: '',
        verificationCode: '',
      },
      Code:'',
      isCountingDown:false,
      countdown:60,
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
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.avatar = e.target.result;
        };
        reader.readAsDataURL(file);
      }
    },

    validateId() {
        if (!this.id) {
          this.errors.id = '账号 不能为空';
        } else if (/\s/.test(this.id)) {
          this.errors.id = '账号 不能包含空格';
        } else if (this.id.length < 6 || this.id.length > 12) {
          this.errors.id = '账号 长度必须在 6 到 12 个字符之间';
        } else if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(this.id)) {
          this.errors.id = '账号 开头只能是字母或下划线，且只能包含字母、数字和下划线';
        } else {
          this.errors.id = '';
        }
    },

    // 验证性别
    validateGender() {
      if (!this.gender) {
        this.errors.gender = '性别不能为空';
      } else {
        this.errors.gender = '';
      }
    },
    
    // 验证昵称
    validateNickname() {
      if (!this.nickname) {
        this.errors.nickname = '昵称不能为空';
      } else {
        this.errors.nickname = '';
      }
    },


    // 验证生日格式
    validateBirthday() {
      if (!this.birthday) {
        this.errors.birthday = '请选择生日！';
      }else {
        this.errors.birthday = '';
      }
    },
    
    // 验证手机号
    validatePhoneNumber() {
      if ((this.phoneNumber)&&!/^1[3-9]\d{9}$/.test(this.phoneNumber)) {
        this.errors.phoneNumber = '手机号格式不正确';
      } else {
        this.errors.phoneNumber = '';
      }
    },
    
    // 验证邮箱
    validateEmail() {
        if (!this.email) {
          this.errors.email = '邮箱不能为空';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.email)) {
          this.errors.email = '邮箱格式不正确';
        } else {
          this.errors.email = '';
        }
    },

    validatePassword() {
        if (!this.password) {
          this.errors.password = '密码不能为空';
        } else if (this.password.length < 6) {
          this.errors.password = '密码长度不能少于6位';
        } else if (/\s/.test(this.password)) {
          this.errors.password = '密码不能包含空格';
        } else if (!/^[a-zA-Z0-9]+$/.test(this.password)) {
          this.errors.password = '密码只能是数字和字母的组合';
        } else {
          this.errors.password = '';
        }
    },
    
    // 验证确认密码
    validateConfirmPassword() {
      if (!this.confirmPassword) {
        this.errors.confirmPassword = '确认密码不能为空';
      } else if (this.confirmPassword !== this.password) {
        this.errors.confirmPassword = '两次输入的密码不一致';
      } else {
        this.errors.confirmPassword = '';
      }
    },
    
    // 验证验证码
    validateVerificationCode() {
      if (!this.verificationCode) {
        this.errors.verificationCode = '验证码不能为空';
      } else {
        this.errors.verificationCode = '';
      }
    },

    async validateCode(){
        if(Code){
          if(Code !== this.verificationCode){
              alert('验证码错误');
              return;
          }
        }
        else{
          alter('请先获取验证码！');
          return;
        }
      },
    
    // 发送验证码
    async sendSmsCode() {
      if (!this.validateEmail()) {
        return;
      }
      
      try {
        const response = await sendSmsCode({
          command:'register',
          email: this.email,
        });
        if (response.success) {
          alert('验证码已发送');
          this.Code = response.code;
          this.startCountdown();//启动一分钟倒计时
        } else {
          alert(response.message || '发送验证码失败');
        }
      } catch (error) {
        alert(error || '发送验证码失败');
      }
    },

    // 启动倒计时
    startCountdown() {
      this.isCountingDown = true;
      this.countdown = 60;

      const timer = setInterval(() => {
        this.countdown--;
        if (this.countdown <= 0) {
          clearInterval(timer);
          this.isCountingDown = false;
        }
      }, 1000);
    },
    
    // 注册方法
    async register() {
      this.validateId();
      this.validateGender();
      this.validateNickname();
      this.validatePhoneNumber();
      this.validateEmail();
      this.validatePassword();
      this.validateConfirmPassword();
      this.validateVerificationCode();
      this.validateBirthday();
      
      if (Object.values(this.errors).some(error => error)) {
        return;
      }

      this.validateCode();
      
      try {
        const response = await Register({
          avatar :this.avatar,
          gender:this.gender,
          id: this.id,
          nickname: this.nickname,
          birthday:this.birthday,
          phone: this.phoneNumber,
          email:this.email,
          password: this.password,
        });
        
        if (response.success) {
          alter(response.message);
          this.$router.push('/loginth');
        } else {
          alert(response.message || '注册失败');
        }
      } catch (error) {
        alert(error || '注册失败');
      }
    },
  },
};
</script>

<style scoped>

.register {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  align-items: center;
  justify-items: center;
  height: 100vh;
  background-color: #f9f9f9;
  padding: 20px;
  box-sizing: border-box;
}

.container {
    display:flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
    border-radius: 8px; /* 添加圆角效果 */
    background-color: #fff; /* 添加背景色 */
    padding: 20px; /* 添加内边距 */
    margin-top: 40px;
    width: 400px;
  }

h2 {
  margin-bottom: 10px;
  align-items: center;
  font-size: 24px;
  color: #333;
}

.avatar-container {
  position: relative;
  cursor: pointer;
}

.headavatar {
  width: 100px;
  height: 100px;
  margin-top: 10px;
  margin-bottom: 20px;
  margin-left: 10px;
  border-radius: 100%;
  cursor: pointer;
}

.radio-group {
  display: flex;
  align-items: center; /* 垂直居中对齐 */
  gap: 10px; /* 设置单选框和文字之间的间距 */
}

.radio-group label {
  display: flex;
  align-items: center; /* 垂直居中对齐 */
  font-size: 14px;
  color: #666;
  margin-right: 5px; /* 设置标签之间的间距 */
}

.radio-group input[type="radio"] {
  margin-left: 5px; /* 设置单选框和文字之间的间距 */
  vertical-align: middle; /* 垂直居中对齐 */
}

.birthday-input-group {
  display: flex;
  align-items: center;
  width: 100%;
}

.date-picker {
  margin-left: 18px;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: pointer;
}

.calendar-icon {
  font-size: 30px;
  cursor: pointer;
  color: #42b983;
}

.input-group {
  display: flex;
  grid-template-columns: 100px 1fr;
  gap: 10px;
  align-items:center;
  width: 100%;
  max-width: 400px;
  margin-bottom: 20px;
}

.input-group label {
  font-size: 14px;
  color: #666;
  text-align: left;
  width: 89px;
}

.input-group input {
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
  width: 100%;
  flex:1;
}

.send-verification-code {
  padding: 10px;
  width: 20%;
  font-size: 14px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-left: 0px;
}

.send-verification-code.counting-down {
  background-color: #ccc; /* 倒计时时的背景色 */
  cursor: not-allowed; /* 倒计时时的不可点击状态 */
}

.error {
  color: red;
  font-size: 12px;
  margin-top: 5px;
  text-align: left;
  width: 100%;
  max-width: 400px;
}

.register-button {
  width: 100%;
  max-width: 100px;
  padding: 10px;
  font-size: 16px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
}

.register-button:hover {
  background-color: #369f6e;
}


</style>