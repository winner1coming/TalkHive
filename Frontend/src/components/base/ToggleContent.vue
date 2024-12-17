<template>
  <div>
    <div class="content-preview">
      {{ previewText }}
      <span class="toggle" @click="toggleContent">
        <i :class="{'arrow-down': showFullContent, 'arrow-right': !showFullContent}"></i>
      </span>
    </div>
    <transition name="fade">
      <div v-if="showFullContent" class="full-content">
        <slot></slot>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  props: {
    previewText: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      showFullContent: false
    };
  },
  methods: {
    toggleContent() {
      this.showFullContent = !this.showFullContent;
    }
  }
};
</script>

<style scoped>
.toggle {
  cursor: pointer;
  margin-left: 5px;
}
.arrow-right::before {
  content: '<';
  display: inline-block;
  transition: transform 0.3s ease;
}
.arrow-down::before {
  content: '▼';
  display: inline-block;
  transition: transform 0.3s ease;
}
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>