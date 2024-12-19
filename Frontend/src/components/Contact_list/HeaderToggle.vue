<template>
	<div style="width: 100%;">
	  <div class="header-toggle">
			<span class="toggle" @click="toggleContent">
				<i :class="{'arrow-down': showFullContent, 'arrow-right': !showFullContent}"></i>
			</span>
			{{ previewText }}
			<button @click="manageDevide" class="toggle-button">管理</button>
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
	  },
	  manageDevide(event) {
		this.$emit('manage-devide', event);
	  },
	}
  };
</script>
  
<style scoped>
.header-toggle {
	width: 100%;
	display: inline-block;
}
.toggle {
	cursor: pointer;
	margin-left: 5px;
	float: left;
}
.arrow-right::before {
	content: '>';
	display: inline-block;
	transition: transform 0.3s ease;
}
.arrow-down::before {
	content: '▼';
	display: inline-block;
	transition: transform 0.3s ease;
}
.toggle-button{
	float: right;
}
.fade-enter-active, .fade-leave-active {
	transition: opacity 0.3s ease;
}
.fade-enter, .fade-leave-to {
	opacity: 0;
}
</style>