<!--用来显示提示消息的弹窗-->
<template>
    <div v-if="visible" class="modal-overlay">
      <div class="modal-content">
        <img src="@/assets/icon/tips.png" alt="Tips" class="icon"/>
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
    background:var(--background-color);
    padding: 20px;
    border-radius: 8px;
    text-align: center;
  }

  .modal-content p{
    color: var(--text-color);
    font-size: var(--font-size);
    gap:20px;
    margin-bottom: 10px;
  }
  
  button {
    padding: 5px 8px;
    background-color: var(--button-background-color1);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    margin-top: 10px;
    cursor: pointer;
  }
  
  button:hover,button.active {
    background-color: var(--button-background-color2);
  }

  .icon{
    width:25px;
    height: 25px;
    margin-top: -5px;
    margin-bottom: 5px;
  }
  </style>
