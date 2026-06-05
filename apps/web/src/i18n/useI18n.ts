import { computed, ref } from 'vue'

type Language = 'en' | 'zh'

const lang = ref<Language>('en')

export const translations = {
  en: {
    nav: { home: 'Home', login: 'Login', register: 'Register' },
    hero: {
      h1: 'AI Powers<br>Directing <em>Everything.</em>',
      p: 'E-Director is an editorial-grade platform designed to help you orchestrate films, visual stories, and narratives through the power of artificial intelligence.',
      cta: 'Get Started',
    },
    collage: { vision: 'The Vision', motion: 'The Motion', result: 'The Result' },
    section: {
      h2: 'Innovating the Art of Film',
      p: 'We bridge the gap between creative vision and technical execution, allowing anyone to direct like a pro.',
    },
    login: {
      h2: 'Welcome Back',
      p: 'Enter your credentials to access your director dashboard.',
      email: 'Email Address',
      password: 'Password',
      submit: 'Sign In',
      footer: 'New to E-Director?',
      registerLink: 'Create an account',
    },
    register: {
      h2: 'Create Account',
      p: 'Join our community of AI-powered filmmakers today.',
      name: 'Full Name',
      email: 'Email Address',
      password: 'Password',
      submit: 'Get Started',
      footer: 'Already have an account?',
      loginLink: 'Sign in instead',
    },
    footer: '© 2026 E-DIRECTOR STUDIO. ALL RIGHTS RESERVED.',
  },
  zh: {
    nav: { home: '首页', login: '登录', register: '注册' },
    hero: {
      h1: 'AI 助力<br>导演 <em>一切.</em>',
      p: 'E-Director 是一个专业级的平台，旨在帮助您通过人工智能的力量编排电影、视觉故事和叙事。',
      cta: '开始导演',
    },
    collage: { vision: '视觉', motion: '动态', result: '成果' },
    section: {
      h2: '创新电影艺术',
      p: '我们弥合了创意视觉与技术执行之间的差距，让任何人都能像专业人士一样导演。',
    },
    login: {
      h2: '欢迎回来',
      p: '输入您的凭据以访问您的导演控制面板。',
      email: '电子邮箱',
      password: '密码',
      submit: '登录',
      footer: '第一次来到 E-Director？',
      registerLink: '创建账号',
    },
    register: {
      h2: '创建账号',
      p: '今天就加入我们的 AI 驱动电影制作人社区。',
      name: '姓名',
      email: '电子邮箱',
      password: '密码',
      submit: '开始导演',
      footer: '已经有账号了？',
      loginLink: '直接登录',
    },
    footer: '© 2026 E-DIRECTOR 工作室。保留所有权利。',
  },
} as const

export function useI18n() {
  const t = computed(() => translations[lang.value])

  const toggleLang = () => {
    lang.value = lang.value === 'en' ? 'zh' : 'en'
  }

  return { lang, t, toggleLang }
}
