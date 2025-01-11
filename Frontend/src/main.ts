import {createApp} from 'vue';  

import App from './App.vue'
import router from './router';
import store from './store';
//import './services/mock/mock';  // todo debug: 项目完成时删除
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';
import { library } from '@fortawesome/fontawesome-svg-core';
import {faComment, faAddressBook, faCog,faBriefcase,faSignOutAlt} from '@fortawesome/free-solid-svg-icons';

library.add(faComment,faAddressBook,faCog,faBriefcase,faSignOutAlt);


const app = createApp(App);
app.use(router);
app.use(store);
app.component('font-awesome-icon', FontAwesomeIcon);

app.mount('#app');


const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
  store.dispatch('setTheme', savedTheme);
}

store.watch(
    (state) => state.settings.theme,
    (theme) => {
      document.documentElement.className = theme;
    }
  );
