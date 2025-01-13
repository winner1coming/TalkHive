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
					v-for="(divide, index) in divides" 
					:key="index"
					@click="toggleDivideSelection(divide)"
				>
					<input type="checkbox" v-model="selectedDivides" :value="divide">
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
			selectAll: false,
			selectedDivides: [],
			watchSection: null,
	  };
	},
	methods: {
	  async confirmSelection() {
			this.$emit('delete-divides', this.selectedDivides);
			//await deleteFriendGroup(this.selectedDivides);  // selectedDivides是被选择的分组的名称的数组
			this.close();
	  },
	  close() {
			this.selectAll = false;
			this.selectedDivides = [];
			this.$refs.selectAllCheckbox.indeterminate = false;
			this.$emit('close');
	  },
		toggleDivideSelection(divide) {
        const index = this.selectedDivides.indexOf(divide);
        if (index === -1) {
            this.selectedDivides.push(divide);
        } else {
            this.selectedDivides.splice(index, 1);
        }
    },
	  toggleSelectAll() {
			this.selectedDivides = this.selectAll ? [...this.divides] : [];
	  }
	},
	mounted() {
		this.watchSection = this.$watch('selectedDivides', (val) => {
			const totalDivides = this.divides.length;
			this.selectAll = val.length === totalDivides && totalDivides > 0;
			this.$refs.selectAllCheckbox.indeterminate = val.length > 0 && val.length < totalDivides;
		}, { deep: true, immediate: true });
	},
	beforeUnmount() {
		this.watchSection();
	}
};
</script>
  
<style scoped src="@/assets/css/contactList.css"></style>
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

.confirm-button {
	flex: 1;
	margin-top: 10px;
	border: none;
	border-radius: 4px;
	padding: 10px;
	cursor: pointer;
}
</style>