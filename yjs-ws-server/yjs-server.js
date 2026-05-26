const WebSocket = require("ws");
const { setupWSConnection } = require("y-websocket/bin/utils");
const Y = require("yjs");
const axios = require("axios");

const port = 1234;
const wss = new WebSocket.Server({ port });

const docs = new Map();
const SNAPSHOT_INTERVAL = 30000; // 30s

wss.on("connection", (ws, req) => {
  const url = req.url || "";
  const [, roomName] = url.split("?room=");
  const docName = roomName || "default-room";

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
  // Periodically save snapshot to Go backend
  const timer = setInterval(() => {
    const snapshot = Y.encodeStateAsUpdate(doc);
    sendSnapshotToGo(roomName, snapshot);
  }, SNAPSHOT_INTERVAL);

  // Clean up timer when doc is no longer needed
  doc.on('destroy', () => clearInterval(timer));
}

async function sendSnapshotToGo(roomName, snapshot) {
  const base64 = Buffer.from(snapshot).toString('base64');
  try {
    await axios.post("http://localhost:8080/workspace/collabdocs/savesnapshot", {
      doc_id: roomName,
      snapshot: base64,
    });
    console.log(`Saved snapshot for room ${roomName}`);
  } catch (err) {
    // Frontend also saves snapshots directly, so this is non-critical.
    // 404 means the Go backend snapshot endpoint needs the correct route.
    if (err.response && err.response.status !== 404) {
      console.error("Failed to send snapshot:", err.message);
    }
  }
}

