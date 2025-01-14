<template>
  <div class="ddl-list">
    <div class="ddl-container">
      <!-- 左侧：待完成的 DDL -->
      <div class="ddl-left" :style="leftDdlStyle">
        <!-- 标题部分 -->
        <div class="header">
          <h2>
            <div class="ddl_header">
              <img src="@/assets/icon/date-icon.png" alt="ddl图标" class="icon"/>
              我的 DDL
              <img
                src="@/assets/icon/create_note.png"
                alt="添加ddl"
                class="create_note_icon"
                @click="showCreateDdl = true"
              />
            </div>
          </h2>
          <!-- <button @click="showCreateDdl = true" class="new-btn">+</button> -->
        </div>

        <!-- 新建 DDL 编辑框 -->
        <div v-if="showCreateDdl" class="create-ddl-modal">
          <div class="modal-content">
            <h3>新建 DDL</h3>
            <div class="deadline-inputs">
              <input
                type="number"
                v-model="newDdl.deadline.year"
                placeholder="年"
                min="1900"
                max="2100"
              />
              <label for="input_year">年</label>
              <input
                type="number"
                v-model="newDdl.deadline.month"
                placeholder="月"
                min="1"
                max="12"
              />
              <label for="input_month">月</label>
              <input
                type="number"
                v-model="newDdl.deadline.day"
                placeholder="日"
                min="1"
                max="31"
              />
              <label for="input_day">日</label>
              <input
                type="number"
                v-model="newDdl.deadline.hour"
                placeholder="小时"
                min="0"
                max="23"
              />
              <label for="input_hour">：</label>
              <input
                type="number"
                v-model="newDdl.deadline.minute"
                placeholder="分钟"
                min="0"
                max="59"
              />
            </div>

            <label for="task-content">任务内容：</label>
            <textarea
              v-model="newDdl.task_content"
              placeholder="输入任务内容"
              rows="3"
              style="width: 100%;"
            ></textarea>

            <div style="display: flex; align-items: center; justify-content: center;">
              <label for="important" style="margin-right: 10px;">是否设为重要：</label>
              <input
                type="checkbox"
                style="transform: scale(1.6); padding: 2px 5px;"
                v-model="newDdl.important"
              />
            </div>

            <div class="modal-actions">
              <button @click="saveDdl" class="save-btn">保存</button>
              <button @click="cancelCreate" class="cancel-btn">取消</button>
            </div>
          </div>
        </div>

        <!-- 编辑 DDL 弹框 -->
        <div v-if="showEditDdl" class="edit-ddl-modal">
          <div class="modal-content">
            <h3>编辑 DDL</h3>
            <label for="deadline-year">截止时间：</label>
            <div class="deadline-inputs">
              <input
                type="number"
                v-model="editingDdl.deadline.year"
                placeholder="年"
                min="1900"
                max="2100"
              />
              <label for="input_year">年</label>
              <input
                type="number"
                v-model="editingDdl.deadline.month"
                placeholder="月"
                min="1"
                max="12"
              />
              <label for="input_month">月</label>
              <input
                type="number"
                v-model="editingDdl.deadline.day"
                placeholder="日"
                min="1"
                max="31"
              />
              <label for="input_day">日</label>
              <input
                type="number"
                v-model="editingDdl.deadline.hour"
                placeholder="小时"
                min="0"
                max="23"
              />
              <label for="input_hour">：</label>
              <input
                type="number"
                v-model="editingDdl.deadline.minute"
                placeholder="分钟"
                min="0"
                max="59"
              />
            </div>

            <label for="task-content">任务内容：</label>
            <textarea
              v-model="editingDdl.task_content"
              placeholder="输入任务内容"
              rows="3"
              style="width: 100%;"
            ></textarea>

            <div style="display: flex; align-items: center; justify-content: center;">
              <label for="important" style="margin-right: 10px;">是否设为重要：</label>
              <input
                type="checkbox"
                style="transform: scale(1.6); padding: 2px 5px;"
                v-model="editingDdl.important"
              />
            </div>

            <div class="modal-actions">
              <button @click="saveEditDdl" class="save-btn">保存</button>
              <button @click="cancelEdit" class="cancel-btn">取消</button>
            </div>
          </div>
        </div>

        <!-- 删除确认框 -->
        <div v-if="showDeleteConfirm" class="confirm-modal">
          <div class="confirm-content">
            <p>是否确认删除此DDL？</p>
            <div class="modal-actions">
              <button @click="deleteDdl" class="confirm-btn">确认</button>
              <button @click="cancelDelete" class="cancel-btn">取消</button>
            </div>
          </div>
        </div>

        <!-- 更新为已完成确认框 -->
        <div v-if="showUpdateConfirm" class="confirm-modal">
          <div class="confirm-content">
            <p>是否确认已完成此DDL？</p>
            <div class="modal-actions">
              <button @click="updateDdlStatus" class="confirm-btn">确认</button>
              <button @click="cancelUpdate" class="cancel-btn">取消</button>
            </div>
          </div>
        </div>


        <ul>
          <li v-for="item in ddlList" :key="item.task_id" class="ddl-item">
            <input 
              type="checkbox" 
              class="ddl-checkbox" 
              v-model="item.completed" 
              @change="confirmUpdate(item)" 
            />
            <span class="deadline" @click="editDdl(item.task_id)">{{ formatDeadline(item.deadline) }}</span>
            <img v-if="item.important" src="@/assets/icon/important.png" title="重要" class="important-icon"/>
            <!-- <span v-if="item.important" class="important-label">重要</span> -->
            <span v-if="!item.important" class="invisible_important-label">  </span>
            <span class="task-content" @click="editDdl(item)">{{ item.task_content }}</span>
            <img src="@/assets/icon/recycle_delete.png" alt="垃圾图标" class="trash-icon" @click="confirmDelete(item)"/>
          </li>
        </ul>
      </div>

      <!-- 右侧：已完成 DDL -->
      <div class="ddl-right" v-if="showCompleted">
        <h2>已完成</h2>
        <ul>
          <li v-for="item in completedDdl" :key="item.task_id" class="ddl-item completed">
            <span class="deadline">{{ formatDeadline(item.deadline) }}</span>
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
import * as WorkSpaceAPI from '@/services/workspace_api';

