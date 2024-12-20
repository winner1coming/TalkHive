import apiClient from '@/services/api';

  // 登录接口
  export const login = async (payload) => {
    try {
      const response = await apiClient.post('/login', payload);
      return response.data;
    } catch (error) {
      throw error.response ? error.response.data.message : error.message;
    }
  };
  //短信接收码的接口
  export const sendSmsCode = async (data) => {
    try {
      const response = await apiClient.post('/sendSmsCode',data);
      return response.data;
    } catch (error) {
      throw error.response ? error.response.data.message : error.message;
    }
  };

  export const smsLogin = async(email) => {
    try{
      const response =await apiClient.post('/smslogin',email);
        return response.data;
    }catch(error){
        throw error.response?.data.message || error.message;
      };
  };

  // 注册接口
  export const Register =async (data) => {
    try{
      const response = await apiClient.post('/register', data);
      return response.data;
    } 
    catch(error) {
      throw error.response?.data.message || error.message;
    };
  };


  // 重置密码接口
  export const resetPassword =async (msg) => {
    try{
      const response = await apiClient.post('/resetPassword', msg);
      return response.data;
    } 
    catch(error) {
      throw error.response?.data.message || error.message;
    };
  };
