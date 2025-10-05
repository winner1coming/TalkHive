import apiClient from '@/services/api';

// DDL——————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
// 获取待完成的 DDL 列表
export const getDdlList = () => {
    return apiClient.get('/workspace/ddl/pending');
  };
  
// 获取已完成的 DDL 列表
export const getCompletedDdlList = () => {
  return apiClient.get('/workspace/ddl/completed');
};

// 保存修改后的 DDL
export const saveEditDdl = (task_id, deadline, task_content, important) => {
  return apiClient.post('/workspace/ddl/update',{task_id, deadline, task_content, important});
};

// 保存新建的 DDL
export const saveDdl = (deadline, task_content, important) => {
  return apiClient.post('/workspace/ddl/create',{deadline, task_content, important});
};

// 更新 DDL 状态为已完成
export const updateDdl = (task_id) => {
  return apiClient.post('/workspace/ddl/complete',{task_id});
};

// 删除ddl
export const deleteDDL = (task_id) => {
  return apiClient.post('/workspace/ddl/delete',{task_id});
};

// 收藏—————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
export const getFavorites = () => {
  return apiClient.get('/workspace/favorites/list');
};

export const deleteFavorites = (msg_ids) => {
  return apiClient.post('/workspace/favorites/delete',msg_ids);
};

// 笔记—————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
//获取分类列表
export const getCategories = () => {
  return apiClient.get('/workspace/notes/categories');
};

// 新建分类
export const saveCategory = (type_name) => {
  return apiClient.post('/workspace/notes/categories/new',{type_name});
};

// 删除分类
export const deleteCategory = (type_name) => {
  return apiClient.post('/workspace/notes/categories/delete', {type_name});
};

// 保存分类修改
export const saveEditCategory = (old_type_name, new_type_name) => {
  return apiClient.post('/workspace/notes/categories/edit',{old_type_name, new_type_name});
};

// 获取笔记列表
export const getNotes = () => {
  return apiClient.get('/workspace/notes/list');
};

// 删除笔记
export const deleteNote = (note_id) => {
  return apiClient.post('/workspace/notes/deletenote', {note_id});
};

// 新建笔记
export const createNote = (note_name, type, content) => {
  return apiClient.post('/workspace/notes/newnote',{note_name, type, content});
};

// 获取笔记内容
export const getNoteContent = (note_id) => {
  return apiClient.post('/workspace/notes/get', {note_id}, { responseType: 'text',});
};

//保存笔记修改
export const saveEditNote = (NoteID, NoteName, Type, Content) => {
  console.log(NoteID, NoteName, Type, Content);
  return apiClient.post('/workspace/notes/editnote',{NoteID, NoteName, Type, Content});
};

// 普通文档————————————————————————————————————————————————————————————————————————————————————————————————————
// // 获取文档列表
// export const getQuillDocs = () => {
//   return apiClient.get('/workspace/docs/list');
// };

// // 创建新文档
// export const createQuillDoc = (doc_name) => {
//   return apiClient.post('/workspace/docs/newdoc',{doc_name});
// };

// // 保存delta
// export const saveQuillContent = (doc_id, delta) => {
//   return apiClient.post('workspace/docs/savesnapshot', {doc_id, delta});
// }

// // 获取delta
// export const getQuillContent = (doc_id) => {
//   return apiClient.post('workspace/docs/getsnapshot', {doc_id});
// }

// 协作文档————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
// 获取参与过的文档列表
export const getDocs = () => {
  return apiClient.get('/workspace/collabdocs/list');
};

// 创建新文档
export const createDoc = (doc_name) => {
  return apiClient.post('/workspace/collabdocs/newdoc',{doc_name});
};

// 保存快照
export const saveSnapshot = (doc_id, snapshot) => {
  return apiClient.post('workspace/collabdocs/savesnapshot', {doc_id, snapshot});
}

// 获取快照
export const getSnapshot = (doc_id) => {
  return apiClient.post('workspace/collabdocs/getsnapshot', {doc_id});
}

// 邀请新成员
export const joinMember = (doc_id) => {
  return apiClient.post('workspace/collabdocs/join', {doc_id});
}

// 回收站————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
// 获取回收站列表
export const getRecycles = () => {
  return apiClient.get('/workspace/recycle/files');
};

// 恢复文件
export const restoreFile = (type, recycle_id) => {
  return apiClient.post('/workspace/recycle/restore-file', {type, recycle_id});
};

// 彻底删除文件
export const deleteFile = (type, recycle_id) => {
  return apiClient.post('/workspace/recycle/delete-permanent', {type, recycle_id});
};

// 代码——————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
// 获取代码列表
export const getCodes = () => {
  return apiClient.post('/workspace/code/list');
};

//删除代码
export const deleteCode = (code_id) => {
  return apiClient.post('/workspace/code/delete', {code_id});
};
// 新建代码
export const createCode = (Name, Suffix) => {
  return apiClient.post('/workspace/code/new',{Name, Suffix});
};
// 获取代码内容
export const getCodeContent = (code_id) => {
  return apiClient.post('/workspace/code/get', {code_id}, { responseType: 'text',});
};
// 保存编辑
export const saveEditCode = (code_id, code_name, suffix, content) => {
  console.log(code_id, code_name, suffix, content);
  return apiClient.post('/workspace/code/edit',{code_id, code_name, suffix, content});
};

