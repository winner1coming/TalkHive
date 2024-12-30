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
  height: 130px;
	display: flex;
	flex-direction: column;
	align-items: center;
}
input {
	margin: 10px;
	width: 80%;
	padding: 10px;
	border: 1px solid #ddd;
	border-radius: 4px;
}

button {
	background-color: #007bff;
	color: white;
	border: none;
	border-radius: 4px;
	padding: 10px;
	cursor: pointer;
	margin: 10px;
}

button:hover {
  background-color: #0056b3;
}
</style>