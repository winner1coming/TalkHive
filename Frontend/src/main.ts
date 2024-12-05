import {createApp} from 'vue';  

import App from './App.vue'
import router from './router';
import store from './store';

const app = createApp(App);
app.use(router);
app.use(store);

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
