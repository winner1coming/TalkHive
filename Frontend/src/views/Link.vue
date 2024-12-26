<template>
  <div class="web-linker-container" 
    @mouseover="handleMouseOver" 
    @mouseleave="handleMouseLeave"
  >
    <a href="#" @mouseover="showWebLinksDropdown = true">网页链接器</a>
    <div v-if="showWebLinksDropdown || isDropdownPinned" class="web-links-dropdown" @mouseleave="handleDropdownMouseLeave">
      <div class="pin-icon" @click="togglePinDropdown" :style="{ backgroundColor: isDropdownPinned ? 'red' : '#706a6a' }">
        <i class="fas fa-thumbtack">📍</i>
      </div>
      <div class="put">
      <SearchLink
        v-model="searchQuery"
        @input="handleSearchInput"
        @focus="handleFocus"
        @blur="handleBlur"
        placeholder="搜索网页链接"
      />
      <button @click="handleAddLinkClick" class="add-button">➕</button>
      </div>
      <ul class="web-links-list">
        <li v-for="(link, index) in filteredWebLinks" :key="index">
          <a :href="link.url" target="_blank">
            <img :src="link.icon" alt="Icon" class="web-icon" @error="handleIconError" />
            <div class="web-info">
              <span class="web-name">{{ link.name }}</span>
              <span class="web-url">{{ link.url }}</span>
            </div>
          </a>
          <button @click="removeWebLink(index)" class="delete-button">删除</button>
        </li>
      </ul>
    </div>
    <div v-if="showAddLinkModal" class="modal-overlay">
      <div class="modal-content">
        <h3>添加新链接</h3>
        <form @submit.prevent="addWebLink">
          <div class="form-group">
            <label for="linkName">名称:</label>
            <input type="text" id="linkName" v-model="newlink.name" required />
          </div>
          <div class="form-group">
            <label for="linkUrl">网址:</label>
            <input type="url" id="linkUrl" v-model="newlink.url" required />
          </div>
          <div class="form-actions">
            <button type="submit" class="add-button">添加</button>
            <button @click="showAddLinkModal = false,isAddingLink= false" class="cancel-button">取消</button>
          </div>
        </form>
      </div>
    </div>
    <Windows 
      :visible="showModal"
      :message="modalMessage"
      @close="showModal = false"
    />
  </div>
</template>

<script>
import SearchLink from '@/components/base/SearchLink.vue';
import Windows from '@/components/base/Windows.vue';
import link from '@/assets/images/link.png';
import { DelLinks, AddLinks } from '@/services/settingView.js';

