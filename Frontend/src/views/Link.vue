<template>
  <div class="web-linker-container" 
    @mouseover="handleMouseOver" 
    @mouseleave="handleMouseLeave"
  >
  <div class="click-button">
    <a href="#" @mouseover="showWebLinksDropdown = true">
      <img src="@/assets/icon/Link.png" alt="Link" class="icon"/>
    </a>
  </div>
    <div v-if="showWebLinksDropdown || isDropdownPinned" class="web-links-dropdown" @mouseleave="handleDropdownMouseLeave">
      <div class="pin-icon" @click="togglePinDropdown" >
        <img v-if="isDropdownPinned === true" src="@/assets/icon/ding1.png" alt="Pin" class="ding_icon"/>
        <img v-else src="@/assets/icon/ding0.png" alt="Pin" class="ding_icon"/>
      </div>
      <div class="put">
        <SearchLink
          v-model="searchQuery"
          @input="handleSearchInput"
          @focus="handleFocus"
          @blur="handleBlur"
          placeholder="搜索网页链接"
        />
        <img src="@/assets/icon/add_tag.png" alt="Add" class="add-icon" @click="handleAddLinkClick"/>
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
          <img src="@/assets/icon/delete_tag.png" alt="Delete" class="delete-icon" @click="removeWebLink(index)"/>
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
            <button type="submit">
              <img src="@/assets/icon/comfirm.png" alt="Comfirm" class="comfirm-icon" />
            </button>
            <img src="@/assets/icon/cancel.png" alt="Cancel" class="cancel-icon" @click="CancelClick"/>
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
      defaultwebLinks: this.$store.state.links,//从vuex获取 this.$store.state.links;
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
          url_name: this.newlink.name,
          url: this.newlink.url,
          icon: this.newlink.icon,
        });
        if (response.success) {
          this.defaultwebLinks.push({ ...this.newlink });
          this.$store.commit('SET_LINKS',this.defaultwebLinks);  //把更改后的数组存入全局
          this.newlink = { name: '', url: '', icon: '' };
        }
        this.showModal = true;
        this.modalMessage = `${this.newlink.name}\n${this.newlink.url}\n添加网址成功`;
        this.isAddingLink = false;
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
          this.$store.commit('SET_LINKS',this.defaultwebLinks);
        }
        this.showModal = true;
        this.modalMessage = `${this.defaultwebLinks[index].name}\n${this.defaultwebLinks[index].url}\n$删除链接成功`;
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

    CancelClick(){
      this.showAddLinkModal = false;
      this.isAddingLink= false;
      this.newlink = { name: '', url: '', icon: '' };
    },
  },
};
</script>

<style scoped>

.click-button:hover{
  background-color: var(--sidebar-background-color1);
}

.web-linker-container {
  position: absolute; /* 为下拉列表提供定位参考 */
  width: 100%;
  z-index: 10px;
}

.web-linker-container li{
  margin:10px 0;
}

.web-linker-container a {
  color: var(--background-color);
  text-decoration: none;
  display: block;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}


/* 网页链接器下拉列表 */
.web-links-dropdown {
  background-color: var(--background-color);
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 4px var(--background-color2);
  width: 300px;
  position: absolute;
  top: 0%; /* 下拉列表紧贴按钮 */
  left: 60px;
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
  justify-content: space-around;
  align-items: center;
  background-color: var(--background-color);
  border-radius: 8px;
}

.web-links-list li:hover{
  background-color: var(--background-color2);
}

.web-links-list a {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--text-color);
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
  text-align: left;
}

.web-url {
  font-size: var(--font-size-small);
  color: #666;
}

.add-icon {
    width: 45px;
    height: 45px;
    margin-top: 5px;
  }

.delete-icon{
  width: 40px;
  height: 40px;
}

.add-icon:hover, .delete-icon:hover{
  background-color: var(--button-background-color);
  cursor: pointer;
  border-radius: 30px;
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
    background: var(--background-color);
    padding: 10px;
    border-radius: 8px;
    width: 300px;
    box-shadow: 0 2px 10px var(--background-color2);
  }

  .modal-content h3 {
    margin-top: 10px;
    font-size: var(--font-size-mlarge);
    color: var(--text-color);
    text-align: center;
    margin-bottom: 15px;
  }

  .form-group {
    margin-bottom: 15px;
    display: flex;
    align-items: center;
    gap:5px;
  }

  .form-group label {
    width: 20%;
    display: block;
    margin-bottom: 5px;
    color: var(--text-color);
  }

  .form-group input {
    width: 80%;
    padding: 6px;
    border-radius: 4px;
  }

  .form-actions {
    display: flex;
    justify-content: space-around;
    margin-top: 20px;
  }

  .form-actions button {
    border: none;
    background-color: var(--background-color);
  }

  .comfirm-icon {
    width: 40px;
    height: 40px;
  }

  .cancel-icon {
    width: 40px;
    height: 40px;
  }

  .comfirm-icon:hover, .cancel-icon:hover {
    background-color: var(--button-background-color);
    border-radius: 16px;
    cursor: pointer;
  }

  .put{
    display: flex;
    flex-direction: row;
    margin-left: 20px;
    align-items: center;
  }
  
  .pin-icon{
    width: 25px;
    margin-bottom: px;
    position:absolute;
  }

  .ding_icon{
    width: 25px;
    height: 25px;
  }

  a .icon {
  width: 25px; /* 图标大小 */
  height: 25px; /* 图标大小 */
  margin-right: 5px; /* 图标和文字之间的间距 */
}

</style>