<template>
	<div class="modal-overlay" @click.self="close">
	  <div class="modal-content">
			<h4>{{ this.type==="in"? "移入" : "移出" }}</h4>
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
						<p class="name">{{ person.remark }}</p>
						<p class="remark">{{ (`签名：${person.signature.length > this.maxChars ? person.signature.slice(0, this.maxChars) + '...' : person.signature}`)}}</p>
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
			maxChars: 20,
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
			if(this.type === "in"){
				this.$emit('divide-in', this.selectedPersons);
			}
			else{
				this.$emit('divide-out', this.selectedPersons);
			}
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