// src/mocks/mock.js
import Mock from 'mockjs';
import './chatListMock';
import './contactListMock';

const baseURL = 'http://localhost:8080';


Mock.setup({
  timeout: '200-300', // 设置模拟延迟（可选）
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization'
  }
});

// 模拟收藏列表接口
Mock.mock(`${baseURL}/workspace/favorites`, 'get', () => {
  // // 获取用户 ID
  // const userId = new URLSearchParams(options.url.split('?')[1]).get('id');
  // if (!userId) {
  //   return Mock.mock({
  //     status: 400,
  //     message: '用户ID不可为空',
  //   });
  // }

  // 模拟返回数据
  return Mock.mock({
    data: [
      {
        message_table_name: 'message_table1',
        message_id: '1',
        type: 'message',
        object_name: '这是一条消息内容',
        sender_name: 'Alice',
        time: '10:00',
      },
      {
        message_table_name: 'note_table1',
        message_id: '2',
        type: 'note',
        object_name: '如何学习 Vue.js',
        sender_name: 'Bob',
        time: '14:00',
      },
      {
        message_table_name: 'code_table1',
        message_id: '3',
        type: 'code',
        object_name: '排序算法实现',
        sender_name: 'Charlie',
        time: '18:00',
      },
    ],
  });
});

//删除收藏
Mock.mock(`${baseURL}/workspace/favorite/delete`, 'post', (options) =>{
  const { items } = JSON.parse(options.body);

  // 打印接收到的参数
  console.log(options.body);

  // 验证请求是否合法
  if (!Array.isArray(items) || items.length === 0) {
    return options.status(400).json({ message: '无效的删除请求' });
  }

  // 提取所有 message_id 和 message_type
  const messageIds = items.map(item => item.message_id);
  const messageTables = items.map(item => item.message_table_name);
  console.log("messageIds:", messageIds);
  console.log("messageTables:", messageTables);
  return{
    status: 200,
  }
});

// 模拟返回 DDL 列表
Mock.mock(`${baseURL}/workspace/ddl/pending`, 'get', (options) => {
    //const userId = new URLSearchParams(options.url.split('?')[1]).get('id');
    // if (!userId) {
    //   return Mock.mock({
    //     status: 400,
    //     message: '用户ID不可为空',
    //   });
    // }
  
    // 模拟 DDL 列表数据
    return Mock.mock({
      status: 200,
      data: [
        {
          task_id: '1',
          deadline: '2024-12-20 17:00',
          task_content: '完成前端开发任务',
          important: false,
        },
        {
          task_id: '2',
          deadline: '2024-12-22 12:00',
          task_content: '撰写项目需求文档',
          important: false,
        },
        {
          task_id: '3',
          deadline: '2024-12-25 09:00',
          task_content: '进行代码审查',
          important: true,
        },
      ],
    });
  });

// 模拟返回已完成的 DDL 列表
Mock.mock(`${baseURL}/workspace/ddl/completed`, 'get', (options) => {
    // const userId = new URLSearchParams(options.url.split('?')[1]).get('id');
    // if (!userId) {
    //   return Mock.mock({
    //     status: 400,
    //     message: '用户ID不可为空',
    //   });
    // }
    
    // 模拟已完成的 DDL 列表数据
    return Mock.mock({
      status: 200,
      data: [
        {
          task_id: '4',
          deadline: '2024-12-15 09:00',
          task_content: '完成数据库设计',
        },
        {
          task_id: '5',
          deadline: '2024-12-18 18:00',
          task_content: '完成后台 API 开发',
        },
      ],
    });
  });

// 模拟新建ddl列表接口
Mock.mock(`${baseURL}/workspace/ddl/create`, 'post', (options)=>{
  const { deadline, task_content, important } = JSON.parse(options.body);
  return{
    status: 200,
    data: {
      deadline: deadline,
      task_content: task_content,
      important: important,
    }
  }
});

// 模拟修改ddl列表接口
Mock.mock(`${baseURL}/workspace/ddl/edit`, 'post', (options)=>{
  const { task_id, deadline, task_content, important } = JSON.parse(options.body);
  return{
    status: 200,
    data: {
      task_id: task_id,
      deadline: deadline,
      task_content: task_content,
      important: important,
    }
  }
});

// 模拟更新 DDL 状态为已完成接口
Mock.mock(`${baseURL}/workspace/ddl/update`, 'post', {
  status: 200
});

// 模拟删除ddl
Mock.mock(`${baseURL}/workspace/ddl/delete`, 'post', {
  status: 200
});


// 获取分类列表
Mock.mock(`${baseURL}/workspace/notes/categories`, 'get', {
  status: 200,
  message: '获取成功',
  categories: [
    '工作', 
    '学习', 
    'cp' , 
    'Matriarchy',
  ]
});

// 模拟删除分类列表接口
Mock.mock(`${baseURL}/workspace/notes/categories/delete`, 'post', (options)=>{
  const { type_name } = JSON.parse(options.body);
  return{
    status: 200,
    data:{
      deletedCategory: type_name
    }
  }
});

// 模拟新建分类列表接口
Mock.mock(`${baseURL}/workspace/notes/categories/new`, 'post', {
  status: 200,
  message: '新建成功'
});

// 模拟修改分类列表接口
Mock.mock(`${baseURL}/workspace/notes/categories/edit`, 'post', (options)=>{
  const { old_type_name, new_type_name } = JSON.parse(options.body);
  return{
    status: 200,
    data: {
      old: old_type_name,
      new: new_type_name,
    }
  }
});

  // 模拟获取笔记列表接口
