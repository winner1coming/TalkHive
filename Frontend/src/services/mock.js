import Mock from 'mockjs';
const baseURL = 'http://your-api-url.com';

// chatlist
const chats = Mock.mock({
  'chats|5-10': [{
    'id|1': /[0-9]{10}/,
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'name': '@name',
    'lastMessage': '@sentence',
    'lastMessageTime': '@time("HH:mm")',
    'unreadCount|0-10': 1,
    'tags': function() {
        return Mock.mock({
          'array|1-3': ['@pick(["friend", "group", "unread", "pinned", "blocked"])']
        }).array;
      },
  }]
});
const messages = Mock.mock({
  'messages|5-10': [{
    'message_id|+1': 1,
    'send_account_id|1-100': 1,
    'content': '@sentence',
    'sender': '@name',
    'create_time': '@time("HH:mm")',
    'type': 'text',
  }]
});
Mock.mock(`${baseURL}/chatlist`, 'get', () => {
  return chats.chats;   // 失败时返回 reason: 失败原因
});
Mock.mock(`${baseURL}/chatlist/createChat`, 'post', (req) => {
  const chat = Mock.mock({
    'chat': {
      'id|1': /[0-9]{10}/,
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'name': '@name',
      'lastMessage': '@sentence',
      'lastMessageTime': '@time("HH:mm")',
      'unreadCount|0-10': 1,
      'tags': function() {
        return Mock.mock({
          'array|1-3': ['@pick(["friend", "group", "unread", "pinned", "blocked"])']
        }).array;
      },
    }
  });
  return {
    status: 200,
    data: chat.chat,
  };
});
Mock.mock(new RegExp(`${baseURL}/chatlist/search/\\w+`), 'get', (options) => {
  const keyword = options.url.split('/').pop();
  const filteredChats = chats.chats.filter(chat => chat.name.includes(keyword));
  return {
    status: 200,
    data: filteredChats,
  };
});
Mock.mock(`${baseURL}/chatlist/pin`, 'post', (options) => {
  console.log("api pin");
  const { tid, is_pinned } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_pinned) {
      chat.tags.push('pinned');
    } else {
      chat.tags = chat.tags.filter(tag => tag !== 'pinned');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});
Mock.mock(`${baseURL}/messages/read`, 'post', (options) => {
  const { tid, is_read } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_read) {
      chat.tags = chat.tags.filter(tag => tag !== 'unread');
    } else {
      chat.tags.push('unread');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});

Mock.mock(new RegExp(`${baseURL}/chatlist/\\d+`), 'delete', (options) => {
  const tid = parseInt(options.url.split('/').pop());
  chats.chats = chats.chats.filter(chat => chat.id !== tid);
  return {
    status: 200,
    data: null,
  };
});

Mock.mock(`${baseURL}/chatlist/mute`, 'post', (options) => {
  const { tid, is_muted } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_muted) {
      chat.tags.push('mute');
    } else {
      chat.tags = chat.tags.filter(tag => tag !== 'mute');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});

Mock.mock(`${baseURL}/chatlist/block`, 'post', (options) => {
  const { tid, is_blocked } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_blocked) {
      chat.tags.push('blocked');
    } else {
      chat.tags = chat.tags.filter(tag => tag !== 'blocked');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});

Mock.mock(new RegExp(`${baseURL}/messages/\\d+`), 'get', (options) => {
  const tid = parseInt(options.url.split('/').pop());
  return {
    status: 200,
    data: messages.messages,
  };
});

Mock.mock(/\/messages\/\d+\/send/, 'post', (options) => {
  const tid = parseInt(options.url.split('/')[2]);
  const { content } = JSON.parse(options.body);
  const newMessage = Mock.mock({
    'message_id': Mock.Random.guid(),
    'send_account_id': Mock.Random.integer(1, 100),
    'content': content,
    'sender': '@name',
    'create_time': '@time("HH:mm")',
    'type': 'text',
  });
  return {
    status: 200,
    data: newMessage,
  };
});