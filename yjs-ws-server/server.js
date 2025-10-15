// server.js
const WebSocket = require("ws");
const { setupWSConnection } = require("y-websocket/bin/utils");

const port = 1234;
const wss = new WebSocket.Server({ port });

wss.on("connection", (ws, req) => {
  // 可选：解析房间名称（默认为 "default-room"）
  const url = req.url || "";
  const [, roomName] = url.split("?room=");  
  setupWSConnection(ws, req, { docName: roomName || "default-room" });
});

console.log(`WebSocket server running on ws://localhost:${port}`);
