<template>
  <div class="container">

    <!-- 显示文件名、在线用户等信息 -->
    <div class="doc-info">
      <div class="left-infos">
        <img src="@/assets/icon/return.png" alt="返回图标" class="icon" @click="returnToWorkspace"/>
        <button class="ql-undo" :disabled="!canUndo" @click="handleUndo">↶</button>
        <button class="ql-redo" :disabled="!canRedo" @click="handleRedo">↷</button>
      </div>
      <div class="doc_name">{{ currentDoc.doc_name }}</div>
      <!-- <button @click="saveSnapshot"> 保存快照</button>
      <button @click="loadSnapshot"> 获取快照</button> -->
      <div class="user-info">
        <div class="other-user">正在编辑：
          <div class = "other-user"
            v-for="[key, user] in remoteUsers"
            :key="key"
            :style="{ color: user.color }"
          >
            <!-- 用户头像 -->
            <img
              :src="user.userIcon"
              :alt="user.name"
              style="width: 18px; height: 18px; border-radius: 50%;"
            />
            {{ user.name }}  <!-- 显示用户名 -->
          </div>
        </div>
        <div class="me-user">
            <img
                :src="currentUser.avatar"
                :alt="currentUser.username"
                style="width: 18px; height: 18px; border-radius: 50%; margin-top:3px; margin-right:5px;"
            />
            <span>{{ currentUser.username }}</span>
        </div>
      </div>
    </div>
    <!--  Quill 编辑器 -->
    <div ref="quillEditor" class="quill-editor"></div>

  </div>
</template>
<script>
import * as Y from "yjs";  // 导入 Yjs 核心库
import { WebsocketProvider } from "y-websocket";  // 导入 Yjs WebSocket 提供者
import { UndoManager } from 'yjs';
import Quill from "quill";
import { QuillBinding} from "y-quill";
import QuillCursors from "quill-cursors";
import 'quill/dist/quill.snow.css';
import * as WorkSpaceAPI from '@/services/workspace_api';
import { Buffer } from "buffer";

if (!Quill.imports['modules/cursors']) {
  Quill.register("modules/cursors", QuillCursors);
}

