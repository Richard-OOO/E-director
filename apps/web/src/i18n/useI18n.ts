import { computed, ref } from 'vue'

type Language = 'en' | 'zh'

const lang = ref<Language>('en')

export const translations = {
  en: {
    nav: { home: 'Home', login: 'Login', register: 'Register', workspace: 'Workspace', logout: 'Logout' },
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
      p: 'Use your email and password to access your director dashboard.',
      email: 'Email Address',
      password: 'Password',
      submit: 'Sign In',
      footer: 'New to E-Director?',
      registerLink: 'Create an account',
      submiting: 'Signing In...',
      success: 'Signed in. Opening your workspace...',
    },
    register: {
      h2: 'Create Account',
      p: 'Verify your email to join our community of AI-powered filmmakers today.',
      name: 'Full Name',
      email: 'Email Address',
      password: 'Password',
      verificationCode: 'Verification Code',
      sendCode: 'Send Code',
      sendingCode: 'Sending...',
      resendIn: 'Resend in',
      codeSent: 'Verification code sent. Please check your email.',
      submit: 'Get Started',
      footer: 'Already have an account?',
      loginLink: 'Sign in instead',
      submiting: 'Creating Account...',
      success: 'Account created. Opening your workspace...',
    },
    workspace: {
      eyebrow: 'Director Workspace',
      h1: 'Your production desk is ready.',
      p: 'Auth is connected. Next we will wire novel import, generation jobs, and Python AI orchestration here.',
      loading: 'Checking your session...',
      signedInAs: 'Signed in as',
      logout: 'Logout',
    },
    footer: '© 2026 E-DIRECTOR STUDIO. ALL RIGHTS RESERVED.',
  },
  zh: {
    nav: { home: '首页', login: '登录', register: '注册', workspace: '工作台', logout: '退出' },
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
      p: '使用邮箱和密码访问你的导演控制面板。',
      email: '电子邮箱',
      password: '密码',
      submit: '登录',
      footer: '第一次来到 E-Director？',
      registerLink: '创建账号',
      submiting: '正在登录...',
      success: '登录成功，正在打开工作台...',
    },
    register: {
      h2: '创建账号',
      p: '验证邮箱后加入我们的 AI 驱动电影制作人社区。',
      name: '姓名',
      email: '电子邮箱',
      password: '密码',
      verificationCode: '验证码',
      sendCode: '发送验证码',
      sendingCode: '发送中...',
      resendIn: '重新发送还需',
      codeSent: '验证码已发送，请查收邮箱。',
      submit: '开始导演',
      footer: '已经有账号了？',
      loginLink: '直接登录',
      submiting: '正在创建账号...',
      success: '账号已创建，正在打开工作台...',
    },
    workspace: {
      eyebrow: '导演工作台',
      h1: '你的创作控制台已就绪。',
      p: 'Auth 已完成连接。下一步会在这里接入小说导入、生成任务和 Python AI 编排。',
      loading: '正在校验会话...',
      signedInAs: '当前登录',
      logout: '退出登录',
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