export default {
  name: 'DdlList',
  data() {
    return {
      ddlList: [], // 待完成的 DDL
      completedDdl: [], // 已完成的 DDL
      showEditDdl: false, // 是否显示编辑框
      editingDdl: { // 当前编辑的 DDL
        task_id: '',
        deadline: {
          year: '',
          month: '',
          day: '',
        },
        task_content: '',
        important: false,
      },
      showCompleted: false, // 控制是否显示已完成的 DDL
      showCreateDdl: false, // 控制是否显示新建 DDL 编辑框
      showDeleteConfirm: false,
      showUpdateConfirm: false,
      itemToDelete: null,
      itemToUpdate: null,
      newDdl: {
        deadline: {
          year: new Date().getFullYear(),
          month: new Date().getMonth() + 1,
          day: new Date().getDate(),
          hour: new Date().getHours(),
          minute: new Date().getMinutes(),
        },
        task_content: '',
        important: false,
      },
      userId: 'your-user-id', // 用户ID
      importantInt:0,
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
        const response = await WorkSpaceAPI.getDdlList();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          this.ddlList = response.data.map(item => {
            item.deadline = new Date(item.deadline); // 转换为时间类型
            return item;
          });
          this.ddlList.sort((a, b) => a.deadline - b.deadline); // 按时间排序
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
        const response = await WorkSpaceAPI.getCompletedDdlList();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          this.completedDdl = response.data.map(item => {
            item.deadline = new Date(item.deadline); // 转换为时间类型
            return item;
          });
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法获取已完成的 DDL 列表:', error);
        alert('获取已完成的 DDL 列表失败！');
      }
    },

    // 格式化截止时间
    formatDeadline(deadline) {
      return deadline.toISOString().slice(0, 16).replace('T', ' '); // 格式化为 "YYYY-MM-DD HH:mm"
    },

    // 切换显示已完成的 DDL
    toggleCompletedDdl() {
      this.showCompleted = !this.showCompleted;
    },

    // 验证日期是否合法且不小于当前日期
    validateDeadline() {
      const inputDeadline = new Date(
        `${this.newDdl.deadline.year}-${String(this.newDdl.deadline.month).padStart(2, '0')}-${String(this.newDdl.deadline.day).padStart(2, '0')} ${String(this.newDdl.deadline.hour).padStart(2, '0')}:${String(this.newDdl.deadline.minute).padStart(2, '0')}`
      );
      const currentDate = new Date();
      console.log("currentDate：",currentDate);
      console.log("inputDate：",inputDeadline);
      if (inputDeadline < currentDate) {
        this.errorMessage = '截止日期不能小于当前日期和时间。';
        return false;
      } else {
        this.errorMessage = ''; // 清空错误信息
        return true;
      }
    },

    isImportantTrue(important){
      if(important){
        console.log("important:",1);
        return 1;
      }
      console.log("important:",0);
      return 0;
    },

    // 编辑 DDL
    editDdl(item) {
      this.editingDdl = { ...item }; // 复制当前 DDL 数据到 editingDdl
      this.editingDdl.deadline.year = item.deadline.toISOString().slice(0, 16).replace('T', ' ').substring(0,4);
      this.editingDdl.deadline.month = item.deadline.toISOString().slice(0, 16).replace('T', ' ').substring(5,7);
      this.editingDdl.deadline.day = item.deadline.toISOString().slice(0, 16).replace('T', ' ').substring(8,10);
      this.editingDdl.deadline.hour = item.deadline.toISOString().slice(0, 16).replace('T', ' ').substring(11,13);
      this.editingDdl.deadline.minute = item.deadline.toISOString().slice(0, 16).replace('T', ' ').substring(14,16);
      console.log(this.editingDdl);
      this.showEditDdl = true; // 显示编辑框
    },

    // 保存修改后的 DDL
    async saveEditDdl() {
      try {
        // const response = await axios.post('/workspace/ddl/edit', {
        //   id: this.userId,
        //   task_id: this.editingDdl.task_id,
        //   deadline: `${this.editingDdl.deadline.year}-${this.editingDdl.deadline.month}-${this.editingDdl.deadline.day}`,
        //   task_content: this.editingDdl.task_content,
        //   important: this.editingDdl.important,
        // });
        this.importantInt = this.isImportantTrue(this.editingDdl.important);
        const response = await WorkSpaceAPI.saveEditDdl(this.editingDdl.task_id,
         `${this.editingDdl.deadline.year}-${this.editingDdl.deadline.month}-${this.editingDdl.deadline.day}`,
          this.editingDdl.task_content,
          this.importantInt,);
        console.log(response.data);
        if (response.status === 200) {
          alert('DDL 修改成功');
          this.showEditDdl = false; // 关闭编辑框
          this.fetchDdlList(); // 刷新待完成 DDL 列表
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法修改 DDL:', error);
        alert('修改 DDL 失败！');
      }
    },

    // 取消编辑
    cancelEdit() {
      this.showEditDdl = false; // 取消编辑
    },

    // 保存新建的 DDL
    async saveDdl() {
      if (!this.validateDeadline()) {
        alert(this.errorMessage);
        return;
      }
      try {
        const deadline = `${this.newDdl.deadline.year}-${String(this.newDdl.deadline.month).padStart(2, '0')}-${String(this.newDdl.deadline.day).padStart(2, '0')} ${String(this.newDdl.deadline.hour).padStart(2, '0')}:${String(this.newDdl.deadline.minute).padStart(2, '0')}`;
        // const response = await axios.post('/workspace/ddl/create', {
        //   id: this.userId,
        //   deadline: deadline,
        //   task_content: this.newDdl.task_content,
        //   important: this.newDdl.important,
        // });
        this.importantInt = this.isImportantTrue(this.newDdl.important);
        const response = await WorkSpaceAPI.saveDdl(
          deadline,
          this.newDdl.task_content,
          this.importantInt,);
        console.log(response.data);
        if (response.status === 200) {
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
    async updateDdlStatus() {
      try {
        const response = await WorkSpaceAPI.updateDdl(this.itemToUpdate.task_id);
        if (response.status === 200) {
          this.fetchDdlList();
          this.fetchCompletedDdl();
        } else {
          alert(response.message);
        }
      } catch (error) {
        console.error('无法更新 DDL 状态:', error);
        alert('更新 DDL 状态失败！');
      }
      this.showUpdateConfirm = false;
      this.itemToUpdate = null;
    },

    // 删除某条 DDL
    async deleteDdl() {
      try {
        const response = await WorkSpaceAPI.deleteDDL(this.itemToDelete.task_id);
        if (response.status === 200) {
          // 刷新待完成和已完成的
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
      this.showDeleteConfirm = false;
      this.itemToDelete = null;
    },

    // 显示删除确认框
    confirmDelete(item) {
      this.itemToDelete = item;
      this.showDeleteConfirm = true;
    },
    // 取消删除操作
    cancelDelete() {
      this.showDeleteConfirm = false;
      this.itemToDelete = null;
    },
    confirmUpdate(item) {
      this.itemToUpdate = item;
      this.showUpdateConfirm = true;
    },
    cancelUpdate() {
      this.showUpdateConfirm = false;
      this.itemToUpdate = null;
      this.fetchDdlList();
    },
  },
  computed: {
    // 根据 showCompleted 控制左侧待完成的 DDL 宽度
    leftDdlStyle() {
      return {
        width: this.showCompleted ? '55%' : '100%', // 显示已完成时为55%，否则占满整个页面
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

.edit-ddl-modal,
.create-ddl-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  align-items: center;
  z-index: 1001;
}

.modal-content {
  background-color:var(--background-color);
  color: var(--text-color);
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

.important-label {
  color: rgb(247, 115, 115);
  font-weight: bold;
  margin-right: 20px;
  padding: 2px 5px;
  border-radius: 3px;
}

.invisible_important-label {
  font-weight: bold;
  margin-left: 10px;
  margin-right: 20px;
  padding: 2px 5px;
  border-radius: 3px;
}

.deadline-inputs input {
  margin-left: 4px;
  margin-right: 1px;
  width: 65px;
}

.modal-actions {
  margin-top: 20px;
  display: flex;
  justify-content: space-around;
}

.save-btn, .cancel-btn {
  padding: 10px;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border: none;
  cursor: pointer;
}

.save-btn:hover, .cancel-btn:hover {
  background-color: var(--button-background-color2);
}

.ddl-container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.ddl-left, .ddl-right {
  width: 44%;
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
  border-bottom: 1px solid var(--background-color2);
  font-size: var(--font-size);
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
  color: var(--text-color);
  opacity: 70%;
  flex-grow: 1;
}

.ddl-item .delete-btn {
  background-color: var(--button-background-color);
  color: var(--button-background-color);
  border: none;
  padding: 5px 10px;
  cursor: pointer;
}

.ddl-item .delete-btn:hover {
  background-color:var(--button-background-color2);
}

.completed {
  background-color: var(--sidebar-background-color);
  color: var(--sidebar-text-color);
}

.toggle-btn-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.toggle-btn {
  padding: 10px;
  background-color: var(--button-background-color);
  color: var(--text-color);
  border: none;
  cursor: pointer;
}

.toggle-btn:hover {
  background-color: var(--button-background-color2);
}

.trash-icon:hover, 
.create_note_icon:hover:hover{
  cursor: pointer;
  /* background-color: #dacfdb; */
}

.create_note_icon {
  width: 35px;
  height: 35px;
  cursor: pointer;
  margin-left: 10px;
  object-fit: contain; /* 确保图片按比例缩放 */
  flex-shrink: 0;
  align-self: center;
}

.important-icon, .trash-icon {
  width: 18px;
  height: 18px;
  background-color: none;
  color: var(--text-color);
  margin-right: 13px;
}


.icon{
  width: 50px;
  height: 50px;
  margin-right: 5px;
}
.ddl_header {
  display: flex;
  align-items: center; /* 垂直居中图标和文字 */
  justify-content: center; /* 水平居中 */
}
.confirm-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1001;
}

.confirm-content {
  background-color: var(--background-color);
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

.confirm-btn, .cancel-btn {
  padding: 10px;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border: none;
  cursor: pointer;
}

</style>
