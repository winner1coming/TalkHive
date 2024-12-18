<template>
	<div v-if="visible" class="context-menu" :style="{ top: `${y}px`, left: `${x}px` }">
		<ul>
			<li v-for="item in items" :key="item" @click="handleClick(item)">
				{{ item }}
			</li>
		</ul>
	</div>
</template>

<script>
import { EventBus } from '@/components/base/EventBus';
export default {
	data() {
		return {
			visible: false,
			x: 0,
			y: 0,
			items: [],
			obj: null,  // 当前触发菜单的对象
		};
	},
	methods: {
		show(event, items, obj, boundD, boundR, w=null, h=null) { // boundD, boundR 为边界的坐标
			if(boundD === null){
				this.x = event.clientX;
				this.y = event.clientY;
			}else{
				const componentWidth = w ? w : 60;
				const componentHeight = h ? h : 400;
				const x = event.clientX + componentWidth > boundR ? event.clientX - componentWidth : event.clientX;
				const y = event.clientY + componentHeight > boundD ? boundD - componentHeight : event.clientY;
				this.x = x;
				this.y = y;
			}
			this.items = items;
			this.visible = true;
			this.obj = obj;
			EventBus.emit('float-component-open', this); // 通知其他组件
		},
		hide() {
			this.visible = false;
			EventBus.emit('hide-float-component'); // 通知其他组件
		},
		// 选中某个选项（item是选中的选项的string，obj是指是谁触发的菜单）
		handleClick(item) {
			this.$emit('select-item', item, this.obj);
			this.hide();
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
.context-menu {
	position: absolute;
	background-color: white;
	border: 1px solid #ccc;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
	z-index: 1000;
}
.context-menu ul {
	list-style: none;
	margin: 0;
	padding: 0;
}
.context-menu li {
	padding: 8px 12px;
	cursor: pointer;
}
.context-menu li:hover {
	background-color: #f5f5f5;
}
</style>