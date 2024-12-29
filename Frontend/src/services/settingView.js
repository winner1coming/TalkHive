import apiClient from '@/services/api';


  // 获取用户信息接口
  export const showProfile = async () => {
    try {
      const response = await apiClient.get(`/Settings/profile`);
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
  export const getUserInfo = async()=>{
    try{
      const response = await apiClient.get('/Settings/getInfo');
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
      const response = await apiClient.post('/Settings/saveEmail',data);
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
      const response = await apiClient.post('/Settings/isIDAdd',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //是否可以通过昵称查找并添加我
  export const isNicknameAdd = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/isNicknameAdd',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //更改字体大小
  export const changeFontsize = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/fontsize',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //更改字体风格
  export const changeFontstyle = async(data)=>{
    try{
      const response = await apiClient.post('Settings/fontStyle',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //更改主题
  export const changeTheme = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/changeTheme',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //更改聊天背景
  export const changeBackground = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/changeBackground',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //开启消息通知
  export const isNotice = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/isNotice',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //设置群聊通知
  export const isNoticeGroup = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/noticeGroup',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //更换消息提示音
  export const changeSound = async(data)=>{
    try{
      const response = await apiClient.post('/Settings/changeSound',data);
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //获取系统设置
  export const getSystemSetting = async()=>{
    try{
      const response = await apiClient.get('/systemSetting');
      return response.data;
    }catch(error){
      throw error.response?.data.message || error.message;
    }
  };

  //注销账号
  export const confirmDeactivation = async()=>{
    try{
      const response = await apiClient.post('/Settings/deactivate');
      return response.data;
    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };

  export const logout = async()=>{
    try{
      const response = await apiClient.post('/Logout');
      return response.data;
    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };

  //链接器的添加
  export const AddLinks = async(links)=>{
    try{
      const response = await apiClient.post('/addLinks',{links});
      return response.data;
    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };

  export const DelLinks = async(links)=>{
    try{
      const response = await apiClient.post('/delLinks',{links});
      return response.data;
    }catch(error){
      throw error.response?.data.message ||error.message;
    }
  };