<template>
	<div class="modal-overlay" @click.self="close">
	  <div class="modal-content">
			<h4>{{ this.type==="in"? "批量移入" : "批量移出" }}</h4>
			<div class="select-all">
				<input 
					type="checkbox" 
					v-model="selectAll" 
					@change="toggleSelectAll"
					ref="selectAllCheckbox"
				>
			</div>
			<ul class="items">
				<li 
					v-for="person in persons" 
					:key="acount_id"
					@click="togglePersonSelection(person)"
				> 
					<input type="checkbox" v-model="selectedPersons" :value="person">
					<img :src="person.avatar" alt="avatar" width="50" height="50" />
					<div class="left">
						<p class="name">{{ person.name }}</p>
						<p class="remark">{{ (`签名：${person.signature}`)}}</p>
					</div>
				</li>
			</ul>
			<button class="confirm-button" @click="confirmSelection">确认</button>
	  </div>
	</div>
</template>
  
<script>

export default {
	props:['type', 'persons'],
	data() {
	  return {
			selectAll: false,
			selectedPersons: [],
	  };
	},
	watch:{
		selectedPersons:{
			handler(val){
				const totalPersons = this.persons.length;
				this.selectAll = val.length === totalPersons;
				if (this.$refs.selectAllCheckbox) {
					this.$refs.selectAllCheckbox.indeterminate = val.length > 0 && val.length < totalPersons;
				}
			},
			deep: true,
			immediate: true,
		}
	},
	methods: {
	  async confirmSelection() {
			this.$emit('confirm', this.selectedPersons);
			this.close();
	  },
	  close() {
			this.selectAll = false;
			this.selectedPersons = [];
			this.$refs.selectAllCheckbox.indeterminate = false;
			this.$emit('close');
	  },
	  togglePersonSelection(person) {
			const index = this.selectedPersons.indexOf(person);
			if (index === -1) {
				this.selectedPersons.push(person);
			} else {
				this.selectedPersons.splice(index, 1);
			}
	  },
	  toggleSelectAll() {
			this.selectedPersons = this.selectAll ? [...this.persons] : [];
	  }
	},
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