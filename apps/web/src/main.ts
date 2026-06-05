import { createApp } from 'vue'

import App from './App.vue'
import { revealDirective } from './directives/reveal'
import router from './router'
import './styles/global.css'

const app = createApp(App)

app.directive('reveal', revealDirective)
app.use(router)
app.mount('#app')
