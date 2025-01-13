<template>
  <div class="codes">
    <!-- 标题部分 -->
    <div class="header">
      <h2>
        <div class="code_header">
          <img src="@/assets/icon/code.png" alt="代码图标" class="icon"/>
          我的代码
          <img
            src="@/assets/icon/create_note.png"
            alt="添加代码"
            class="create_note_icon"
            @click="showCreateFile = true"
          />
        </div>
      </h2>
      <!-- <button @click="showCreateFile = true" class="new-btn">+</button> -->
    </div>

    <!-- 新建文件编辑框 -->
    <div v-if="showCreateFile" class="create-file-modal">
      <div class="modal-content">
        <h3>新建文件</h3>
        <label for="filename">文件名：</label>
        <input
          type="text"
          id="filename"
          v-model="newFile.filename"
          placeholder="输入文件名"
        />
        <label for="filetype">文件格式：</label>
        <select v-model="newFile.filetype" id="filetype">
          <option value=".cpp">.cpp</option>
          <option value=".js">.js</option>
          <option value=".html">.html</option>
          <option value=".css">.css</option>
          <option value=".txt">.txt</option>
          <!-- 可以根据需求继续添加文件类型 -->
        </select>

        <div class="modal-actions">
          <button @click="saveFile" class="save-btn">确定</button>
          <button @click="cancelCreate" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 文件列表 -->
    <ul>
      <li v-for="code in codes" :key="code.code_id" class="code-item">
        <span class="filename" @click="editCode(code)">{{ code.code_name + code.Suffix }}</span>
        <span class="modified"> - 上次修改时间: {{ code.last_modified_time }}</span>
        <button class="more-btn" @click="toggleDropdown(code.code_id)">...</button>
        <div v-if="code.showDropdown" class="dropdown">
          <button class="dropdown_delete_btn" @click="confirmDelete(code.code_id)">删除</button>
        </div>
      </li>
    </ul>

    <!-- 删除确认框 -->
    <div v-if="showDeleteConfirm" class="confirm-modal">
      <div class="confirm-content">
        <p>是否确认删除此文件？</p>
        <div class="modal-actions">
          <button @click="deleteFile" class="confirm-btn">确认</button>
          <button @click="cancelDelete" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as WorkSpaceAPI from '@/services/workspace_api';

