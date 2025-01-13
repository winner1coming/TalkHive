<template>
  <transition name="fade">
    <div v-if="visible" class="notification" :class="type">
      {{ message }}
    </div>
  </transition>
</template>

<script>
export default {
  data() {
    return {
      visible: false,
      message: '',
      type: 'success' // 'success' or 'error'
    };
  },
  methods: {
    show(message, type = 'success') {
      this.message = message;
      this.type = type;
      this.visible = true;
      setTimeout(() => {
        this.visible = false;
      }, 3000); // 显示3秒后自动隐藏
    }
  }
};
</script>

<style scoped>
.notification {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  border-radius: 5px;
  color: rgb(37, 35, 35);
  z-index: 1000;
  font-size: var(--font-size);
}

.notification.success {
  background-color: #f4f6f4e9;
}

.notification.error {
  background-color: #f6c5c2e0;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s;
}

.fade-enter, .fade-leave-to  {
  opacity: 0;
}
</style>