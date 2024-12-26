<template>
	<div style="width: 100%;">
	  <div class="header-toggle">
			<p class="toggle" @click="toggleContent">
				<p :class="{'arrow-down': showFullContent, 'arrow-right': !showFullContent}"></p>
			</p>
			<p class="title">{{ previewText }}</p>
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
		this.$emit('manage-divide', event);
	  },
	}
  };
</script>
  
<style scoped src="@/assets/css/contactList.css"></style>
<style scoped>
.header-toggle {
	width: 100%;
	display: flex;
	justify-content: space-between;
	align-items: center;
	border: 0.5px solid #e0e0e0;
}
.toggle {
	display: flex;
	cursor: pointer;
	margin: 5px;
	align-self: start;
	align-items: center;
	height: 100%;
}
.arrow-right::before {
	content: '>';
	transition: transform 0.2s ease;
}
.arrow-down::before {
	content: '▼';
	transition: transform 0.2s ease;
}
.title{
	align-self: center;
}
.toggle-button{
	align-self: flex-end;
}
.fade-enter-active, .fade-leave-active {
	transition: opacity 0.2s ease;
}
.fade-enter, .fade-leave-to {
	opacity: 0;
}
</style>