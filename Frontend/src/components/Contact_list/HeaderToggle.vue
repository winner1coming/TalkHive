<template>
	<div style="width: 100%;">
	  <div class="header-toggle">
			<p class="toggle" @click="toggleContent">
				<span :class="{'arrow-down': showFullContent, 'arrow-right': !showFullContent}"></span>
			</p>
			<p class="title">{{ previewText }}</p>
			<button 
				:style="{visibility: visibility}" 
				@click="manageDevide" 
				class="toggle-button"
			>管理</button>
	  </div>
	  <transition-group name="drop" tag="div" v-if="showFullContent">
				<slot></slot>
	  </transition-group>
	</div>
</template>
  
<script>
export default {
	props: {
	  previewText: {
			type: String,
			required: true
	  },
		manageable: {
			type: Boolean,
			default: true,
			required: false
		},
	},
	data() {
	  return {
			showFullContent: false,
			
	  };
	},
	computed: {
	  visibility() {
			return this.manageable ? 'visible' : 'hidden';
	  }
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
	border: 0.5px solid #888282;
	margin-bottom: 3px;
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
	transition: transform 0.2s ease, opacity 0.2s ease;
	opacity: 1;
}
.arrow-down::before {
	content: '▼';
	transition: transform 0.2s ease, opacity 0.2s ease;
	transform: rotate(180deg);
	opacity: 1;
}
.title{
	align-self: center;
}
.toggle-button{
	align-self: flex-end;
}

.drop-enter-active, .drop-leave-active {
  transition: all 0.5s ease;
}
.drop-enter, .drop-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>