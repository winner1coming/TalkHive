import Mock from 'mockjs';
const baseURL = 'http://your-api-url.com';

import './chatListMock';
import './contactListMock';

// // 请求拦截器
// Mock.XHR.prototype.proxy_send = Mock.XHR.prototype.send;
// Mock.XHR.prototype.send = function() {
//   console.log('Request:', this.custom);
//   this.proxy_send(...arguments);
// };

// // 响应拦截器
// Mock.XHR.prototype.proxy_open = Mock.XHR.prototype.open;
// Mock.XHR.prototype.open = function() {
//   this.addEventListener('load', function() {
//     console.log('Response:', this.responseText);
//   });
//   this.proxy_open(...arguments);
// };