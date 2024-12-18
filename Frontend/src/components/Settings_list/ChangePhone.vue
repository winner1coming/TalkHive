<template>
    <div class="change-phone">
      <h3>更改手机号</h3>
      <div class="text_group">
        <label for="oldPhone">原号码:</label>
        <input type="text" v-model="oldPhone" :placeholder="oldPhone" />
      </div>
      <div class="text_group">
        <label for="newPhone">新号码:</label>
        <input type="text" v-model="newPhone" :placeholder="newPhone" />
      </div>
      <div class="verificate">
        <label for="intput_v">验证码:</label>
        <input type="text" v-model="intput_v" />
        <button  @click="sendSmsCode">获取</button>
      </div>
    <p>注：验证发将发送到新号码中，请注意查收</p>
    <div class="button_container">
        <button @click="savePhone">保存</button>
        <button @click="cancle">取消</button>
    </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        oldPhone:'',
        newPhone: '',
      };
    },
    methods: {
      savePhone() {
        this.$emit('updateUser', { phone: this.newPhone });
      },
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