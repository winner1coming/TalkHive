<!--用来显示提示消息的弹窗-->
<template>
    <div v-if="visible" class="modal-overlay">
      <div class="modal-content">
        <p>{{ message }}</p>
        <button @click="closeModal">确定</button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      visible: {
        type: Boolean,
        default: false,
      },
      message: {
        type: String,
        default: '',
      },
    },
    methods: {
      closeModal() {
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
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .modal-content {
    background: rgb(252, 251, 251);
    padding: 20px;
    border-radius: 8px;
    text-align: center;
  }

  .modal-content p{
    color: black;
    gap:10px;
  }
  
  button {
    padding: 10px 20px;
    background-color: #42b983;
    color: rgb(24, 22, 22);
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #369f6d;
  }

  
  </style>