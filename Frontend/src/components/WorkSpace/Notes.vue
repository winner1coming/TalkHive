<template>
  <div class="notes">
    <div class="header">
      <h2>
        <div class="note_header">
          <img src="@/assets/icon/edit.png" alt="笔记图标" class="icon"/>
          我的笔记
          <img
            src="@/assets/icon/create_note.png"
            alt="添加笔记"
            class="create_note_icon"
            @click="showCreateFile = true"
          />
        </div>
      </h2>
        <!-- <img src="@/assets/icon/note_icon.png" alt="笔记图标" class="icon"/> -->
      <!-- <button @click="showCreateFile = true" class="new-btn">+</button> -->

    </div>

    <!-- 分类标签 -->
    <div class="category-tags">
      <span
        v-for="category in categories"
        :key="category"
        class="category-tag"
        :class="{ active: selectedCategory === category }"
        @click="filterByCategory(category)"
        @contextmenu="showEditCategoryDialog($event, category)"
      >
        {{ category }}
      </span>

      <!-- 新建标签图标 -->
      <img
        src="@/assets/icon/add_tag.png"
        alt="新增标签"
        class="tag_icon"
        @click="showCreateCategory = true"
      />

      <!-- 删除标签图标 -->
      <img
        src="@/assets/icon/delete_tag.png"
        alt="删除标签"
        class="tag_icon"
        @click="showDeleteCategory = true"
      />
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
        </select>

        <label for="category">分类：</label>
        <select v-model="newFile.category" id="category">
          <option v-for="category in categories" :key="category" :value="category">{{ category }}</option>
        </select>

        <div class="modal-actions">
          <button @click="saveFile" class="save-btn">确定</button>
          <button @click="cancelCreate" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 新建分类编辑框 -->
    <div v-if="showCreateCategory" class="create-category-modal">
      <div class="modal-content">
        <h3>新建分类</h3>
        <input
          type="text"
          v-model="newCategory"
          placeholder="输入新分类名称"
        />
        <div class="modal-actions">
          <button @click="saveCategory" class="save-btn">保存</button>
          <button @click="cancelCreateCategory" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 删除分类确认框 -->
    <div v-if="showDeleteCategory" class="delete-category-modal">
      <div class="modal-content">
        <h3>删除分类</h3>
        <p>请选择要删除的分类：</p>
        <select v-model="selectedCategoryToDelete">
          <option v-for="category in categories" :key="category" :value="category">
            {{ category }}
          </option>
        </select>
        <div class="modal-actions">
          <button @click="deleteCategory" class="delete-btn">删除</button>
          <button @click="cancelDeleteCategory" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 修改分类名称编辑框 -->
    <div v-if="showEditCategory" class="edit-category-modal">
      <div class="modal-content">
        <h3>修改分类名称</h3>
        <input
          type="text"
          v-model="editedCategoryName"
          placeholder="输入新的分类名称"
        />
        <div class="modal-actions">
          <button @click="saveCategoryName" class="save-btn">保存</button>
          <button @click="cancelEditCategory" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>
    
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

    <ul>
      <li v-for="note in filteredNotes" :key="note.id" class="note-item">
        <span class="filename" @click="editNote(note)">{{ note.filename+".md" }}</span>
        <span class="modified"> - 上次修改时间: {{ note.lastModified }}</span>
        <button class="more-btn" @click="toggleDropdown(note.id)">...</button>
        <div v-if="note.showDropdown" class="dropdown">
          <button class="dropdown_delete_btn" @click="showFriendSelect(note.id,note.filename, '.md')">转发</button>
          <button class="dropdown_delete_btn" @click="confirmDelete(note.id)">删除</button>
        </div>
      </li>
    </ul>
    <SelectFriend
      v-if="showSelectFriend"
      @close="showSelectFriend = false"
      @forwordNote="forwardNote"
    />
  </div>
  
</template>

<script>
import * as WorkSpaceAPI from '@/services/workspace_api';
import * as chatListAPI from '@/services/chatList';
import SelectFriend from '@/components/WorkSpace/SelectFriend.vue';

