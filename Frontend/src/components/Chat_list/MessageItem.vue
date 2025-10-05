<template>
  <div class="message-item" >
    <div v-if="this.$store.state.user.id !== message.send_account_id" class="friend-message">
      <div class="avatar" @contextmenu.prevent="showBanned($event)">
        <img :src="message.avatar" alt="avatar" @click="showProfileCard($event)"/>
      </div>
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ formatTime(message.create_time) }}</span>
        </div>
        <!--文本消息-->
        <div class="message-content" 
            v-if="message.type==='text'"
            v-html="message.content" 
            @contextmenu.prevent="showContextMenu($event, message)">
        </div>
        <!--图片消息-->
        <div v-else-if="message.type==='image'" @contextmenu.prevent="showContextMenu($event, message)">
          <img :src="message.content" alt="image" style="max-width: 100%; max-height: 200px;"/>
        </div>
        <!--协作文档消息-->
        <div class="message-collab-doc" 
            v-else-if="message.type==='collab_doc'"
            @contextmenu.prevent="showContextMenu($event, message)">
          <div class="doc-card" @click="editCollabDoc">
            <div class="doc-icon">
              <img src="@/assets/icon/text.png" alt="文档" />
            </div>
            <div class="doc-info">
              <div class="doc-name">{{ JSON.parse(message.content).doc_name }}</div>
            </div>
          </div>
        </div>
        <!--文件消息-->
        <div class="message-file" v-else-if="message.type==='file'" @contextmenu.prevent="showContextMenu($event, message)">
          <div class="file-item">
            <img src="@/assets/images/default-file.png" alt="file"/>
            <div class="file-header">
              <div class="file-name">{{ message.content.name }}</div>
              <span class="file-size">{{ message.content.size }}</span>
            </div>
          </div>
          <span class="file-buttons">
            <button class="file-button" @click="downloadFile">下载</button>
            <!-- <button class="file-button" @click="previewFile">预览</button> -->
            <!-- <a ref="link" style="visibility: hidden" :href="message.content" download>下载</a> -->
          </span>
        </div>
        <!--代码消息-->
        <div v-else class="editor-container" @contextmenu.prevent="showContextMenu($event, message)">
          <div ref="editor" class="editor"></div>
        </div>
      </div>
    </div>

    <!--我的消息-->
    <div v-else class="my-message">
      <div class="message-content-wrapper">
        <div class="message-header">
          <span class="message-sender">{{ message.sender }}</span>
          <span class="message-time">{{ formatTime(message.create_time) }}</span>
        </div>
        <!--文本消息-->
        <div class="message-content" 
             v-if="message.type==='text'"
             v-html="message.content" 
             @contextmenu.prevent="showContextMenu($event, message)">
        </div>
        <!--图片消息-->
        <div v-else-if="message.type==='image'" @contextmenu.prevent="showContextMenu($event, message)">
          <img :src="message.content" alt="image" style="max-width: 100%; max-height: 200px;"/>
        </div>
        <!--协作文档消息-->
        <div class="message-collab-doc" 
            v-else-if="message.type==='collab_doc'"
            @contextmenu.prevent="showContextMenu($event, message)">
          <div class="doc-card" @click="editCollabDoc">
            <div class="doc-icon">
              <img src="@/assets/icon/text.png" alt="文档" />
            </div>
            <div class="doc-info">
              <div class="doc-name">{{ JSON.parse(message.content).doc_name }}</div>
            </div>
          </div>
        </div>
        <!--文件消息-->
        <div class="message-file" v-else-if="message.type==='file'" @contextmenu.prevent="showContextMenu($event, message)">
          <div class="file-item">
            <img src="@/assets/images/default-file.png" alt="file"/>
            <div class="file-header">
              <div class="file-name">{{ message.content.name }}</div>
                <span class="file-size">{{ formatFileSize(message.content.size) }}</span>
            </div>
          </div>
          <span class="file-buttons">
            <button class="file-button" @click="downloadFile">下载</button>
            <!-- <button class="file-button" @click="previewFile">预览</button> -->
            <!-- <a ref="link" style="visibility: hidden" :href="message.content" download>下载</a> -->
          </span>
        </div>
        <!--代码消息-->
        <div v-else class="editor-container" @contextmenu.prevent="showContextMenu($event, message)">
          <div ref="editor" class="editor"></div>
        </div>
      </div>
      <div class="avatar" @contextmenu.prevent="showBanned($event)">
        <img :src="message.avatar" alt="avatar" @click="showProfileCard($event)"/>
      </div>
    </div>

    <!--预览文件-->

  </div>
</template>

