const WebSocket = require("ws");
const { setupWSConnection } = require("y-websocket/bin/utils");
const Y = require("yjs");
const axios = require("axios"); // 向 Go 后端发请求

const port = 1234;
const wss = new WebSocket.Server({ port });

// 保存每个 room 的文档和计数器
const docs = new Map();

wss.on("connection", (ws, req) => {
  const url = req.url || "";
  const [, roomName] = url.split("?room=");
  const docName = roomName || "default-room";

  // 获取/创建该房间对应的 Y.Doc
  let doc = docs.get(docName);
  if (!doc) {
    doc = new Y.Doc();
    docs.set(docName, doc);
    setupAutoSync(doc, docName);
  }

  setupWSConnection(ws, req, { docName, doc });
});

console.log(`Yjs WebSocket Server running on ws://localhost:${port}`);

function setupAutoSync(doc, roomName) {
  let updateBuffer = [];
  let updateCount = 0;

  // 每当文档有更新时
  doc.on('update', update => {
    updateBuffer.push(update);
    updateCount++;

    // 每10次操作就上报一次增量
    if (updateCount >= 10) {
      sendUpdatesToGo(roomName, updateBuffer);
      updateBuffer = [];
      updateCount = 0;
    }
  });

  // 每30秒保存一次 snapshot
  setInterval(() => {
    const snapshot = Y.encodeStateAsUpdate(doc);
    sendSnapshotToGo(roomName, snapshot);
  }, 30000);
}

async function sendUpdatesToGo(roomName, updates) {
  const merged = Y.mergeUpdates(updates);
  const base64 = Buffer.from(merged).toString('base64');
  try {
    await axios.post("http://localhost:8080/api/docs/updates", {
      doc_id: roomName,
      update: base64,
    });
    console.log(`Sent updates for room ${roomName}`);
  } catch (err) {
    console.error("Failed to send updates:", err.message);
  }
}

async function sendSnapshotToGo(roomName, snapshot) {
  const base64 = Buffer.from(snapshot).toString('base64');
  try {
    await axios.post("http://localhost:8080/api/docs/snapshot", {
      doc_id: roomName,
      snapshot: base64,
    });
    console.log(`Saved snapshot for room ${roomName}`);
  } catch (err) {
    console.error("Failed to send snapshot:", err.message);
  }
}