export default {
  name: "Notes",
  components: { SelectFriend },
  data() {
    return {
      showCreateFile: false,
      showCreateCategory: false,
      showDeleteCategory: false,
      showEditCategory: false,  
      editedCategoryName: '',  // 用于编辑的分类名称
      beforeEditedCategoryName:'',
      newCategory: '',  // 用于新建分类的输入框
      selectedCategory: '',  // 当前选中的分类
      selectedCategoryToDelete: '',  // 用户选择删除的分类
      showDeleteConfirm: false,
      fileToDelete: null, // 用于存储待删除文件的id
      newFile: {
        filename: '',
        filetype: '.md',
        category: '',  // 选择的分类
      },
      categories: [],  // 模拟的分类列表
      notes: [],  // 所有笔记数据
      showSelectFriend: false,
      forwardCodeContent:{
        code_id: null,
        name: null,
        Suffix: null,
      }
    };
  },
  computed: {
    // 根据选中的分类筛选笔记
    filteredNotes() {
      if (this.selectedCategory) {
        return this.notes.filter(note => note.category === this.selectedCategory);
      }
      return this.notes;
    },
  },
  mounted() {
    this.fetchCategories();
    this.fetchNotes();  // 获取笔记数据
  },
  methods: {
    async fetchCategories() {
      // 从后端获取分类列表
      try {
        const response = await WorkSpaceAPI.getCategories();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          this.categories = response.data.categories;
        } else {
          alert('获取分类列表失败');
        }
      } catch (error) {
        console.error('无法获取分类列表:', error);
        alert('获取分类列表失败');
      }
    },
    async fetchNotes() {
      // 从后端获取笔记列表
      try {
        const response = await WorkSpaceAPI.getNotes();
        if (response.status === 200) {
          if(!response.data)
          {
            return;
          }
          // 为每个 note 增加 showDropdown: false
          this.notes = response.data.map(note => ({
            ...note,       // 保留原来的属性
            showDropdown: false // 添加新的属性
          }));
          console.log(this.notes);
          //this.notes = response.data.notes;
        } else {
          alert('获取笔记列表失败');
        }
      } catch (error) {
        console.error('无法获取笔记列表:', error);
        alert('获取笔记列表失败');
      }
    },

    filterByCategory(category) {
      if(this.selectedCategory == category)
      {
        this.selectedCategory = '';
      }
      else{
        this.selectedCategory = category;  // 设置选中的分类
      }
    },

    showCreateFileModal() {
      this.showCreateFile = true;
    },

    cancelCreate() {
      this.showCreateFile = false;
    },

    //保存新建的文件
    async saveFile() {
      try {
        const response = await WorkSpaceAPI.createNote(this.newFile.filename, this.newFile.category);

        if (response.status === 200) {
          this.showCreateFile = false;
          this.fetchNotes();
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法创建文件:', error);
        alert('创建文件失败！');
      }
    },

    // 显示转发好友选择框
    showFriendSelect(code_id, name, Suffix) {
      this.forwardCodeContent.code_id = code_id;
      this.forwardCodeContent.name = name;
      this.forwardCodeContent.Suffix = Suffix;
      this.showSelectFriend = true;
    },
    // 转发笔记
    async forwardNote(tid) {
      this.showFriendSelect = false;
      try {
        // 获取文件
        const response = await WorkSpaceAPI.getNoteContent(this.forwardCodeContent.code_id);
        const content = response.data;
        const blob = new Blob([content], { type: this.forwardCodeContent.Suffix.slice(1) });
        const file = new File([blob], this.forwardCodeContent.name+this.forwardCodeContent.Suffix, { type: this.forwardCodeContent.Suffix.slice(1) });
        console.log(file);
        // 转发文件
        const formData = new FormData();
        formData.append('tid', tid);
        formData.append('content', content);
        formData.append('is_group', false);
        response = await chatListAPI.sendFile(formData);
        if (response.status === 200) {
          this.$root.notify('转发成功', 'success');
        } else {
          this.$root.notify(response.data.message, 'error');
        }
      } catch (error) {
        console.error('无法转发笔记:', error);
      }
    },

    // 编辑笔记
    editNote(note) {
      // 使用 Vuex 更新 currentNote 对象
      this.$store.dispatch('updateCurrentNote', {
        note_id: note.id,
        filename: note.filename,
        category: note.category
      });
      this.$store.dispatch('updateCategories', this.categories);

      // 跳转到编辑页面
      this.$router.push(`/workspace/notes/editor`);
    },
    
    // 切换下拉框显示/隐藏
    toggleDropdown(note_id) {
      const note = this.notes.find(n => n.id === note_id);
      note.showDropdown = !note.showDropdown;
    },
    // 显示删除确认框
    confirmDelete(id) {
      this.fileToDelete = id;
      this.showDeleteConfirm = true;
      console.log("this.fileToDelete:",this.fileToDelete);
    },

    // 删除文件
    async deleteFile() {
      try {
        // 发送删除请求到后端
        const response = await WorkSpaceAPI.deleteNote(this.fileToDelete);
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
      this.fetchNotes();
    },
    // 取消删除操作
    cancelDelete() {
      this.showDeleteConfirm = false;
      this.toggleDropdown(this.fileToDelete);
      this.fileToDelete = null;
    },

    showCreateCategoryModal() {
      this.showCreateCategory = true;
    },

    cancelCreateCategory() {
      this.showCreateCategory = false;
    },

    //新建分类
    async saveCategory() {
      if (this.newCategory.trim() === '') {
        alert('分类名称不能为空');
        return;
      }
      else if (this.categories!= null && this.categories.indexOf(this.newCategory)!=-1)
      {
        alert('已存在该分类');
        return;
      }
      try {
        const response = await WorkSpaceAPI.saveCategory(this.newCategory);
        if (response.status === 200) {
          //this.categories.push(this.newCategory);  // 添加新分类
          this.fetchCategories();
          this.newCategory = '';  // 清空输入框
        } else {
          alert(response.message);
        }
      } catch (error) {
        console.error('无法创建分类:', error);
        alert('创建分类失败！');
      }
      this.showCreateCategory = false;  // 关闭新建分类窗口
    },
    
    // 显示删除分类的模态框
    cancelDeleteCategory() {
      this.showDeleteCategory = false;
    },

    // 删除选中的分类
    async deleteCategory() {
      const index = this.categories.indexOf(this.selectedCategoryToDelete);
      if (index > -1) {
        try {
        const response = await WorkSpaceAPI.deleteCategory(this.selectedCategoryToDelete);
        if (response.status === 200) 
        {
          // this.categories.splice(index, 1);  // 从分类列表中删除
          this.fetchCategories();
          this.fetchNotes();
        } else {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('删除分类列表失败:', error);
        alert('删除分类列表失败');
      }
      }
      this.selectedCategoryToDelete = '';  // 清空选择的分类
      this.showDeleteCategory = false;
    },

    // 显示编辑分类名称的模态框
    showEditCategoryDialog(event, category) {
      event.preventDefault();  // 阻止右键菜单的默认行为
      this.editedCategoryName = category;  // 设定当前编辑的分类名称
      this.beforeEditedCategoryName = category;
      this.showEditCategory = true;
    },

    // 保存修改的分类名称
    async saveCategoryName() {
      const index = this.categories.indexOf(this.beforeEditedCategoryName);
      if (index > -1) {
        try {
          const response = await WorkSpaceAPI.saveEditCategory(this.beforeEditedCategoryName,this.editedCategoryName);
          if (response.status === 200) {
            this.categories[index] = this.editedCategoryName;  // 更新分类名称
            this.beforeEditedCategoryName = '';  // 清空输入框
            this.editedCategoryName = '';
          } else {
            alert(response.data.message);
          }
        } catch (error) {
          console.error('无法修改分类:', error);
          alert('修改分类失败！');
        }
      }
      this.showEditCategory = false;  // 关闭编辑框
    },

    // 取消编辑分类名称
    cancelEditCategory() {
      this.showEditCategory = false;  // 关闭编辑框
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
  width: 45px;
  font-size: 36px;
}

.new-btn:hover {
  background-color: #0056b3;
  color: rgb(134, 154, 233);
}

.category-tags {
  display: flex;
  gap: 10px;
  margin: 20px 0;
}

.category-tag {
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  background-color: var(--button-background-color);
}

.category-tag:hover {
  background-color: var(--button-background-color1);
}

.category-tag.active, .category-tag.active:hover {
  background-color: var(--button-background-color2);
  color: var(--button-text-color);
}

.new-category {
  border-radius: 50%;
  height: 25px;
  margin-top: 4px;
  background-color: #86d799;
  color: white;
}

.delete-category {
  border-radius: 50%;
  height: 25px;
  width: 15px;
  margin-top: 4px;
  background-color: #f88d93;
  color: white;
}

.create-file-modal,
.edit-category-modal,
.create-category-modal,
.delete-category-modal,
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
  z-index: 2000; /* 保证弹出框在遮罩层上方 */
}

.modal-content {
  background-color:var(--background-color);
  color: var(--text-color);
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

input[type="text"],
select {
  width: 100%;
  padding: 8px;
  margin: 10px 0;
  border: 1px solid var(--background-color2);
  border-radius: 4px;
}

.modal-actions {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
}

.confirm-btn, 
.delete-btn,
.save-btn,
.cancel-btn {
  padding: 10px;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border: none;
  cursor: pointer;
}

.confirm-btn:hover, 
.save-btn:hover,
.cancel-btn:hover {
  background-color: var(--button-background-color2);
}

.notes ul {
  list-style: none;
  padding: 0;
}

.note-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--background-color2);
  justify-content: space-between;
  position: relative; /* 给父元素设置相对定位 */
  color: var(--text-color);
}

.note-item .filename {
  font-weight: bold;
  margin-right: 10px;
  cursor: pointer;
}

.note-item .modified {
  color: var(--text-color);
  opacity: 70%;
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
  border: 1px solid var(--background-color2);
  border-radius: 5px;
  box-shadow: 0 2px 8px var(--background-color2);
  z-index: 10;  /* 确保 dropdown 在按钮上方显示 */
}

.dropdown_delete_btn{
  margin: 5px; 
  padding:5px;
  border: none; 
  color:var(--text-color); 
  background-color: var(--background-color);
}

.dropdown_delete_btn:hover{
  background-color: var(--background-color1);
}

/* .confirm-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
} */

.confirm-content {
  background-color: var(--background-color);
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

.tag_icon {
  width: 35px;
  height: 35px;
  cursor: pointer;
  margin-left: 5px;
  object-fit: contain; /* 确保图片按比例缩放 */
  flex-shrink: 0;
  align-self: center;
}
.create_note_icon:hover, .tag_icon:hover{
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
.icon{
  width: 50px;
  height: 50px;
  margin-right: 5px;
}
.note_header {
  display: flex;
  align-items: center; /* 垂直居中图标和文字 */
  justify-content: center; /* 水平居中 */
  color: var(--text-color);
}
</style>
