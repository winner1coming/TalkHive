<template>
  <div v-if="visible" class="member-select" :style="{ top: `${y}px`, left: `${x}px` }">
    <SearchBar :isImmidiate="false" @search="search" @button-click="search"/>
      <ul class="items" v-show="members.length">
        <li 
          v-for="member in members" 
          :key="member.account_id"
          @click="selectMember(member)"
        >
          <input type="radio" v-model="selectedMember" :value="member">
          <div class="avatar">   <!-- 头像-->
            <img :src="member.avatar" alt="avatar" />
          </div>
          <div class="name">{{ member.remark?member.remark:(member.group_nickname?member.group_nickname:member.nickname) }}</div>
        </li>
      </ul>
  </div>
</template>

<script>
import { EventBus } from '@/components/base/EventBus';
export default {
  props:['members'],
  data() {
    return {
      visible: false,
      x: 0,
      y: 0,
      selectedMember: null,
    };
  },
  methods: {
    selectMember(member) {
      this.selectedMember = member;
      this.$emit('select-member', this.selectedMember);
      this.hide();
    },
    show(event, boundD, boundR) {  // boundD, boundR 为边界的坐标
      EventBus.emit('float-component-opened', this); // 通知其他组件
      const cardWidth = 200;
      const cardHeight = 400;
      const x = event.clientX + cardWidth > boundR ? event.clientX - cardWidth : event.clientX;
      const y = event.clientY + cardHeight > boundD ? boundD - cardHeight : event.clientY;
      this.x = x;
      this.y = y;
      this.visible = true;
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
.member-select {
  position: absolute;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 1px solid #ddd;
  border-radius: 8px;
  width: 200px;
  height: 300px;
  background-color: var(--background-color);
  color: var(--text-color);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.items {
	list-style: none;
	padding: 0;
	margin: 0;
	flex: 9;
	overflow-y: auto;
	border: 1px solid #ddd;
	border-radius: 4px;
  width: 100%;
}
.items li {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #ddd;
	cursor: pointer;
}

.avatar img {
  width: 30px;
  height: 30px;
  border-radius: 50%;
}


.name {
  font-weight: bold;
  font-size: var(--font-size-small);
}

.remark {
  margin-top: 5px;
  font-size: var(--font-size-small);
  color: #666;
}

button {
  margin-top: 10px;
  padding: 5px 10px;
  border: none;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: var(--button-background-color1);
}
</style>