export default {
  name: "CollaborativeInput",  // 组件名称
  computed: {
    currentDoc() {
        return this.$store.getters.getCurrentDoc;  // 从 Vuex 获取 currentDoc
        // 格式如下：    
        // currentDoc: {
        //   doc_id: null,
        //   doc_name: '',
        // }
    },
    currentUser() {
      return this.$store.state.user;
    }
  },
  data() {
    return {
      // currentUser: { name: "", color: "", userIcon: "" },  // 当前用户信息
      remoteUsers: new Map(),  // 远程用户信息 Map，key: clientID, value: user对象
      canUndo: false,
      canRedo: false,
      snapshotTimer: null,
    };
  },
  created() {
    // 从Go后端获取快照
    this.loadSnapshot();
    
    // 创建 Yjs 文档
    this.ydoc = new Y.Doc();

    // 创建 WebSocket provider，连接到本地服务器文档,唯一文档名是doc_id
    // this.provider = new WebsocketProvider("ws://localhost:1234", this.currentDoc.doc_id, this.ydoc);
    // console.log(this.currentDoc.doc_id);
    this.provider = new WebsocketProvider(
        `ws://localhost:1234?room=${this.currentDoc.doc_id}`,
        this.currentDoc.doc_id,  // 这个参数可以随便，但建议保持同一个
        this.ydoc
    );
    // 随机生成用户颜色
    const randomColor = () => `rgb(${Math.random()*255|0},${Math.random()*255|0},${Math.random()*255|0})`;

    // 随机生成用户名
    // const names = ["Alice", "Bob", "Charlie", "David", "Eve"];
    // const randomName = () => names[Math.floor(Math.random()*names.length)] + Math.floor(Math.random()*100);

    // const userName = randomName();

    // 设置本地用户状态，Yjs awareness 用于广播用户信息（姓名、颜色、头像）
    this.provider.awareness.setLocalStateField("user", {
      id: this.currentUser.id,
      name: this.currentUser.username,
      color: randomColor(),
      userIcon: this.currentUser.avatar,
    });
    
  },
  mounted() {
    // 获取共享文本对象
    const yText = this.ydoc.getText("shared-text");
    this.yTextRef = yText;

    const toolbarOptions = [
      ['bold', 'italic', 'underline', 'strike'],        // toggled buttons
      ['blockquote', 'code-block'],
      ['link', 'image', 'video', 'formula'],

      [{ 'header': 1 }, { 'header': 2 }],               // custom button values
      [{ 'list': 'ordered'}, { 'list': 'bullet' }, { 'list': 'check' }],
      [{ 'script': 'sub'}, { 'script': 'super' }],      // superscript/subscript
      [{ 'indent': '-1'}, { 'indent': '+1' }],          // outdent/indent
      [{ 'direction': 'rtl' }],                         // text direction

      [{ 'size': ['small', false, 'large', 'huge'] }],  // custom dropdown
      [{ 'header': [1, 2, 3, 4, 5, 6, false] }],

      [{ 'color': [] }, { 'background': [] }],          // dropdown with defaults from theme
      [{ 'font': [] }],
      [{ 'align': [] }],

      ['clean']                                         // remove formatting button
    ];

    // 获取可用视口的高度，用来在style里设置容器高度
    const vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty('--vh', `${vh}px`);

    // quill初始化
    this.quill = new Quill(this.$refs.quillEditor, {
        theme: "snow",
        modules: {
        toolbar: toolbarOptions,
        cursors: true, // 保留，但不自己控制
        },
    });

    // // 加两个按钮在toolbar里
    // this.$nextTick(() => {
    //     const toolbarEl = this.$el.querySelector('.ql-toolbar');
    //     if (toolbarEl) {
    //         // 创建撤销按钮
    //         this.undoBtn = document.createElement('button');
    //         this.undoBtn.className = 'ql-undo';
    //         this.undoBtn.innerHTML = '↶';
    //         this.undoBtn.onclick = () => this.handleUndo();
    //         this.undoBtn.disabled = !this.canUndo;  // 初始状态

    //         // 创建恢复按钮
    //         this.redoBtn = document.createElement('button');
    //         this.redoBtn.className = 'ql-redo';
    //         this.redoBtn.innerHTML = '↷';
    //         this.redoBtn.onclick = () => this.handleRedo();
    //         this.redoBtn.disabled = !this.canRedo;  // 初始状态

    //         toolbarEl.appendChild(this.undoBtn);
    //         toolbarEl.appendChild(this.redoBtn);
    //     }
    // });

    // 只用 QuillBinding 自动处理内容 + 光标同步
    this.binding = new QuillBinding(yText, this.quill, this.provider.awareness);

    // 撤销
    this.undoManager = new UndoManager(yText, {
        captureTimeout: 500,
        ignoreRemoteMapChanges: true,
        trackedOrigins: new Set([this.binding]),
    });
    this.ydoc.on("update", this.updateUndoRedoState);

    // 撤销/重做快捷键
    this.quill.keyboard.addBinding({ key: 'z', shortKey: true }, () => this.handleUndo());
    this.quill.keyboard.addBinding({ key: 'y', shortKey: true }, () => this.handleRedo());

    // 删除了手动 selection-change 和 cursor 广播部分

    // 使用 IndexedDB 持久化文档数据
    //this.persistence = new IndexeddbPersistence("collab-input-db", this.ydoc);

    // 获取本地用户信息
    // this.currentUser = this.provider.awareness.getLocalState()?.user || {};

    // 监听远程用户状态变化
    this.provider.awareness.on("change", () => {
        const states = this.provider.awareness.getStates();
        const remote_users = new Map();
        states.forEach((state, clientID)=>{
            const user = state.user;
            if(user.id!==this.currentUser.id && !remote_users.has(user.id))
            {
                remote_users.set(user.id, user);
            }
        })
        this.remoteUsers = remote_users;
    });

    // 每隔30秒，保存文档快照到数据库
    this.snapshotTimer = setInterval(()=>{
        this.saveSnapshot();
    }, 30000);

    window.addEventListener("beforeunload", this.destroy_handler);
  },

  beforeUnmount(){
    window.removeEventListener("beforeunload", this.destroy_handler);
    this.destroy_handler();
    console.log("进入了beforeUnmount");
  },

  // beforeRouteLeave() {
  //   window.removeEventListener("beforeunload", this.destroy_handler);
  //   this.destroy_handler();
  //   console.log("进入了beforeRouteLeave");
  // },

  methods: {
    // 让Go后端保存快照
    async saveSnapshot() {
        const update = Y.encodeStateAsUpdate(this.ydoc); // 用Y.encodeStateAsUpdate而不是Y.encodeSnapshot
        const update_base64 = Buffer.from(update).toString('base64');
        try{
            const response = await WorkSpaceAPI.saveSnapshot(this.currentDoc.doc_id, update_base64);
            if(response.status!=200){
                console.log("保存协作文档快照失败");
            }
        }
        catch(err){
            console.log("保存快照出错：");
            console.log(err);
        }
    },

    // 从Go后端获取快照
    async loadSnapshot() {
        const response = await WorkSpaceAPI.getSnapshot(this.currentDoc.doc_id);
        const data = response.data;
        if (data.snapshot) {
            const update = Uint8Array.from(atob(data.snapshot), c => c.charCodeAt(0));
            Y.applyUpdate(this.ydoc, update);
            console.log('Snapshot restored from server')
        } else {
            console.log('ℹ️ No snapshot found, starting fresh')
        }
    },

    returnToWorkspace()
    {
      // 跳转到编辑页面
      this.$router.push(`/workspace/collabdocs`);
    },

    destroy_handler()
    {
        // 路由变换/组件销毁/关闭标签页前，移除监听并断开连接
        if (this.provider?.awareness) {
            // 清空本地状态，通知其他客户端该用户已离开
            this.provider.awareness.setLocalState(null);
        }
        this.provider?.disconnect();

        // 清除 undoManager
        if(this.undoManager)
        {
        this.undoManager.destroy();
        this.undoManager = null;
        }
        if(this.binding)
        {
            this.binding.destroy();
        }
        // 清除定时器
        if (this.snapshotTimer) {
            clearInterval(this.snapshotTimer);
            this.snapshotTimer = null;
        }
    },

    updateUndoRedoState()
    {
      this.canUndo = this.undoManager.undoStack.length > 0;
      this.canRedo = this.undoManager.redoStack.length > 0;
      if (this.undoBtn) this.undoBtn.disabled = !this.canUndo;
      if (this.redoBtn) this.redoBtn.disabled = !this.canRedo;
    },

    // 撤销
    handleUndo(event)
    {
      // 快捷键的Ctrl Z事件，需要阻止浏览器的撤销文本行为
      if(event)
      {
        event.preventDefault();
        event.stopPropagation();
      }
      if(this.canUndo)
      {
        this.undoManager.undo();
      }
    },

    // 恢复
    handleRedo(event)
    {
      // 快捷键的Ctrl Y事件，需要阻止浏览器的行为
      if(event)
      {
        event.preventDefault();
        event.stopPropagation();
      }
      if(this.canRedo)
      {
        this.undoManager.redo();
      }
    }
  },
};
</script>
<style scoped>
  html, body{
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
    overflow: hidden;   /* 防止全局滚动 */
    box-sizing: border-box;
  }
  *{
    box-sizing: inherit;
  }
  .container{
    display: flex;
    flex-direction: column;
    height: calc(var(--vh, 1vh) * 100 - 40px);
    width: 100%;
    padding: 20px;
  }
  .doc-info{
    display: flex;
    flex-direction: row;
    flex-shrink:0;
    align-items: center;
    justify-content: space-between;
    background-color:lightgrey;
  }
  .left-infos{
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
    align-items: center;
  }
  .ql-undo, .ql-redo {
    font-size: 18px;
    padding: 4px 8px;
    border: none;
    background: none;
    color:var(--button-text-color);
    cursor: pointer;
  }
  .ql-undo:disabled, .ql-redo:disabled {
    color:gray;
  }
  .doc_name{
    font-size: 20px;
    margin-left: 20px;
  }
  .user-info {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    align-items: center;
  }
  .other-user, .me-user{
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    align-items: center;
    padding: 8px;
    margin-right: 5px;
    font-size: 16px;
  }
  .quill-editor{
    flex: 1;
    min-height: 0;
    padding: 20px;
  }
  .icon{
    width: 25px;
    height: 25px;
    margin-left: 5px;
    cursor: pointer;
  }

</style>
