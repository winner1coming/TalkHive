<template>
  <div class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <h2>新建分组</h2>
      <input v-model="newDevide" placeholder="输入分组名称">
      <button @click="addDevide">确认</button>
    </div>
  </div>
</template>

<script>
import { addFriendGroup, searchFriendGroup } from '@/services/api';
export default {
  data() {
    return {
      newDevide: '',
    };
  },
  methods: {
    async addDevide() {
			if (this.newDevide.trim()) {
				// todo 增加分组的api
				this.close();
			}
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