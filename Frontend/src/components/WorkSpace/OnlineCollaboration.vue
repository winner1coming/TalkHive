<template>
  <div class="docs">
    <!-- 标题部分 -->
    <div class="header">
      <h2>
        <div class="doc_header">
          <img src="@/assets/icon/edit.png" alt="代码图标" class="icon"/>
          共享文档
          <img
            src="@/assets/icon/create_note.png"
            alt="添加文档"
            class="create_note_icon"
            @click="showCreateDoc = true"
          />
        </div>
      </h2>
      <!-- <button @click="showCreateDoc = true" class="new-btn">+</button> -->
    </div>

    <!-- 新建文件编辑框 -->
    <div v-if="showCreateDoc" class="create-doc-modal">
      <div class="modal-content">
        <h3>新建文件</h3>
        <label for="docname">文件名：</label>
        <input
          type="text"
          id="docname"
          v-model="newDoc.docname"
          placeholder="输入文件名"
        />

        <div class="modal-actions">
          <button @click="saveDoc" class="save-btn">确定</button>
          <button @click="cancelCreate" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 文件列表 -->
    <ul>
      <li v-for="doc in docs" :key="doc.doc_id" class="doc-item">
        <div class="docname" @click="editDoc(doc)">{{ doc.doc_name }}</div>
        <p class="modified"> - 修改时间: {{ doc.last_modified_time }}</p>
        <p class="owner-name"> - 创建者：{{ doc.owner_name }}</p>
        <!-- 展示下拉框（转发、删除等），目前只做了转发功能，美观起见就只放一个转发好了-->
        <!-- <button class="more-btn" @click="toggleDropdown(doc.doc_id)">...</button>
        <div v-if="doc.showDropdown" class="dropdown">
          <button class="dropdown_delete_btn" @click="showFriendSelect(doc.doc_id,doc.doc_name)">转发</button>
          <button class="dropdown_delete_btn" @click="confirmDelete(doc.doc_id)">删除</button>
        </div> -->
        <img
          src="@/assets/icon/share.png"
          alt="转发"
          class="share_icon"
          @click="showFriendSelect(doc.doc_id,doc.doc_name)"
        />        
      </li>
    </ul>

    <!-- 删除确认框 -->
    <div v-if="showDeleteConfirm" class="confirm-modal">
      <div class="confirm-content">
        <p>是否确认删除此文件？</p>
        <div class="modal-actions">
          <button @click="deleteDoc" class="confirm-btn">确认</button>
          <button @click="cancelDelete" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>
    <SelectFriend
      v-if="showSelectFriend"
      @close="showSelectFriend = false"
      @forwordNote="forwardDoc"
    />
  </div>
  
</template>

<script>
import * as WorkSpaceAPI from '@/services/workspace_api';
import * as chatListAPI from '@/services/chatList';
import SelectFriend from '@/components/WorkSpace/SelectFriend.vue';

