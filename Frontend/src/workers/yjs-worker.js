// Web Worker for CPU-intensive Yjs operations.
// Runs off the main thread to prevent UI freezing during epoch compaction.

import * as Y from 'yjs'

self.onmessage = (e) => {
  const { id, type, data } = e.data

  try {
    switch (type) {
      case 'epochCompact': {
        // data is a transferable Uint8Array buffer containing the full Yjs state update
        const update = new Uint8Array(data)
        const freshDoc = new Y.Doc()
        Y.applyUpdate(freshDoc, update)
        const compacted = Y.encodeStateAsUpdate(freshDoc)
        freshDoc.destroy()
        // Transfer the buffer back to main thread (zero-copy)
        self.postMessage(
          { id, type: 'epochCompact', data: compacted.buffer },
          [compacted.buffer]
        )
        break
      }

      default:
        self.postMessage({ id, type: 'error', error: `Unknown worker command: ${type}` })
    }
  } catch (err) {
    self.postMessage({ id, type: 'error', error: err.message })
  }
}
