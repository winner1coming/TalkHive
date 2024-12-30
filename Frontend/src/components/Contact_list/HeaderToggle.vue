<template>
	<div style="width: 100%;">
	  <div class="header-toggle">
			<p class="toggle" @click="toggleContent">
				<span :class="{'arrow-down': showFullContent, 'arrow-right': !showFullContent}">
					<img v-show="!showFullContent" src="@/assets/images/arrow-right.png" class="arrow-right-img"/>
					<img v-show="showFullContent" src="@/assets/images/arrow-down.png" class="arrow-down-img"/>
				</span>
			</p>
			<p class="title">{{ previewText }}</p>
			<button 
				:style="{visibility: visibility}" 
				@click="manageDivide" 
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
	  manageDivide(event) {
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
	height: 100%;
}
.toggle {
	display: flex;
	cursor: pointer;
	margin: 8px;
	align-self: stretch;
	align-items: center;
	height: 100%;
}
.arrow-right::before {
	transition: transform 0.2s ease, opacity 0.2s ease;
	opacity: 1;
}
.arrow-right-img{
	width: 20px;
	height: 20px;
}
.arrow-down::before {
	transition: transform 0.2s ease, opacity 0.2s ease;
	transform: rotate(180deg);
	opacity: 1;
}
.arrow-down-img{
	width: 20px;
	height: 20px;
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