export default {
  name: "Docs",
  components: {
    SelectFriend,
  },
  data() {
    return {
      showCreateDoc: false, // 控制新建文件编辑框的显示
      newDoc: {
        docname: '', // 用户输入的文件名
        doctype: '.md', // 默认文件格式
      },
      // docs: [
      //   { id: 1, docname: "Vue学习笔记.md", lastModified: "2024-12-01 10:30", showDropdown: false },
      //   { id: 2, docname: "pythontest.py", lastModified: "2024-12-05 14:15", showDropdown: false },
      //   { id: 3, docname: "javascripttest.js", lastModified: "2024-12-10 09:45", showDropdown: false }
      // ],
      docs:[],
      showDeleteConfirm: false,
      docToDelete: null, // 用于存储待删除文件的id
      showSelectFriend: false,
      docToSend:{
        doc_id: null,
        doc_name: null,
      }
    };
  },
  mounted() {
    this.fetchDocs();  // 获取笔记数据
  },
  methods: {
    async fetchDocs() {
      // 从后端获取列表
      try {
        const response = await WorkSpaceAPI.getDocs();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          this.docs = response.data.map(doc => ({
            ...doc,       // 保留原来的属性
            showDropdown: false // 添加新的属性
          }));
        } else {
          alert('获取协作文档列表失败');
        }
      } catch (error) {
        console.error('无法获取代码列表:', error);
        alert('获取代码列表失败');
      }
    },

    // 显示新建文件编辑框
    showCreateDocModal() {
      this.showCreateDoc = true;
    },

    // 取消新建文件
    cancelCreate() {
      this.showCreateDoc = false;
    },

    // 新建文件
    async saveDoc() {
      try {
        const response = await WorkSpaceAPI.createDoc(this.newDoc.docname);

        if (response.status === 200) {

          this.showCreateDoc = false;
          this.fetchDocs();
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法创建文件:', error);
        alert('创建文件失败！');
      }
      // this.editDoc(1);
    },

    // 显示转发好友选择框
    showFriendSelect(doc_id, name) {
      this.docToSend.doc_id = doc_id;
      this.docToSend.doc_name = name;
      this.showSelectFriend = true;
    },
    // 转发文档
    async forwardDoc(tid) {
      this.showSelectFriend = false;
      try {
        // 构造json对象，包括docid，docname
        const json_doc = JSON.stringify({doc_id: this.docToSend.doc_id, doc_name: this.docToSend.doc_name});
        const response = await chatListAPI.sendMessage(tid, json_doc, "collab_doc", false);
        if (response.status === 200) {
          this.$root.notify('转发成功', 'success');
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      } catch (error) {
        console.error('无法转发文档:', error);
      }
    },

    // 跳转到文件编辑页
    editDoc(doc) {
      // this.$router.push(`/workspace/doc/${id}`);
      
      // 使用 Vuex 更新 currentDoc 对象
      this.$store.dispatch('updateCurrentDoc', {
        doc_id: doc.doc_id,
        doc_name: doc.doc_name,
      });

      // 跳转到编辑页面
      this.$router.push(`/workspace/collabdocs/editor`);
    },

    // 切换下拉框显示/隐藏
    toggleDropdown(doc_id) {
      const doc = this.docs.find(n => n.doc_id === doc_id);
      doc.showDropdown = !doc.showDropdown;
    },

    // 显示删除确认框
    confirmDelete(id) {
      this.docToDelete = id;
      this.showDeleteConfirm = true;
    },

    // 删除文件
    async deleteDoc() {
      try {
        // 发送删除请求到后端
        const response = await WorkSpaceAPI.deleteDoc(this.docToDelete);
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
      this.toggleDropdown(this.docToDelete);
      this.docToDelete = null;
      this.fetchDocs();
    },

    // 取消删除操作
    cancelDelete() {
      this.showDeleteConfirm = false;
      this.toggleDropdown(this.docToDelete);
      this.docToDelete = null;
    },
  },
};
</script>

<style scoped>
/* 样式部分 */
.docs {
  padding: 20px;
}

.docs h2 {
  margin-bottom: 20px;
}

.header {
  display: flex;
  justify-content: center;
  align-items: center;
}

.new-btn {
  background-color: var(--button-background-color);
  color: var(--text-color);
  border: none;
  cursor: pointer;
  margin-left: 20px;
  margin-bottom: 15px;
  border-radius: 50%;
  width: 45px;
  font-size: 36px;
}

.new-btn:hover {
  background-color: var(--button-background-color2);
  color: rgb(134, 154, 233);
}

.create-doc-modal {
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
  background-color: var(--background-color);
  color: var(--text-color);
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

.save-btn, .cancel-btn{
  padding: 10px;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border: none;
  cursor: pointer;
}

.save-btn:hover, .cancel-btn:hover {
  background-color: var(--button-background-color2);
}

.docs ul {
  list-style: none;
  padding: 0;
}

.doc-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--background-color2);
  justify-content: space-between;
  position: relative; /* 给父元素设置相对定位 */
  color: var(--text-color);
}

.doc-item .docname {
  font-weight: bold;
  margin-right: 10px;
  cursor: pointer;
}

.doc-item .modified, 
.doc-item .owner-name {
  color: var(--text-color);
  opacity: 80%;
  cursor: default;
}

.more-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 23px;
  color: var(--text-color);
}

.more-btn:hover {
  color: var(--button-background-color2);
}

.dropdown {
  display: flex;
  flex-direction: column;
  position: absolute;
  right: -25px;
  bottom: 23px;
  background-color: var(--background-color);
  border: 1px solid #ddd;
  border-radius: 5px;
  box-shadow: 0 2px 8px var(--background-color2);
  z-index: 10;  /* 确保 dropdown 在按钮上方显示 */
}

/* .dropdown:hover {
  background-color: rgb(208, 208, 208);
} */

.dropdown_delete_btn{
  margin: 5px; 
  padding:5px;
  border: none; 
  color:var(--text-color); 
  background-color:var(--background-color);
}

.dropdown_delete_btn:hover{
  background-color: var(--background-color1);
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

.confirm-btn:hover, .cancel-btn:hover {
  background-color: var(--button-background-color2);
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

.doc_header {
  display: flex;
  align-items: center; /* 垂直居中图标和文字 */
  justify-content: center; /* 水平居中 */
}

.icon{
  width: 50px;
  height: 50px;
  margin-right: 5px;
}

.share_icon {
  width: 20px;
  height: 20px;
  margin-right: 10px;
}

.share_icon:hover {
  cursor: pointer;
}
</style>
