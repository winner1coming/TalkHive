// 本地存储工具函数

// 保存数据到本地存储
export const saveToLocalStorage = (key, value) => {
  localStorage.setItem(key, JSON.stringify(value));
};

// 从本地存储获取数据
export const getFromLocalStorage = (key) => {
  const value = localStorage.getItem(key);
  return value ? JSON.parse(value) : null;
};

// 从本地存储删除数据
export const removeFromLocalStorage = (key) => {
  localStorage.removeItem(key);
};