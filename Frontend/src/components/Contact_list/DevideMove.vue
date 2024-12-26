<template>
	<div class="modal-overlay" @click.self="close">
	  <div class="modal-content">
			<h4>移动至</h4>
			<ul class="items">
				<!-- 每个分组 -->
				<li 
					v-for="(divide, index) in divides" 
					:key="index"
					@click="selectedDevide = divide"
				>
					<input type="radio" v-model="selectedDevide" :value="divide">
					<span>{{ divide }}</span>
				</li>
			</ul>
			<button class="confirm-button" @click="confirmSelection">确认</button>
	  </div>
	</div>
</template>
  
<script>

export default {
	props:['divides'],
	data() {
	  return {
			//divides:['家人', '好友', '同事','a','b','c','d','e','f', 'g'], 
			selectAll: false,
			selectedDevide: null,
			multiple: false,
	  };
	},
	methods: {
	  async confirmSelection() {
			if(this.multiple===false) this.$emit('divide-move', this.selectedDevide);
			else this.$emit('divides-move', this.selectedDevide);
			this.close();
	  },
	  close() {
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
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 2000; /* 确保在最上层 */
}

.modal-content {
	background-color: #fff;
	padding: 20px;
	border-radius: 8px;
	width: 300px;
	height: 400px;
	display: flex;
	flex-direction: column;
}

.items {
	list-style: none;
	padding: 0;
	margin: 0;
	flex: 9;
	overflow-y: auto;
	border: 1px solid #ddd;
	border-radius: 4px;
}

.items li {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #ddd;
	cursor: pointer;
}

.select-all {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #ddd;
}

.items li input[type="checkbox"] {
	margin-right: 10px;
}

.confirm-button {
	flex: 1;
	margin-top: 10px;
	background-color: #007bff;
	color: white;
	border: none;
	border-radius: 4px;
	padding: 10px;
	cursor: pointer;
}

.confirm-button:hover {
	background-color: #0056b3;
}
</style>