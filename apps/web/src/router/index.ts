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
    { path: '/workspace', redirect: '/workspace/create', meta: { requiresAuth: true } },
    { path: '/workspace/create', name: 'workspace-create', component: WorkspacePage, meta: { requiresAuth: true } },
    { path: '/workspace/progress', name: 'workspace-progress', component: WorkspacePage, meta: { requiresAuth: true } },
    { path: '/workspace/editor', name: 'workspace-editor', component: WorkspacePage, meta: { requiresAuth: true } },
    { path: '/workspace/export', name: 'workspace-export', component: WorkspacePage, meta: { requiresAuth: true } },
    { path: '/workspace/history', name: 'workspace-history', component: WorkspacePage, meta: { requiresAuth: true } },
    { path: '/workspace/prompts', name: 'workspace-prompts', component: WorkspacePage, meta: { requiresAuth: true } },
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
