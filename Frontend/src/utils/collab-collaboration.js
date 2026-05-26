import { Extension } from '@tiptap/core'
import { Plugin, PluginKey } from '@tiptap/pm/state'
import {
  ySyncPluginKey,
  ySyncPlugin,
  yCursorPlugin,
  ProsemirrorBinding,
  getRelativeSelection,
  updateYFragment,
} from 'y-prosemirror'
import { isChangeOrigin } from './collab-utils'

// Prevent the editor viewport from auto-scrolling when remote Yjs updates
// add content at the document end. Without this, collaborative editing in
// large documents causes the local view to jump to the bottom whenever
// another user types.
//
// Note: `scrolledIntoView` is a getter-only property (computed from the
// `updated` bitmask), so we clear bit 4 (UPDATED_SCROLL) directly instead.
const PreventRemoteScroll = new Plugin({
  key: new PluginKey('preventRemoteScroll'),
  view(view) {
    const origDispatch = view.dispatch.bind(view)
    view.dispatch = function (tr) {
      if (tr.getMeta(ySyncPluginKey)) {
        tr.updated &= ~4  // clear UPDATED_SCROLL bit
      }
      return origDispatch(tr)
    }
    return {
      destroy() {
        view.dispatch = origDispatch
      },
    }
  },
})

// Wrap beforeAllTransactions with error handling to prevent the
// "AbsolutePositionToRelativePosition → findRootTypeKey" crash when
// concurrent Yjs updates modify the document structure while ProseMirror
// still has the old selection. Without this, remote Yjs updates can't
// sync to the local ProseMirror editor (breaks collaboration).
//
// IMPORTANT: origInitView internally OVERWRITES beforeAllTransactions,
// so our try-catch wrapper must be assigned AFTER calling origInitView.
const origInitView = ProsemirrorBinding.prototype.initView
ProsemirrorBinding.prototype.initView = function (view) {
  origInitView.call(this, view)
  this.beforeAllTransactions = () => {
    if (this.beforeTransactionSelection === null && this.prosemirrorView != null) {
      try {
        this.beforeTransactionSelection = getRelativeSelection(
          this,
          this.prosemirrorView.state
        )
      } catch (e) {
        // Selection capture failed (concurrent Yjs/Pm state mismatch).
        // Reset selection — the transaction still proceeds and content syncs,
        // cursor may jump to start which is acceptable.
        this.beforeTransactionSelection = null
      }
    }
  }
}

// Wrap _prosemirrorChanged with error handling to prevent the same
// absolutePositionToRelativePosition crash when pressing Enter (paragraph
// split triggers updateYFragment then getRelativeSelection during the
// transaction, and the do-while loop in absolutePositionToRelativePosition
// can encounter n._item === null when climbing parent chain).
const origProsemirrorChanged = ProsemirrorBinding.prototype._prosemirrorChanged
ProsemirrorBinding.prototype._prosemirrorChanged = function (doc) {
  this.doc.transact(() => {
    updateYFragment(this.doc, this.type, doc, this)
    try {
      this.beforeTransactionSelection = getRelativeSelection(this, this.prosemirrorView.state)
    } catch (e) {
      this.beforeTransactionSelection = null
    }
  }, ySyncPluginKey)
}

export const Collaboration = Extension.create({
  name: 'collaboration',
  priority: 1000,

  addOptions() {
    return {
      document: null,
      field: 'default',
      fragment: null,
      provider: null,
    }
  },

  addProseMirrorPlugins() {
    const fragment = this.options.fragment
      ? this.options.fragment
      : this.options.document.getXmlFragment(this.options.field)

    const ySyncPluginInstance = ySyncPlugin(fragment)

    // Remote user cursors via y-prosemirror's yCursorPlugin
    const cursorPlugin = this.options.provider
      ? yCursorPlugin(this.options.provider.awareness)
      : null

    return [
      ySyncPluginInstance,
      PreventRemoteScroll,
      cursorPlugin,
      new Plugin({
        key: new PluginKey('filterInvalidContent'),
        filterTransaction: (transaction) => {
          if (!isChangeOrigin(transaction)) return true
          if (!transaction.docChanged) return true
          try {
            transaction.doc.check()
            return true
          } catch {
            return false
          }
        },
      }),
    ].filter(Boolean)
  },
})
