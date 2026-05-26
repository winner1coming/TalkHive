import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import TextAlign from '@tiptap/extension-text-align'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { Collaboration } from './collab-collaboration'
import { ViewportPlugin } from '../plugins/viewport-plugin'

// 获取所有扩展（协作模式：带 Yjs）
export function getCollabExtensions(ydoc, provider, userInfo) {
  return [
    StarterKit.configure({
      // 使用 ProseMirror 内置 history（自定义 collaboration 未包含 yUndoPlugin）
      history: true,
      link: false,
      underline: false,
    }),
    Underline,
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Image.configure({
      inline: false,
      allowBase64: true,
    }),
    Link.configure({
      openOnClick: true,
    }),
    Placeholder.configure({
      placeholder: '开始编辑...',
    }),
    // Yjs 协作（自定义扩展，包含 ySyncPlugin + yCursorPlugin）
    Collaboration.configure({
      document: ydoc,
      field: 'content',
      provider,
    }),
    // 大文档视口优化
    ViewportPlugin,
  ]
}

// 获取所有扩展（单机模式：无 Yjs，用于普通笔记）
export function getLocalExtensions() {
  return [
    StarterKit.configure({
      history: true,
      link: false,
      underline: false,
    }),
    Underline,
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Image.configure({
      inline: false,
      allowBase64: true,
    }),
    Link.configure({
      openOnClick: true,
    }),
    Placeholder.configure({
      placeholder: '开始编辑...',
    }),
    // 大文档视口优化
    ViewportPlugin,
  ]
}