<script>
import * as monaco from 'monaco-editor';
import * as WorkSpaceAPI from '@/services/workspace_api';
export default {
  props: ['message'],
  data() {
    return {
      showMenu: false,
      axis: {
        x: 0,
        y: 0
      },
    };
  },
  methods: {
    downloadFile(){
      const blob = new Blob([this.message.content], { type: 'application/octet-stream' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', this.message.content.name);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url); // 释放 URL 对象
    },
    async editCollabDoc()
    {
      // this.$router.push(`/workspace/doc/${id}`);
      try{
        // 先加入文档编辑成员
        const response = await WorkSpaceAPI.joinMember(JSON.parse(this.message.content).doc_id);
        if(response.status === 200)
        {
          // 使用 Vuex 更新 currentDoc 对象
          this.$store.dispatch('updateCurrentDoc', {
            doc_id: JSON.parse(this.message.content).doc_id,
            doc_name: JSON.parse(this.message.content).doc_name,
          });

          // 跳转到编辑页面
          this.$router.push(`/workspace/collabdocs/editor`);
        }
      }catch(err){
        console.error('无法加入文档编辑:', err);
        alert('加入文档编辑失败');
      }
    },
    showContextMenu(event, message) {
      this.$emit('show-context-menu',event, message);
    },
    showProfileCard(event){
      this.$emit('show-profile-card', event, this.message.send_account_id);
    },
    // 格式化时间
    formatTime(time) {
      if(!time) return '';
      const now = new Date();
      const messageTime = new Date(time);
      const isToday = now.toDateString() === messageTime.toDateString();
      const isYesterday = new Date(now.setDate(now.getDate() - 1)).toDateString() === messageTime.toDateString();
      const isThisYear = now.getFullYear() === messageTime.getFullYear();

      if (isToday) {
        return messageTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else if (isYesterday) {
        return '昨天' + messageTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else if(isThisYear){
        return messageTime.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' }) + ' ' + messageTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      }else{
        return messageTime.toLocaleDateString();
      }
    },
    formatFileSize(size) {
      if (size < 1024) {
        return size + 'B';
      } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + 'KB';
      } else if (size < 1024 * 1024 * 1024) {
        return (size / 1024 / 1024).toFixed(2) + 'MB';
      } else {
        return (size / 1024 / 1024 / 1024).toFixed(2) + 'GB';
      }
    },
    showBanned(event){
      // this.$emit('show-banned', event, this.message.send_account_id);
    },
    
  },

  mounted() {
    if(this.message.type=== 'text'||this.message.type=== 'file'||this.message.type==='image' ||this.message.type==='collab_doc'){
      return;
    }
    this.$nextTick(()=>{
      if (!this.$refs.editor) {
        console.log("message type:", this.message.type);
        console.warn('Editor DOM 未渲染，无法初始化 monaco');
        return;
      }
      this.editor = monaco.editor.create(this.$refs.editor, {
        value: this.message.content,
        language: this.message.type,
        automaticLayout: true,
        readOnly: true,
        lineNumbersMinChars: 2, // 设置行号的最小字符数
        tabSize: 2, // 设置制表符宽度
        minimap: {
          enabled: false, // 禁用右侧的迷你地图
        },
        fontSize: 14, // 设置字体大小
        lineHeight: 20, // 设置行高
        padding: {
          top: 10,
          bottom: 10,
        },
      });

      this.editor.onDidChangeModelContent(() => {
        this.$emit('input', this.editor.getValue());
      });
    });
  },
  // watch: {
  //   language(newLang) {
  //     monaco.editor.setModelLanguage(this.editor.getModel(), newLang);
  //   },
  //   value(newValue) {
  //     if (newValue !== this.editor.getValue()) {
  //       this.editor.setValue(newValue);
  //     }
  //   },
  // },
  beforeDestroy() {
    if (this.editor) {
      this.editor.dispose();
    }
  },
};
</script>

<style scoped>
.message-item {
  display: flex;
  padding: 5px;
  position: relative;
  width: 100%;
}
.friend-message {
  align-self: flex-start; 
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  width: 100%;
}
.my-message {
  align-self: flex-end;
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  width: 100%;
}
.friend-message .avatar {
  align-self: flex-start;
}
.my-message .avatar {
  align-self: flex-start;
}
.avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.message-content-wrapper {
  max-width: 450px;
  display: inline-flex;
  flex-direction: column;
}
.message-header {
  flex:1;
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}
.message-sender {
  color: #888;
  font-size: var(--font-size-small-small);
  text-align: left;
  padding: 0 0 0 10px;
}
.message-time {
  color: #888;
  font-size: var(--font-size-small-small);
  text-align: right;
  padding: 0 0 0 10px;
}


.message-content {
  flex:5;
  background-color: var(--sidebar-background-color);
  color: var(--sidebar-text-color);
  padding: 10px;
  border-radius: 5px;
  text-align: left;
}
.message-file{
  flex:5;
  background-color: var(--sidebar-background-color);
  color:var(--sidebar-text-color);
  padding: 5px 10px 2px 2px;
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
}
.file-item{
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  flex-direction: row;
  padding: 3px 0 3px 0;
}
.message-file img{
  width: 50px;
  height: 50px;
}
.file-header{
  display: flex;
  flex-direction: column;
}
.file-name{
  margin-top: 5px;
  font-size: var(--font-size-small);
  color: #888;
}
.file-size{
  margin-top: 5px;
  font-size: var(--font-size-small);
  color: #888;
}
.file-buttons{
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  width: 100%;
  padding: 5px;
}
.file-button{
  padding: 2px 5px 2px 5px;
  border-radius: 10%;
  border: none;
  cursor: pointer;
}

.editor-container {
  width: 350px;
  max-height: 400px;
  padding: 3px;
}
.editor {
  width: 100%;
  height: 200px;
  text-align: left;
}
.message-image {
  max-width: 70%;
}

.img {
  max-width: 100%;
  border-radius: 5px;
}

/* 协作文档消息样式 */
.message-collab-doc {
  flex: 5;
  /*background-color: var(--sidebar-background-color);*/
  color: var(--sidebar-text-color);
  padding: 12px;
  border-radius: 8px;
}

.doc-card {
  display: flex;
  flex-direction: row;
  justify-content: space-evenly;
  align-items:center;
  background-color: white;
  gap: 12px;
  padding-right:12px;
  border-radius: 5px;
  cursor: pointer;
}

.doc-icon {
  flex-shrink: 0;
}

.doc-icon img {
  width: 35px;
  height: 35px;
  border-radius: 6px;
  margin: 5px;
}

.doc-info {
  flex: 1;
  min-width: 0;
}

.doc-name {
  font-size: var(--font-size-small);
  font-weight: 500;
  color: var(--sidebar-text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>