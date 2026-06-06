import { createRouter, createWebHashHistory } from 'vue-router'

import { fetchMe } from '@/auth'
import LandingPage from '@/pages/LandingPage.vue'
import LoginPage from '@/pages/LoginPage.vue'
import RegisterPage from '@/pages/RegisterPage.vue'
import WorkspacePage from '@/pages/WorkspacePage.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'landing', component: LandingPage },
    { path: '/login', name: 'login', component: LoginPage },
    { path: '/register', name: 'register', component: RegisterPage },
    { path: '/workspace', name: 'workspace', component: WorkspacePage, meta: { requiresAuth: true } },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach(async (to) => {
  if (!to.meta.requiresAuth) {
    return
  }

  try {
    await fetchMe()
  } catch {
    return { name: 'login' }
  }
})

export default router
