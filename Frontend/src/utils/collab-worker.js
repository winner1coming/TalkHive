// Worker wrapper: manages a Web Worker instance and provides async helpers
// for offloading Yjs epoch compaction off the main thread.

let worker = null
let nextId = 0
const pending = new Map()

function getWorker() {
  if (worker) return worker
  worker = new Worker(new URL('@/workers/yjs-worker.js', import.meta.url), { type: 'module' })
  worker.onmessage = (e) => {
    const { id, type, data, error } = e.data
    const cb = pending.get(id)
    if (!cb) return
    pending.delete(id)
    if (type === 'error') {
      cb.reject(new Error(error))
    } else {
      cb.resolve(data)
    }
  }
  worker.onerror = (err) => {
    // For all pending callbacks, reject with the error
    for (const [, cb] of pending) {
      cb.reject(new Error(err.message))
    }
    pending.clear()
  }
  return worker
}

// Send a Yjs state update to the worker for epoch compaction.
// Returns a Promise that resolves with the compacted state as an ArrayBuffer.
// Falls back gracefully: if the worker fails, the caller should fall back to
// main-thread compaction.
export function workerEpochCompact(update) {
  return new Promise((resolve, reject) => {
    const id = nextId++
    pending.set(id, { resolve, reject })
    try {
      const w = getWorker()
      // Transfer the buffer to the worker (zero-copy on supported browsers).
      // Send a CLONE so the original `update` remains valid for fallback path.
      const clone = update.buffer.slice(0)
      w.postMessage({ id, type: 'epochCompact', data: clone }, [clone])
    } catch (err) {
      pending.delete(id)
      reject(err)
    }
  })
}

// Terminate the worker and clean up. Call during editor teardown.
export function terminateWorker() {
  if (worker) {
    worker.terminate()
    worker = null
  }
  for (const [, cb] of pending) {
    cb.reject(new Error('Worker terminated'))
  }
  pending.clear()
}
