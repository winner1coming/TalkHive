import { Extension } from '@tiptap/core'
import { Plugin, PluginKey } from '@tiptap/pm/state'

// Viewport tracking plugin for large-document optimization.
//
// Phase 2: monitors which block nodes are near the viewport using
// IntersectionObserver. Provides visibility data for future enhancements
// (placeholder rendering, progressive loading).
//
// The per-node content-visibility:auto CSS (in CollabEditor.vue) is the
// primary off-screen optimization — this plugin adds observation and
// metadata on top of that.

export const ViewportPlugin = Extension.create({
  name: 'viewportPlugin',

  addOptions() {
    return {
      // Margin around the viewport to consider nodes "near-viewport" (px)
      rootMargin: '2000px 0px',
    }
  },

  addProseMirrorPlugins() {
    const pluginKey = new PluginKey('viewport')

    return [
      new Plugin({
        key: pluginKey,
        state: {
          init() {
            return { visibleSet: new Set(), observer: null }
          },
          apply(tr, prev) {
            // Pass through — state is managed by the DOM observer
            return prev
          },
        },
        view: (view) => {
              const rootMargin = this.options.rootMargin
              let observer = null

              return {
                update: (view, prevState) => {
                  const dom = view.dom
                  if (!dom || !dom.parentElement) return

                  if (observer) observer.disconnect()

                  try {
                    observer = new IntersectionObserver(
                      (entries) => {
                        for (const entry of entries) {
                          const el = entry.target
                          if (entry.isIntersecting || entry.intersectionRatio > 0) {
                            el.removeAttribute('data-viewport-offscreen')
                          } else {
                            el.setAttribute('data-viewport-offscreen', 'true')
                          }
                        }
                      },
                      { root: dom.parentElement, rootMargin }
                    )

                // Observe all direct block children of the ProseMirror editor
                for (const child of dom.children) {
                  if (child instanceof HTMLElement) {
                    observer.observe(child)
                  }
                }
              } catch (e) {
                // IntersectionObserver not supported or failed — degrade gracefully
              }
            },
            destroy() {
              if (observer) observer.disconnect()
            },
          }
        },
      }),
    ]
  },
})