export default {
  name: "Codes",
  data() {
    return {
      showCreateFile: false, // 控制新建文件编辑框的显示
      newFile: {
        filename: '', // 用户输入的文件名
        filetype: '.md', // 默认文件格式
      },
      // codes: [
      //   { id: 1, filename: "Vue学习笔记.md", lastModified: "2024-12-01 10:30", showDropdown: false },
      //   { id: 2, filename: "pythontest.py", lastModified: "2024-12-05 14:15", showDropdown: false },
      //   { id: 3, filename: "javascripttest.js", lastModified: "2024-12-10 09:45", showDropdown: false }
      // ],
      codes:[],
      showDeleteConfirm: false,
      fileToDelete: null, // 用于存储待删除文件的id
    };
  },
  mounted() {
    this.fetchCodes();  // 获取笔记数据
  },
  methods: {
    async fetchCodes() {
      // 从后端获取代码列表
      try {
        const response = await WorkSpaceAPI.getCodes();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          this.codes = response.data.map(code => ({
            ...code,       // 保留原来的属性
            showDropdown: false // 添加新的属性
          }));
        } else {
          alert('获取代码列表失败');
        }
      } catch (error) {
        console.error('无法获取代码列表:', error);
        alert('获取代码列表失败');
      }
    },

    // 显示新建文件编辑框
    showCreateFileModal() {
      this.showCreateFile = true;
    },

    // 取消新建文件
    cancelCreate() {
      this.showCreateFile = false;
    },

    // 保存新建的文件
    async saveFile() {
      // try {
      //   // 发送请求到后端，保存新建的文件
      //   const response = await axios.post('/workspace/create-file', {
      //     filename: this.newFile.filename + this.newFile.filetype,  // 文件名和文件格式拼接
      //   });

      //   if (response.status === 200) {
      //     alert('文件创建成功');
      //     this.showCreateFile = false; // 关闭编辑框
      //   } else {
      //     alert(response.data.message);
      //   }
      // } catch (error) {
      //   console.error('无法创建文件:', error);
      //   alert('创建文件失败！');
      // }
      try {
        const response = await WorkSpaceAPI.createCode(this.newFile.filename, this.newFile.filetype);

        if (response.status === 200) {

          this.showCreateFile = false;
          this.fetchCodes();
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法创建文件:', error);
        alert('创建文件失败！');
      }
      // this.editCode(1);
    },

    // 跳转到文件编辑页
    editCode(code) {
      // this.$router.push(`/workspace/code/${id}`);
      
      // 使用 Vuex 更新 currentCode 对象
      this.$store.dispatch('updateCurrentCode', {
        code_id: code.code_id,
        filename: code.code_name + code.Suffix,
      });

      // 跳转到编辑页面
      this.$router.push(`/workspace/code/editor`);
    },

    // 切换下拉框显示/隐藏
    toggleDropdown(code_id) {
      const code = this.codes.find(n => n.code_id === code_id);
      code.showDropdown = !code.showDropdown;
    },

    // 显示删除确认框
    confirmDelete(id) {
      this.fileToDelete = id;
      this.showDeleteConfirm = true;
    },

    // 删除文件
    async deleteFile() {
      try {
        // 发送删除请求到后端
        const response = await WorkSpaceAPI.deleteCode(this.fileToDelete);
        if (response.status === 200) {
          console.log(response.data.message);
        } else {
          alert('删除失败');
        }
      } catch (error) {
        console.error('无法删除文件:', error);
        alert('删除文件失败');
      }
      this.showDeleteConfirm = false;
      this.toggleDropdown(this.fileToDelete);
      this.fileToDelete = null;
    },

    // 取消删除操作
    cancelDelete() {
      this.showDeleteConfirm = false;
      this.toggleDropdown(this.fileToDelete);
      this.fileToDelete = null;
    },
  },
};
</script>

<style scoped>
/* 样式部分 */
.codes {
  padding: 20px;
}

.codes h2 {
  margin-bottom: 20px;
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
  width: 45px;
  font-size: 36px;
}

.new-btn:hover {
  background-color: #0056b3;
  color: rgb(134, 154, 233);
}

.create-file-modal {
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

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

input[type="text"], select {
  width: 100%;
  padding: 8px;
  margin: 10px 0;
  border: 1px solid #ddd;
  border-radius: 4px;
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

.codes ul {
  list-style: none;
  padding: 0;
}

.code-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #ddd;
  justify-content: space-between;
  position: relative; /* 给父元素设置相对定位 */
}

.code-item .filename {
  font-weight: bold;
  margin-right: 10px;
  cursor: pointer;
}

.code-item .modified {
  color: #666;
}

.more-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 23px;
  color: #333;
}

.more-btn:hover {
  color: #007bff;
}

.dropdown {
  position: absolute;
  right: -25px;
  bottom: 23px;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 5px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 10;  /* 确保 dropdown 在按钮上方显示 */
}

/* .dropdown:hover {
  background-color: rgb(208, 208, 208);
} */

.dropdown_delete_btn{
  margin: 5px; 
  padding:5px;
  border: none; 
  color:#333; 
  background-color: white;
}

.dropdown_delete_btn:hover{
  background-color: rgb(208, 208, 208);
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
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

.confirm-btn, .cancel-btn {
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

.confirm-btn:hover, .cancel-btn:hover {
  background-color: #0056b3;
}

.create_note_icon:hover{
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

.code_header {
  display: flex;
  align-items: center; /* 垂直居中图标和文字 */
  justify-content: center; /* 水平居中 */
}
.icon{
  width: 50px;
  height: 50px;
  margin-right: 5px;
}
</style>
