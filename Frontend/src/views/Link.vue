<template>
    <div class="web-linker-container" @mouseover="showWebLinksDropdown = true" @mouseleave="handleMouseLeave">
      <a href="#" @mouseover="showWebLinksDropdown = true">网页链接器</a>
      <div v-if="showWebLinksDropdown || isDropdownPinned" class="web-links-dropdown" @mouseleave="handleDropdownMouseLeave">
        <!-- 使用 searchBar 组件 -->
        <SearchBar v-model="searchQuery" @input="searchWebLink" placeholder="搜索网页链接" />
        <ul class="web-links-list">
          <li v-for="(link, index) in filteredWebLinks" :key="index">
            <a :href="link.url" target="_blank">
              <img :src="link.icon" alt="Icon" class="web-icon" />
              <div class="web-info">
                <span class="web-name">{{ link.name }}</span>
                <span class="web-url">{{ link.url }}</span>
              </div>
            </a>
            <button @click="removeWebLink(index)" class="delete-button">删除</button>
          </li>
        </ul>
        <button @click="showAddLinkModal = true" class="add-button">添加链接</button>
        <!-- 钉子图标 -->
        <div class="pin-icon" @click="togglePinDropdown">
          <i class="fas fa-thumbtack"></i>
        </div>
      </div>
      <div v-if="showAddLinkModal" class="modal-overlay">
        <div class="modal-content">
            <h3>添加新链接</h3>
            <form @submit.prevent="addWebLink">
            <label for="linkName">名称:</label>
            <input type="text" id="linkName" v-model="newlink.name" required />
            <label for="linkUrl">URL:</label>
            <input type="url" id="linkUrl" v-model="newlink.url" required />
            <button type="submit" class="add-button">添加</button>
            <button @click="showAddLinkModal = false" class="cancel-button">取消</button>
            </form>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import SearchBar from '@/components/base/SearchBar.vue'; // 确保导入 searchBar 组件
  
  export default {
    name: 'WebLinker',
    components: {
      SearchBar, // 注册 searchBar 组件
    },
    data() {
      return {
        showWebLinksDropdown: false,
        isDropdownPinned: false, // 控制下拉列表是否固定
        searchQuery: '',
        webLinks: [
          { name: 'Google', url: 'https://www.google.com', icon: 'https://www.google.com/favicon.ico' },
          { name: 'Github', url: 'https://github.com', icon: 'https://github.githubassets.com/favicon.ico' },
          { name: 'Facebook', url: 'https://www.facebook.com', icon: 'https://www.facebook.com/favicon.ico' },
          { name: 'Vue', url: 'https://vuejs.org', icon: 'https://vuejs.org/favicon.ico' },
        ],
        showAddLinkModal: false, // 控制弹窗的显示
        newlink:{
            name:'',
            url:'',
            icon:'',
        },
        defaulticon:'https://vuejs.org/favicon.ico',
      };
    },
    computed: {
      filteredWebLinks() {
        return this.webLinks.filter(link => {
          return link.name && link.name.toLowerCase().includes(this.searchQuery.toLowerCase());
        });
      },
    },
    methods: {
      handleMouseLeave() {
        if (!this.isDropdownPinned) {
          this.showWebLinksDropdown = false;
        }
      },
      handleDropdownMouseLeave(event) {
        // 如果鼠标在列表中操作，不隐藏下拉列表
        const dropdown = event.currentTarget;
        if (!dropdown.contains(event.relatedTarget)) {
            this.handleMouseLeave();
        }
      },
      async addWebLink() {
        const iconUrl = await this.getFavicon(this.newlink.url);
        this.newlink.icon = iconUrl;

        // 添加新的网页链接
        this.webLinks.push({ ...this.newlink });

        // 重置表单
        this.newLink = { name: '', url: '', icon: '' };
        this.showAddLinkModal = false;
      },
      removeWebLink(index) {
        // 删除指定索引的网页链接
        this.webLinks.splice(index, 1);
      },
      togglePinDropdown() {
        // 切换下拉列表的固定状态
        this.isDropdownPinned = !this.isDropdownPinned;
        if (this.isDropdownPinned) {
          this.showWebLinksDropdown = true; // 固定时确保下拉列表显示
        }
      },
      searchWebLink() {
        // 搜索逻辑在 computed 属性中处理
      },
      async getFavicon(url) {
        try {
            // 获取 URL 的主机名
            const hostname = new URL(url).hostname;

            // 构造 favicon 的 URL
            const faviconUrl = `https://${hostname}/favicon.ico`;

            // 尝试获取图标
            const response = await fetch(faviconUrl);

            // 如果图标存在，返回图标的 URL
            if (response.ok) {
            return faviconUrl;
            } else {
            // 如果图标不存在，返回默认图标
            return this.defaultIcon;
            }
        } catch (error) {
            console.error('Error fetching favicon:', error);
            return this.defaultIcon; // 默认图标 URL
        }
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
    width: 300px;
    position: absolute;
    top: 0%; /* 下拉列表紧贴按钮 */
    left: 120px;
    z-index: 10; /* 确保下拉列表在其他内容之上 */
  }
  
  .web-links-list {
    list-style: none;
    padding: 0;
    margin: 0;
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
    margin-top: 10px;
  }
  
  .delete-button {
    background-color: #dc3545;
    color: white;
  }
  
    /* 钉子图标样式 */
    .pin-icon {
    position: absolute;
    top: 5px;
    right: 5px;
    cursor: pointer;
    width: 20px;
    height: 20px;
    background-color: #706a6a;
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
    }

    .modal-content form {
    display: flex;
    flex-direction: column;
    }

    .modal-content label {
        margin-bottom: 5px;
        display: block; /* 确保 label 是块级元素 */
        font-weight: bold; /* 可选：加粗字体 */
    }

    .modal-content input {
    margin-bottom: 10px;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
    }

    .modal-content button {
    padding: 8px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    }

    .modal-content .add-button {
    background-color: #28a745;
    color: white;
    margin-top: 10px;
    }

    .modal-content .cancel-button {
    background-color: #dc3545;
    color: white;
    margin-top: 10px;
    }
  </style>