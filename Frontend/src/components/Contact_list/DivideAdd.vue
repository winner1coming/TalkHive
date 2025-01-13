<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <h2>{{type==='add'?"新建分组":"重命名分组"}}</h2>
      <input 
        v-model="newDivide" 
        placeholder="输入分组名称"
        @keyup.enter="addDivide"
      />
      <button @click="addDivide">确认</button>
    </div>
  </div>
</template>

<script>

export default {
  data() {
    return {
      newDivide: '',
      type: '',  // add, rename
    };
  },
  methods: {
    async addDivide() {
			if (this.newDivide.trim()) {
        if(this.type === 'rename'){
          this.$emit('rename-divide', this.newDivide);
        }else{
          this.$emit('add-divide', this.newDivide);
        }
				this.close();
			}else{
        alert('分组名不能为空');
      }
		},
    close() {
      this.newDivide = '';
      this.$emit('close');
    },
  },
};
</script>

<style scoped src="@/assets/css/contactList.css"> </style>
<style scoped>
.modal-content {
  width: 300px;
  height: 130px;
}
input {
	margin: 10px;
	width: 80%;
	padding: 10px;
	border: 1px solid #ddd;
	border-radius: 4px;
}

button {
	border: none;
	border-radius: 4px;
	padding: 10px;
	cursor: pointer;
	margin: 10px;
}
</style>