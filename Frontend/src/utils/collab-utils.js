import { ySyncPluginKey } from 'y-prosemirror'

export function isChangeOrigin(transaction) {
  return !!transaction.getMeta(ySyncPluginKey)
}
