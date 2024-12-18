<template>
	<div class="modal-overlay" @click.self="close">
	  <div class="modal-content">
			<h4>删除分组</h4>
			<div class="select-all">
				<input 
					type="checkbox" 
					v-model="selectAll" 
					@change="toggleSelectAll"
					ref="selectAllCheckbox"
				>
			</div>
			<ul class="items">
				<!-- 每个分组 -->
				<li 
					v-for="(devide, index) in devides" 
					:key="index"
					@click="toggleDevideSelection(devide)"
				>
					<input type="checkbox" v-model="selectedDevides" :value="devide">
					<span>{{ devide }}</span>
				</li>
			</ul>
			<button class="confirm-button" @click="confirmSelection">确认</button>
	  </div>
	</div>
</template>
  
<script>

export default {
	props:['devides'],
	data() {
	  return {
			selectAll: false,
			selectedDevides: [],
			watchSection: null,
	  };
	},
	methods: {
	  async confirmSelection() {
			//await deleteFriendGroup(this.selectedDevides);  // selectedDevides是被选择的分组的名称的数组
			this.close();
	  },
	  close() {
			this.selectAll = false;
			this.selectedDevides = [];
			this.$refs.selectAllCheckbox.indeterminate = false;
			this.$emit('close');
	  },
		toggleDevideSelection(devide) {
        const index = this.selectedDevides.indexOf(devide);
        if (index === -1) {
            this.selectedDevides.push(devide);
        } else {
            this.selectedDevides.splice(index, 1);
        }
    },
	  toggleSelectAll() {
			this.selectedDevides = this.selectAll ? [...this.devides] : [];
	  }
	},
	mounted() {
		this.watchSection = this.$watch('selectedDevides', (val) => {
			const totalDevides = this.devides.length;
			this.selectAll = val.length === totalDevides;
			this.$refs.selectAllCheckbox.indeterminate = val.length > 0 && val.length < totalDevides;
		}, { deep: true, immediate: true });
	},
	beforeUnmount() {
		this.watchSection();
	}
};
</script>
  
<style scoped src="@/assets/css/contactList.css"></style>
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