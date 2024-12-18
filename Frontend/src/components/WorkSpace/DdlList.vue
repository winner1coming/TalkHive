<template>
  <div class="ddl-list">
    <div class="ddl-container">
      <!-- 左侧：待完成的 DDL -->
      <div class="ddl-left" :style="leftDdlStyle">
        <!-- 标题部分 -->
        <div class="header">
          <h2>我的 DDL</h2>
          <button @click="showCreateDdl = true" class="new-btn">+</button>
        </div>

        <!-- 新建 DDL 编辑框 -->
        <div v-if="showCreateDdl" class="create-ddl-modal">
          <div class="modal-content">
            <h3>新建 DDL</h3>
            <label for="deadline-year">截止时间：</label>
            <div class="deadline-inputs">
              <input
                type="number"
                v-model="newDdl.deadline.year"
                placeholder="年"
                min="1900"
                max="2100"
              />
              <input
                type="number"
                v-model="newDdl.deadline.month"
                placeholder="月"
                min="1"
                max="12"
              />
              <input
                type="number"
                v-model="newDdl.deadline.day"
                placeholder="日"
                min="1"
                max="31"
              />
            </div>

            <label for="task-content">任务内容：</label>
            <textarea
              v-model="newDdl.task_content"
              placeholder="输入任务内容"
              rows="3"
              style="width: 100%;"
            ></textarea>

            <label for="important">是否设为重要：</label>
            <input
              type="checkbox"
              v-model="newDdl.important"
            />

            <div class="modal-actions">
              <button @click="saveDdl" class="save-btn">保存</button>
              <button @click="cancelCreate" class="cancel-btn">取消</button>
            </div>
          </div>
        </div>

        <ul>
          <li v-for="item in ddlList" :key="item.task_id" class="ddl-item">
            <input 
              type="checkbox" 
              class="ddl-checkbox" 
              v-model="item.completed" 
              @change="updateDdlStatus(item)" 
            />
            <span class="deadline">{{ item.deadline }}</span>
            <span class="task-content">{{ item.task_content }}</span>
            <button @click="deleteDdl(item)" class="delete-btn">删除</button>
          </li>
        </ul>
      </div>

      <!-- 右侧：已完成 DDL -->
      <div class="ddl-right" v-if="showCompleted">
        <h2>已完成</h2>
        <ul>
          <li v-for="item in completedDdl" :key="item.task_id" class="ddl-item completed">
            <span class="deadline">{{ item.deadline }}</span>
            <span class="task-content">{{ item.task_content }}</span>
          </li>
        </ul>
      </div>
    </div>

    <!-- "已完成" 或 "收起" 按钮 -->
    <div class="toggle-btn-container">
      <button @click="toggleCompletedDdl" class="toggle-btn">
        {{ showCompleted ? '> 收起' : '< 已完成' }}
      </button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'DdlList',
  data() {
    return {
      ddlList: [], // 待完成的 DDL
      completedDdl: [], // 已完成的 DDL
      showCompleted: false, // 控制是否显示已完成的 DDL
      showCreateDdl: false, // 控制是否显示新建 DDL 编辑框
      newDdl: {
        deadline: {
          year: '2024',
          month: '12',
          day: '15',
        },
        task_content: '',
        important: false,
      },
      userId: 'your-user-id', // 用户ID
    };
  },
  created() {
    this.fetchDdlList();
    this.fetchCompletedDdl();
  },
  methods: {
    // 获取待完成的 DDL 列表
    async fetchDdlList() {
      try {
        const response = await axios.get('/workspace/ddl', {
          params: { id: this.userId },
        });
        if (response.data.status === 200) {
          this.ddlList = response.data.data;
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法获取 DDL 列表:', error);
        alert('获取 DDL 列表失败！');
      }
    },

    // 获取已完成的 DDL 列表
    async fetchCompletedDdl() {
      try {
        const response = await axios.get('/workspace/ddl/completed', {
          params: { id: this.userId },
        });
        if (response.data.status === 200) {
          this.completedDdl = response.data.data;
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法获取已完成的 DDL 列表:', error);
        alert('获取已完成的 DDL 列表失败！');
      }
    },

    // 切换显示已完成的 DDL
    toggleCompletedDdl() {
      this.showCompleted = !this.showCompleted;
    },

    async saveDdl() {
      try {
        const response = await axios.post('/workspace/ddl/create', {
          id: this.userId,
          deadline: `${this.newDdl.deadline.year}-${this.newDdl.deadline.month}-${this.newDdl.deadline.day}`,
          task_content: this.newDdl.task_content,
          important: this.newDdl.important,
        });

        if (response.data.status === 200) {
          alert('DDL 创建成功');
          this.showCreateDdl = false; // 关闭编辑框
          this.fetchDdlList(); // 刷新待完成 DDL 列表
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法创建 DDL:', error);
        alert('创建 DDL 失败！');
      }
    },

    cancelCreate() {
      this.showCreateDdl = false; // 取消新建 DDL
    },

    // 更新 DDL 状态为已完成
    async updateDdlStatus(item) {
      try {
        const response = await axios.post('/workspace/ddl/update', {
          id: this.userId,
          task_id: item.task_id,
          completed: item.completed, // 如果复选框选中，则为 true
        });
        if (response.data.status === 200) {
          alert('状态更新成功');
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法更新 DDL 状态:', error);
        alert('更新 DDL 状态失败！');
      }
    },

    // 删除某条 DDL
    async deleteDdl(item) {
      try {
        const response = await axios.delete('/workspace/ddl/delete', {
          data: { id: this.userId, task_id: item.task_id },
        });
        if (response.data.status === 200) {
          alert('删除成功');
          // 刷新待完成和已完成的 DDL 列表
          this.fetchDdlList();
          this.fetchCompletedDdl();
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法删除 DDL:', error);
        alert('删除 DDL 失败！');
      }
    },
  },
  computed: {
    // 根据 showCompleted 控制左侧待完成的 DDL 宽度
    leftDdlStyle() {
      return {
        width: this.showCompleted ? '60%' : '100%', // 显示已完成时为48%，否则占满整个页面
      };
    },
  },
};
</script>

<style scoped>
.ddl-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: center;
  align-items: center;
}

.new-btn {
  background-color: #c7d7e9;
  color: rgb(75, 103, 216);
  border: none;
  cursor: pointer;
  margin-left: 20px;
  margin-bottom: 15px;
  border-radius: 50%;
  width: 45px; /* 宽度 */
  font-size: 36px;
}

.new-btn:hover {
  background-color: #0056b3;
  color: rgb(134, 154, 233);
}

.create-ddl-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

.deadline-inputs input {
  margin-right: 5px;
}

.modal-actions {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
}

.save-btn, .cancel-btn {
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

.save-btn:hover, .cancel-btn:hover {
  background-color: #0056b3;
}

.ddl-container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.ddl-left, .ddl-right {
  width: 48%;
}

.ddl-list h2 {
  margin-bottom: 20px;
}

.ddl-list ul {
  list-style: none;
  padding: 0;
}

.ddl-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #ddd;
}

.ddl-item .ddl-checkbox {
  /* 放大复选框 */
  transform: scale(2); /* 使复选框放大2倍 */
  margin-right: 20px; /* 增加复选框与任务内容之间的间距 */
}

.ddl-item .deadline {
  font-weight: bold;
  margin-right: 10px;
}

.ddl-item .task-content {
  color: #666;
  flex-grow: 1;
}

.ddl-item .delete-btn {
  background-color: #007bff;
  color: rgb(255, 255, 255);
  border: none;
  padding: 5px 10px;
  cursor: pointer;
}

.ddl-item .delete-btn:hover {
  background-color:rgb(75, 103, 216);
}

.completed {
  background-color: #f0f8ff;
  color: #888;
}

.toggle-btn-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.toggle-btn {
  padding: 10px;
  background-color: #9fbddf;
  color: white;
  border: none;
  cursor: pointer;
}

.toggle-btn:hover {
  background-color: #0056b3;
}
</style>