Mock.mock(`${baseURL}/workspace/notes`, 'get', {
  status: 200,
  message: '获取成功',
  notes: [
    { id: 1, filename: 'Vue学习笔记.md', lastModified: '2024-12-01 10:30', category: '工作' },
    { id: 2, filename: '项目需求分析.docx', lastModified: '2024-12-05 14:15', category: '学习' },
    { id: 3, filename: '代码优化方案.txt', lastModified: '2024-12-10 09:45', category: '个人' },
    { id: 4, filename: '糖点总结.md', lastModified: '2024-12-12 16:00', category: 'cp' }
  ]
});

// 删除笔记文件
Mock.mock(`${baseURL}/workspace/notes/deletenote`, 'post', (options) => {
  const { note_id } = JSON.parse(options.body);
  return {
    status: 200,
    message: ` 文件 id ${note_id} 已删除`,
  };
});

// 模拟创建文件接口
Mock.mock('/workspace/create-file', 'post', (options) => {
  const { filename } = JSON.parse(options.body);
  return {
    status: 200,
    message: `文件 ${filename} 创建成功`
  };
});

// 获取笔记内容
// 模拟请求返回数据
// Mock.mock(`${baseURL}/workspace/notes/get`, 'post', (options) => {
//   const { note_id } = JSON.parse(options.body); // 解析请求体中的参数
  
//   // 模拟响应数据
//   const note = Mock.mock({
//     'note_id': note_id,
//     'name': 'example.md',
//     'suffix': '.md',
//   });

//   // 模拟文件内容
//   const fileContent = `# Example Markdown File\n\nThis is a sample markdown content for file ${code.name}.`;

//   // 返回模拟的响应数据
//   if (code.is_show) {
//     return {
//       code_id: code.code_id,
//       file_path: code.cachePath,
//       file_exists: true,
//       file_content: fileContent, // 返回文件的内容
//       content_type: 'text/markdown', // 假设返回 markdown 文件的 content-type
//     };
//   } else {
//     return {
//       message: "Code not found or not visible",
//       status: 404,
//     };
//   }
// });
// Mock.mock(`${baseURL}/workspace/notes/get`, 'post', (options) => {
//   // 模拟文件内容
//   const fileContent = 'This is a test file content.';
  
//   // 将文件内容转为 Base64 编码
//   const base64Content = btoa(fileContent); // 使用 btoa 将字符串转为 Base64 编码

//   // 返回模拟的响应
//   return {
//     status: 'success',
//     data: {
//       fileName: 'example.txt',
//       content: base64Content, // 返回 Base64 编码的文件内容
//       contentType: 'text/plain' // 文件的 MIME 类型
//     },
//     message: 'File fetched successfully'
//   };
// });
// Mock.mock(`${baseURL}/workspace/notes/get`, 'post', (options) => {
//   const { codeID } = JSON.parse(options.body);
//   const fileContent = `This is the content of file with codeID: ${codeID}\nThis is line 2 of the file.\nThis is line 3 of the file.`;

//   // 确保返回的是 Blob 类型的数据
//   const blob = new Blob([fileContent], { type: 'text/plain' });

//   return blob;
// });
Mock.mock(`${baseURL}/workspace/notes/get`, 'post', (options) => {
  const { note_id } = JSON.parse(options.body);
  // 模拟后端响应的数据
  // 定义模拟的笔记文件内容
  const markdownText = `codeID: ${note_id}\n
  # 这是一个标题

  这是一些笔记内容，支持 **Markdown** 格式。

  - 项目 1
  - 项目 2
  - 项目 3

  > 这是一个引用块

  \`\`\`javascript
  console.log('这是代码块');
  \`\`\`
  `;

  return {
    status: 200,
    message: '成功',
    data: markdownText,  // 这里返回模拟的文件流内容（Markdown文本）
  };
});


// 代码————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
// 获取代码列表
Mock.mock(`${baseURL}/workspace/code/list`, 'get', {
  status: 200,
  message: '获取成功',
  codes: [
    { code_id: 1, code_name: "Vue学习笔记.md", last_modified_time: "2024-12-01 10:30", Suffix: '.js' },
    { code_id: 2, code_name: "pythontest.py", last_modified_time: "2024-12-05 14:15", Suffix: '.cpp' },
    { code_id: 3, code_name: "javascripttest.js", last_modified_time: "2024-12-10 09:45", Suffix: '.py' }
  ]
});

// 删除代码文件
Mock.mock(`${baseURL}/workspace/code/delete`, 'post', (options) => {
  const { code_id } = JSON.parse(options.body);
  return {
    status: 200,
    message: ` 文件 id ${code_id} 已删除`,
  };
});

//回收站————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
//—————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
// 模拟已删除文件列表接口
Mock.mock(`${baseURL}/workspace/recycle/files`, 'get', () => {
  return {
    files: [
      {
        recycle_id: 1,
        filename: 'file1.txt',
        type: 'note',
        recycle_time: '2024-12-01 10:30',
      },
      {
        recycle_id: 2,
        filename: 'file2.docx',
        type: 'note',
        recycle_time: '2024-12-05 14:15',
      },
      {
        recycle_id: 3,
        filename: 'file3.pdf',
        type: 'code',
        recycle_time: '2024-12-10 09:45',
      },
    ],
  };
});

// 模拟恢复文件接口
Mock.mock(`${baseURL}/workspace/recycle/restore-file`, 'post', (options) => {
  const { type, recycle_id } = JSON.parse(options.body);
  return {
    status: 200,
    message: ` ${type} 文件 id ${recycle_id} 恢复成功`,
  };
});

// 模拟彻底删除文件接口
Mock.mock(`${baseURL}/workspace/recycle/delete-permanent`, 'post', (options) => {
  const { type, recycle_id } = JSON.parse(options.body);
  return {
    status: 200,
    message: ` ${type} 文件 id ${recycle_id} 已彻底删除`,
  };
});


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
	return {data:profileCard.profileCard};
});
