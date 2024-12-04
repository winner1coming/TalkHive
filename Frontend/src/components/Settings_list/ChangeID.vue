<template>
    <div class="changeid">
      <h3>更改ID</h3>
      <div class="hole_page">
        <input type="text" v-model="ID" :placeholder="ID" />
        <p>注：每个用户一年内只能修改一次ID</p>
        <div class="button_container">
            <button @click="saveID">保存</button>
            <button @click="cancle">取消</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        ID: '',
      };
    },
    methods: {
      saveID() {
        this.$emit('updateUser', { id: this.ID });
      },
      cancle(){
        this.$emit('cancle');
      }
    },
  };
  </script>
  
  <style scoped>
  .changeid {
    padding: 20px;
  }

  h3{
    margin-top: 20px;
  }

  .hole_page{
    align-items: center;
    display: flex;
    justify-items: center;
    flex-direction: column;
    margin-top: 20px;
  }

  .hole_page input{
    border: none;
    border-bottom: 1px solid #ccc;
    padding: 10px 0;
    font-size: 16px;
    outline: none;
    text-align: center;
  }

  .button_container{
    display: flex;
    justify-content: space-between;
    width: 100%;
    max-width: 150px;
    margin-top: 50px;
    margin-bottom: 20px;
  }

  p{
    margin-top: 10px;
    font-size: 14px;
    color: #ca3b3b;
  }
  </style>