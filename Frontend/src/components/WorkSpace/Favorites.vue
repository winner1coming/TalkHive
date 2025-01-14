<template>
  <div class="favorites">
    <h2 class="header">
      <div class="left">
        <img src="@/assets/icon/favorite.png" alt="收藏图标" class="icon"/>
        我的收藏
      </div>
      <!-- 更多按钮 -->
      <img src="@/assets/icon/more.png" alt="更多" class="more-icon" @click="toggleDropdown"/>
      <!-- 下拉框 -->
      <div v-if="showDropdown" class="dropdown">
        <ul>
          <li @click="activateDeleteMode">删除</li>
        </ul>
      </div>
    </h2>

    <!-- 删除模式开启后，显示垃圾桶按钮 -->
    <!-- 垃圾桶图标和取消按钮 -->
    <div v-if="deleteMode" class="trash-actions">
      <img src="@/assets/icon/delete.png" alt="垃圾图标" class="trash-icon" @click="deleteSelected"/>
      <button class="cancel-delete-btn" @click="cancelDelete">取消</button>
    </div>
    <ul>
      <li v-for="item in favorites" :key="item.message_id" class="favorite-item">
        <!-- 复选框 -->
        <input v-if="deleteMode" type="checkbox" class="favorite-checkbox" v-model="selectedItems" :value="{ message_id: item.message_id, type: item.type }" />
        <span class="type" style="color:darkgrey">{{ getTypeLabel(item.type) }}</span>
        <span class="object-name"  @click="viewItem(item)">{{ item.object_name }}</span>
        <span class="time"> - {{ item.time }}</span>
        <span class="sender"> - {{ item.sender_name }}</span>
      </li>
    </ul>

    <!-- 删除确认框 -->
    <div v-if="showConfirmDelete" class="confirm-delete">
      <p>是否确认删除所选的收藏内容？</p>
      <button @click="confirmDelete">确认</button>
      <button @click="cancelDelete">取消</button>
    </div>
  </div>
</template>

<script>
import * as WorkSpaceAPI from '@/services/workspace_api';

export default {
  name: 'Favorites',
  data() {
    return {
      favorites: [],
      userId: 'your-user-id',
      showDropdown: false,
      deleteMode: false,
      selectedItems: [], // 存储被勾选的收藏项ID
      showConfirmDelete: false,
    };
  },
  created() {
    this.fetchFavorites();
  },
  methods: {
    //获取收藏列表
    async fetchFavorites() {
      try {
        const response = await WorkSpaceAPI.getFavorites();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          this.favorites = response.data;
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error("无法获取收藏列表:", error);
        alert("获取收藏列表失败！");
      }
    },
    getTypeLabel(type) {
      switch (type) {
        case 'message': return '消息';
        case 'code': return '代码';
        case 'note': return '笔记';
        default: return '未知';
      }
    },
    viewItem(item) {
      if (item.type === 'message') {
        this.$router.push({
          path: `/workspace/favorites/${item.message_id}`,
          query: { table: item.message_list_name },
        });
      } else if (item.type === 'code' || item.type === 'note') {
        this.$router.push({
          path: `/workspace/${item.type}s/${item.message_id}`,
        });
      }
    },
    toggleDropdown() {
      this.showDropdown = !this.showDropdown;
    },
    activateDeleteMode() {
      this.deleteMode = true;
      this.showDropdown = false; // 关闭下拉框
    },
    deleteSelected() {
      if (this.selectedItems.length > 0) {
        this.showConfirmDelete = true; // 显示确认删除框
      } else {
        alert("请先选择要删除的收藏项！");
      }
    },
    async confirmDelete() {
      // 执行删除操作
      const response = await WorkSpaceAPI.deleteFavorites(this.selectedItems);
      console.log(response);
      if (response.status === 200) {
        this.selectedItems = []; // 清空选择
        this.deleteMode = false; // 退出删除模式
        this.showConfirmDelete = false; // 隐藏确认框
        this.fetchFavorites();// 刷新
      } else {
        alert(response.data.message);
      }
      // this.favorites = this.favorites.filter(item => !this.selectedItems.includes(item.message_id));
    },
    cancelDelete() {
      this.selectedItems = []; // 清空选择
      this.deleteMode = false; // 退出删除模式
      this.showConfirmDelete = false; // 隐藏确认框
    }
  },
};
</script>

<style scoped>
.favorites {
  padding: 20px;
}

.favorites h2 {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* .more-btn {
  padding: 5px 10px;
  background-color:none;
  color: rgb(118, 118, 118);
  cursor: pointer;
} */

.more-btn:hover {
  background-color: #c6c6c6;
}

.dropdown {
  position: absolute;
  top: 85px;
  right: 50px;
  border: 1px solid #ddd;
  background-color: white;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
  z-index: 100;
  font-size: 20px;
}

.dropdown ul {
  list-style: none;
  margin: 0;
  padding: 10px;
}

.dropdown li {
  padding: 8px 12px;
  cursor: pointer;
}

.dropdown li:hover {
  background-color: #f0f0f0;
}

.trash-actions {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.cancel-delete-btn {
  margin-left: 10px;
  padding: 5px 10px;
  background-color: #eff0f0;
  color: rgb(91, 91, 91);
  cursor: pointer;
}

.cancel-delete-btn:hover {
  background-color: #bebfc0;
}

.favorite-checkbox {
  /* 放大复选框 */
  transform: scale(2); /* 使复选框放大2倍 */
  margin-right: 20px; /* 增加复选框与任务内容之间的间距 */
}

.favorite-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #ddd;
  justify-content: space-between;
}

.favorite-item .type,
.favorite-item .object-name {
  font-weight: bold;
  margin-right: 10px;
}

.favorite-item .sender {
  color: #666;
}

/* .favorite-item button {
  margin-left: auto;
  padding: 5px 10px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

.favorite-item button:hover {
  background-color: #0056b3;
} */

.confirm-delete {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 20px;
  background-color: white;
  border: 1px solid #ddd;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  z-index: 200;
}

.confirm-delete button {
  margin: 10px;
}
.icon{
  width: 50px;
  height: 50px;
  margin-right: 5px;
}
.trash-icon {
  width: 40px;
  height: 40px;
  background-color: none;
  color: white;
}

.more-icon {
  width: 30px;
  height: 30px;
  margin-right: 5px;
}

.more-icon:hover {
  cursor: pointer;
  /* background-color: #dacfdb; */
}
.trash-icon:hover {
  cursor: pointer;
  /* background-color: #dacfdb; */
}
.header {
  display: flex;
  justify-content: space-between; /* 将内容分布到两端 */
  align-items: center; /* 垂直居中对齐 */
}

.left {
  display: flex;
  align-items: center; /* 垂直居中图标和文字 */
  justify-content: center; /* 水平居中 */
  flex: 1; /* 让左侧部分占据剩余空间 */
}

</style>
