<template>
    <div class="recycle-bin">
      <h2>回收站</h2>
      <ul>
        <li v-for="file in deletedFiles" :key="file.recycle_id" class="file-item">
          <div class="file-info">
            <!--文件类型-->
            <span class="file_type">{{ getTypeLabel(file.type) }}</span>
            <!-- 文件名在左边 -->
            <span class="file-name">{{ file.filename }}</span>
            <!-- 删除时间在右边 -->
            <span class="deleted-time">{{ file.recycle_time }}</span>
          </div>
          <div class="file-actions">
            <!-- 恢复按钮 -->
            <button @click="confirmRestore(file)" class="restore-btn">恢复</button>
            <!-- 删除按钮 -->
            <button @click="confirmDelete(file)" class="delete-btn">彻底删除</button>
          </div>
        </li>
      </ul>
  
      <!-- 恢复确认弹窗 -->
      <div v-if="showRestoreModal" class="modal-overlay">
        <div class="modal">
          <h3>确认恢复文件</h3>
          <p>你确定要恢复 "{{ selectedFile?.filename }}" 吗？</p>
          <button @click="restoreFile" class="confirm-btn">确认</button>
          <button @click="closeModal" class="cancel-btn">取消</button>
        </div>
      </div>
  
      <!-- 删除确认弹窗 -->
      <div v-if="showDeleteModal" class="modal-overlay">
        <div class="modal">
          <h3>确认删除文件</h3>
          <p>你确定要彻底删除 "{{ selectedFile?.filename }}" 吗？</p>
          <button @click="deleteFile" class="confirm-btn">确认</button>
          <button @click="closeModal" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import * as WorkSpaceAPI from '@/services/workspace_api';
  
  export default {
    name: 'RecycleBin',
    data() {
      return {
        deletedFiles: [], // 存储已删除文件列表
        showRestoreModal: false, // 是否显示恢复确认弹窗
        showDeleteModal: false, // 是否显示删除确认弹窗
        selectedFile: null, // 当前操作的文件对象
      };
    },
    mounted() {
      this.fetchDeletedFiles();
    },
    methods: {
      // 从后端获取已删除文件列表
      async fetchDeletedFiles() {
        try {
          const response = await WorkSpaceAPI.getRecycles();
          this.deletedFiles = response.data.files;
        } catch (error) {
          console.error('获取已删除文件列表失败:', error);
        }
      },
  
      // 点击恢复按钮时，显示恢复确认弹窗
      confirmRestore(file) {
        this.selectedFile = file;
        this.showRestoreModal = true;
      },
  
      // 点击删除按钮时，显示删除确认弹窗
      confirmDelete(file) {
        this.selectedFile = file;
        this.showDeleteModal = true;
      },
  
      // 关闭弹窗
      closeModal() {
        this.showRestoreModal = false;
        this.showDeleteModal = false;
        this.selectedFile = null;
      },
  
      // 恢复文件
      async restoreFile() {
        try {
          const response = await WorkSpaceAPI.restoreFile(this.selectedFile.type, this.selectedFile.recycle_id);
  
          if (response.data.status === 200) {
            this.fetchDeletedFiles(); // 重新获取已删除文件列表
          } else {
            alert(response.data.message);
          }
        } catch (error) {
          console.error('恢复文件失败:', error);
          alert('恢复文件失败！');
        } finally {
          this.closeModal();
        }
      },
  
      // 删除文件
      async deleteFile() {
        try {
          const response = await WorkSpaceAPI.deleteFile(this.selectedFile.type, this.selectedFile.recycle_id);
          if (response.data.status === 200) {
            console.log(response.data.message);
            this.fetchDeletedFiles(); // 重新获取已删除文件列表
          } else {
            alert(response.data.message);
          }
        } catch (error) {
          console.error('删除文件失败:', error);
          alert('删除文件失败！');
        } finally {
          this.closeModal();
        }
      },
      getTypeLabel(type) {
        switch (type) {
          case 'code': return '代码';
          case 'note': return '笔记';
          default: return '未知';
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .recycle-bin {
    padding: 20px;
  }
  
  .recycle-bin h2 {
    margin-bottom: 20px;
  }
  
  .file-item {
    display: flex;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid #ddd;
  }
  
  .file-info {
    display: flex;
    justify-content: space-between;
    width: 100%;
  }
  
  .file_type,
  .file-name {
    font-weight: bold;
    text-align: left;
    flex: 1;
  }
  
  .deleted-time {
    color: #666;
    font-size: 20px;
    text-align: right;
    margin-right: 50px;
    margin-top: 8px;
  }
  
  .file-actions {
  display: flex;
  gap: 10px;
  align-items: center; /* 垂直居中按钮 */
}

.file-actions button {
  flex-shrink: 0; /* 防止按钮被压缩 */
  white-space: nowrap; /* 确保按钮文本不换行 */
}

  .restore-btn, .delete-btn {
    padding: 5px 10px;
    background-color: #007bff;
    color: white;
    border: none;
    cursor: pointer;
  }
  
  .restore-btn:hover, .delete-btn:hover {
    background-color: #0056b3;
  }
  
  .modal-overlay {
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
  
  .modal {
    background-color: white;
    padding: 20px;
    border-radius: 10px;
    width: 300px;
    text-align: center;
  }
  
  .confirm-btn, .cancel-btn {
    padding: 10px;
    background-color: #007bff;
    color: white;
    border: none;
    cursor: pointer;
    margin-top: 10px;
    margin-left: 20px;
  }
  
  .confirm-btn:hover, .cancel-btn:hover {
    background-color: #0056b3;
  }
  </style>
  