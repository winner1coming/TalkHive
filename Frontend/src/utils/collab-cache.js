// IndexedDB cache for collaborative document snapshots.
// Enables instant load from local cache (stale-while-revalidate),
// eliminating redundant network download and gzip decompression.

const DB_NAME = 'TalkHive_CollabCache'
const DB_VERSION = 1
const STORE_NAME = 'snapshots'

function openDB() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION)
    request.onupgradeneeded = (e) => {
      const db = e.target.result
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME, { keyPath: 'docId' })
      }
    }
    request.onsuccess = (e) => resolve(e.target.result)
    request.onerror = (e) => reject(e.target.error)
  })
}

// Retrieve cached snapshot for a doc. Returns { docId, snapshot, updatedAt } or null.
export async function getCachedSnapshot(docId) {
  let db
  try {
    db = await openDB()
    const tx = db.transaction(STORE_NAME, 'readonly')
    const store = tx.objectStore(STORE_NAME)
    return await new Promise((resolve, reject) => {
      const req = store.get(docId)
      req.onsuccess = () => resolve(req.result || null)
      req.onerror = () => reject(req.error)
    })
  } finally {
    db?.close()
  }
}

// Store snapshot in IndexedDB cache.
export async function setCachedSnapshot(docId, snapshotBase64) {
  let db
  try {
    db = await openDB()
    const tx = db.transaction(STORE_NAME, 'readwrite')
    const store = tx.objectStore(STORE_NAME)
    await new Promise((resolve, reject) => {
      const req = store.put({ docId, snapshot: snapshotBase64, updatedAt: new Date().toISOString() })
      req.onsuccess = () => resolve()
      req.onerror = () => reject(req.error)
    })
  } finally {
    db?.close()
  }
}

// Remove cached snapshot for a doc.
export async function deleteCachedSnapshot(docId) {
  let db
  try {
    db = await openDB()
    const tx = db.transaction(STORE_NAME, 'readwrite')
    const store = tx.objectStore(STORE_NAME)
    await new Promise((resolve, reject) => {
      const req = store.delete(docId)
      req.onsuccess = () => resolve()
      req.onerror = () => reject(req.error)
    })
  } finally {
    db?.close()
  }
}