export default {
  name: 'WebLinker',
  components: {
    SearchLink,
    Windows,
  },
  data() {
    return {
      showWebLinksDropdown: false,
      isDropdownPinned: false,
      searchQuery: '',
      defaultwebLinks: [
        { name: 'Vue', url: 'https://vuejs.org', icon: 'https://vuejs.org/logo.svg' },
        { name: 'Github', url: 'https://github.com', icon: 'https://github.githubassets.com/favicon.ico' },
        { name: 'baidu', url: 'https://baidu.com', icon: 'https://baidu.com/favicon.ico' },
      ],
      showAddLinkModal: false,
      newlink: {
        name: '',
        url: '',
        icon: '',
      },
      defaulticon: link,
      isTyping: false,
      showModal: false,
      modalMessage: '',
      isAddingLink: false, // 是否正在添加链接
      dropdownTimeout: null, // 用于延迟关闭下拉框
    };
  },
  computed: {
    filteredWebLinks() {
      if (!this.searchQuery && !this.isTyping) {
        return this.defaultwebLinks;
      }
      const query = this.searchQuery.toLowerCase();
      return this.defaultwebLinks.filter(link => {
        return (
          (link.name && link.name.toLowerCase().includes(query)) ||
          (link.url && link.url.toLowerCase().includes(query))
        );
      });
    },
  },
  methods: {
    handleMouseOver() {
      this.showWebLinksDropdown = true;
      clearTimeout(this.dropdownTimeout); // 清除延迟关闭
    },
    handleMouseLeave() {
      if (!this.isDropdownPinned && !this.isTyping && !this.isAddingLink) {
        this.showWebLinksDropdown = false;
      }
    },
    handleDropdownMouseLeave(event) {
      const dropdown = event.currentTarget;
      const relatedTarget = event.relatedTarget;
      if (!dropdown.contains(relatedTarget) && !this.isTyping && !this.isAddingLink) {
        this.handleMouseLeave();
      }
    },
    async addWebLink() {
      this.isAddingLink = true; // 标记正在添加链接
      try {
        this.newlink.icon = `${this.newlink.url}/favicon.ico`;
        const response = await AddLinks({
          name: this.newlink.name,
          url: this.newlink.url,
          icon: this.newlink.icon,
        });
        if (response.success) {
          this.defaultwebLinks.push({ ...this.newlink });
          this.newlink = { name: '', url: '', icon: '' };
        }
        this.showModal = true;
        this.modalMessage = `${this.newlink.name}\n${this.newlink.url}\n${response.message}`;
        this.showAddLinkModal = false;
      } catch (error) {
        this.showModal = true;
        this.modalMessage = `${this.newlink.name}\n${this.newlink.url}\n添加网址失败`;
      } finally {
        // 添加完成后延迟 5 秒关闭下拉框
        this.dropdownTimeout = setTimeout(() => {
          this.isAddingLink = false;
          if (!this.isDropdownPinned && !this.isTyping) {
            this.showWebLinksDropdown = false;
          }
        }, 1000);
      }
    },
    async removeWebLink(index) {
      const url = this.defaultwebLinks[index].url;
      try {
        const response = await DelLinks({ url: url });
        if (response.success) {
          this.defaultwebLinks.splice(index, 1);
        }
        this.showModal = true;
        this.modalMessage = `${this.defaultwebLinks[index].name}\n${this.defaultwebLinks[index].url}\n${response.message}`;
      } catch (error) {
        this.showModal = true;
        this.modalMessage = `${this.defaultwebLinks[index].name}\n${this.defaultwebLinks[index].url}\n请检查网络`;
      }
    },
    togglePinDropdown() {
      this.isDropdownPinned = !this.isDropdownPinned;
      if (this.isDropdownPinned) {
        this.showWebLinksDropdown = true;
      }
    },
    handleSearchInput(query) {
      this.searchQuery = query;
      this.isTyping = true; // 搜索时强制保持下拉框显示
    },
    handleFocus() {
      this.isTyping = true; // 输入框聚焦时强制保持下拉框显示
    },
    handleBlur() {
      if(!this.searchQuery){
        this.searchQuery = '';
      }
      this.isTyping = false; // 输入框失焦时恢复正常逻辑
    },
    handleAddLinkClick() {
      this.isAddingLink = true; // 点击添加链接时强制保持下拉框显示
      this.showAddLinkModal = true;
    },
    handleIconError(event) {
      event.target.src = this.defaulticon;
    },
  },
};
</script>

<style scoped>
.web-linker-container {
  position: absolute; /* 为下拉列表提供定位参考 */
  width: 100%;
  z-index: 10px;
}

.web-linker-container li{
  margin:10px 0;
}

.web-linker-container a {
  color: #fdfdfd;
  text-decoration: none;
  display: block;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.web-linker-container a:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* 网页链接器下拉列表 */
.web-links-dropdown {
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 250px;
  position: absolute;
  top: 0%; /* 下拉列表紧贴按钮 */
  left: 120px;
  z-index: 10; /* 确保下拉列表在其他内容之上 */
  height: 350px;
  overflow: hidden;
}

.web-links-list {
  list-style: none;
  padding: 0;
  margin: 0;
  max-height: 250px;
  overflow-y: auto;
}

.web-links-list li {
  margin: 5px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.web-links-list a {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #333;
}

.web-icon {
  width: 20px;
  height: 20px;
  margin-right: 10px;
}

.web-info {
  display: flex;
  flex-direction: column;
}

.web-name {
  font-weight: bold;
}

.web-url {
  font-size: 0.8em;
  color: #666;
}

/* 添加和删除按钮样式 */
.add-button, .delete-button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8em;
}

.add-button {
  background-color: #28a745;
  color: white;
  margin-top: 20px;
  margin-bottom: 10px;
}

.delete-button {
  background-color: #dc3545;
  color: white;
}

  /* 钉子图标样式 */
  .pin-icon {
  position: absolute;
  top: 5px;
  left: 1px;
  cursor: pointer;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  }

  .pin-icon-content {
  width: 10px;
  height: 10px;
  background-color: #fff;
  border-radius: 2px;
  }

  .pin-icon:hover {
  background-color: #b72424;
  }
  
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    background: white;
    padding: 20px;
    border-radius: 8px;
    width: 300px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }

  .modal-content h3 {
    margin-top: 0;
    font-size: 1.2em;
    text-align: center;
  }

  .form-group {
    margin-bottom: 15px;
    display: flex;
    align-items: center;
    gap:5px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    font-weight: bold;
    color: #333;
  }

  .form-group input {
    width: 80%;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
  }

  .form-actions {
    display: flex;
    justify-content: space-between;
    margin-top: 20px;
  }

  .form-actions button {
    padding: 8px 15px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9em;
  }

  .add-button {
    background-color: #28a745;
    color: white;
  }

  .cancel-button {
    background-color: #dc3545;
    color: white;
  }

  .put{
    display: flex;
  }
</style>