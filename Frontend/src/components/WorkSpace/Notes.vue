<template>
    <div class="notes">
      <div class="header">
        <h2>我的笔记</h2>
        <button @click="showCreateFile = true" class="new-btn">+</button>
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
            <option value=".md">.md</option>
            <option value=".txt">.txt</option>
            <option value=".docx">.docx</option>
            <!-- 可以根据需求继续添加文件类型 -->
          </select>

          <div class="modal-actions">
            <button @click="saveFile" class="save-btn">确定</button>
            <button @click="cancelCreate" class="cancel-btn">取消</button>
          </div>
        </div>
      </div>

      <ul>
        <li v-for="note in notes" :key="note.id" class="note-item" @click="editNote(note.id)">
            <span class="filename">{{ note.filename }}</span>
            <span class="modified"> - 上次修改时间: {{ note.lastModified }}</span>
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  export default {
    name: "Notes",
    data() {
      return {
        showCreateFile: false, // 控制新建文件编辑框的显示
        newFile: {
          filename: '',  // 用户输入的文件名
          filetype: '.md', // 默认文件格式
        },
        notes: [
          { id: 1, filename: "Vue学习笔记.md", lastModified: "2024-12-01 10:30" },
          { id: 2, filename: "项目需求分析.docx", lastModified: "2024-12-05 14:15" },
          { id: 3, filename: "代码优化方案.txt", lastModified: "2024-12-10 09:45" }
        ]
      };
    },
    methods: {
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
          try {
            // 发送请求到后端，保存新建的文件
            const response = await axios.post('/workspace/create-file', {
              filename: this.newFile.filename + this.newFile.filetype,  // 文件名和文件格式拼接
            });

            if (response.data.status === 200) {
              alert('文件创建成功');
              this.showCreateFile = false;  // 关闭编辑框
            } else {
              alert(response.data.message);
            }
          } catch (error) {
            console.error('无法创建文件:', error);
            alert('创建文件失败！');
          }
          this.editNote(1)
        },
        editNote(id) {
            this.$router.push(`/workspace/notes/${id}`);
        },
    },
  };
  </script>
  
  <style scoped>
  .notes {
    padding: 20px;
  }
  
  .notes h2 {
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
    width: 45px; /* 宽度 */
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

  .notes ul {
    list-style: none;
    padding: 0;
  }
  
  .note-item {
    display: flex;
    align-items: center;
    padding: 10px 0;
    border-bottom: 1px solid #ddd;
  }
  
  .note-item .filename {
    font-weight: bold;
    margin-right: 10px;
  }
  
  .note-item .modified {
    color: #666;
  }
  </style>