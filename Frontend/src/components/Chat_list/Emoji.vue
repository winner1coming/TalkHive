<template>
  <div v-show="visible" class="emoji-picker" :style="{ top: `${y}px`, left: `${x}px` }">
    <emoji-picker @emoji-click="emojiClick"></emoji-picker>
  </div>
</template>

<script>
import { EventBus } from '@/components/base/EventBus';
import 'emoji-picker-element';
export default {
  data() {
    return {
      visible: false,
      x: 0,
      y: 0,
    };
  },
  methods:{
    emojiClick(event) {
      EventBus.emit('emoji-click', event.detail.unicode);
      this.hide();
    },
    show(event, boundD, boundR) {  // boundD, boundR 为边界的坐标
      EventBus.emit('float-component-opened', this); // 通知其他组件
      const cardWidth = 343;
      const cardHeight = 400;
      const x = event.clientX + cardWidth > boundR ? event.clientX - cardWidth : event.clientX;
      const y = event.clientY - cardHeight -20 < 0 ? 0 : event.clientY - cardHeight -20;
      this.x = x;
      this.y = y;
      this.visible = true;
      console.log('show');
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide() {
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
  },
  mounted() {
    EventBus.on('other-float-component', (component) => {
      if (this.visible && component !== this) {
        this.hide();
      }
    });
    EventBus.on('close-float-component', (clickedElement) => {
      if (this.visible && !this.$el.contains(clickedElement)) {
        this.hide();
      }
    });
  },
};
</script>

<style scoped>
.emoji-picker {
  position: fixed;
  z-index: 2000;
}
</style>