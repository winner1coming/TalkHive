<template>
	<div class="modal-overlay" @click.self="close">
	  <div class="modal-content">
			<h4>移动至</h4>
			<ul class="items">
				<!-- 每个分组 -->
				<li 
					v-for="(divide, index) in divides" 
					:key="index"
					@click="selectedDivide = divide"
				>
					<input type="radio" v-model="selectedDivide" :value="divide">
					<span class="radiomark"></span>
					{{ divide }}
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
			selectedDivide: null,
			multiple: false,
	  };
	},
	methods: {
	  async confirmSelection() {
			if(this.multiple===false) this.$emit('divide-move', this.selectedDivide);
			else this.$emit('divides-move', this.selectedDivide);
			this.close();
	  },
	  close() {
			this.$emit('close');
	  },
	},
};
</script>
  
<style src="@/assets/css/contactList.css"></style>
<style scoped>
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
.items li input[type="checkbox"]:checked ~ .radiomark{
	color: var(--button-background-color);
}

.confirm-button {
	flex: 1;
	margin-top: 10px;
	border: none;
	border-radius: 4px;
	padding: 10px;
	cursor: pointer;
}
</style>