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


// 资料卡片


// 拦截 /profileCard 的 GET 请求
Mock.mock(new RegExp(`${baseURL}/profileCard/\\d+`), 'get', (options) => {
	const profileCard = Mock.mock({
		'profileCard': {
		'tid|1': '2',
		'remark': '@name',
		'nickname': '@name',
		'groupNickname': '@name',
		'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
		'status': '@pick(["online", "offline"])',
		'signature': '@sentence',
		'tag': '@pick(["家人", "朋友", "同事"])',
		}
	});
	return profileCard.profileCard;
});