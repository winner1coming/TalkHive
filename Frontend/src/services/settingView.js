import apiClient from '@/services/api';


  // 获取用户信息接口
  export const showProfile = async (id) => {
    try {
      const response = await apiClient.get(`/Settings/profile/${id}`);
      return response.data;
    } catch (error) {
      throw error.response?.data.message || error.message;
    }
  };

  // 更新用户信息接口
  export const saveEdit = async (data) => {
    try {
      const response = await apiClient.post('/Settings/saveEdit', data);
      return response.data;
    } catch (error) {
      throw error.response?.data.message || error.message;
    }
  };

  //获取用户的邮箱（安全设置）
  export const getUserInfo = async(id)=>{
    try{
      const response = await apiClient.get('/Settings/getInfo',{id});
      return response.data;

    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };

  //更换邮箱时获取验证码
  export const getCode = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/getCode',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message||error.message;
    }
  };
  //保存Email的修改
  export const saveEmail = async (data) => {
    try {
      const response = await apiClient.post('/Settings/saveEmail', data);
      return response.data;
    } catch (error) {
      throw error.response?.data.message || error.message;
    }
  };

  //保存密码的更改
  export const savePassword =async(data)=>{
    try{
      const response = await apiClient.post('/Settings/savePassword',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //是否可以通过ID添加我
  export const isIDAdd = async(data)=>{
    try{
      const response = await apiClient.get('/Settings/isIDAdd',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //是否可以通过昵称查找并添加我
  export const isNicknameAdd = async(data)=>{
    try{
      const response = await apiClient.get('/Settings/isNicknameAdd',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  }


  //注销账号
  export const confirmDeactivation = async(id)=>{
    try{
      const response = await apiClient.post('/Settings/deactivate',{id});
      return response.data;
    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };

  export const logout = async(id)=>{
    try{
      const response = await apiClient.post('/Logout',{id});
      return response.data;
    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };