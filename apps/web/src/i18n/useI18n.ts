import { computed, ref } from 'vue'

type Language = 'en' | 'zh'

const lang = ref<Language>('en')

export const translations = {
  en: {
    nav: { home: 'Home', login: 'Login', register: 'Register', workspace: 'Workspace', logout: 'Logout' },
    layout: {
      sidebar: { new: 'New Project', archive: 'Archive', prompts: 'Prompts' },
      viewTitles: { new: 'Project Workshop', archive: 'Historical Archive', prompts: 'AI Prompt Library' },
      flow: {
        import: 'Import',
        processing: 'Processing',
        editor: 'Editor',
        export: 'Export',
        project: 'Project Workshop',
        archive: 'Historical Archive',
        prompts: 'AI Prompt Library',
      },
    },
    login: {
      h2: 'Welcome Back',
      p: 'Sign in with your password to continue.',
      submiting: 'Signing in...',
      submit: 'Sign In',
      success: 'Logged in successfully.',
      email: 'Email',
      password: 'Password',
      footer: 'Need an account?',
      registerLink: 'Register now',
    },
    register: {
      h2: 'Create Account',
      p: 'Register with your email verification code.',
      submiting: 'Creating account...',
      submit: 'Create Account',
      sendingCode: 'Sending code...',
      resendIn: 'Resend in {seconds}s',
      sendCode: 'Send Code',
      codeSent: 'Verification code sent.',
      success: 'Account created successfully.',
      name: 'Name',
      email: 'Email',
      password: 'Password',
      verificationCode: 'Verification Code',
      footer: 'Already have an account?',
      loginLink: 'Login now',
    },
    workspace: {
      eyebrow: 'Workspace',
      h1: 'Your project is ready.',
      p: 'Open the workspace to continue editing, reviewing, or exporting your screenplay project.',
      panel: { label: 'Current session', title: 'Server-backed login state', note: 'Cookie session is the source of truth.' },
      loading: 'Loading session...',
      signedInAs: 'Signed in as',
      logout: 'Sign Out',
      cta: 'Enter Workspace',
    },
    import: {
      h2: 'Start a New Project',
      p: 'Import your document to begin the AI-driven directing process.',
      upload: 'Drop Script / Document',
      cta: 'Begin Analysis',
    },
    processing: { h2: 'Deep Reading Scenario...', p: 'Our AI is currently segmenting chapters and analyzing narrative depth.' },
    editor: { outline: 'Script Outline', scenes: 'Scenes', chapters: 'Chapters' },
    archive: { tag: 'Cinematic Work', footer: 'Directorial Proof' },
    prompts: { tag: 'System Directive', footer: 'Context: Narrative' },
    footer: '© 2026 E-DIRECTOR STUDIO. ALL RIGHTS RESERVED.',
  },
  zh: {
    nav: { home: '首页', login: '登录', register: '注册', workspace: '工作台', logout: '退出' },
    layout: {
      sidebar: { new: '新建项目', archive: '历史归档', prompts: '提示词管理' },
      viewTitles: { new: '项目工作台', archive: '历史作品集', prompts: 'AI 提示词库' },
      flow: {
        import: '导入',
        processing: '处理中',
        editor: '编辑器',
        export: '导出',
        project: '项目工作台',
        archive: '历史作品集',
        prompts: 'AI 提示词库',
      },
    },
    login: {
      h2: '欢迎回来',
      p: '使用密码登录后继续。',
      submiting: '登录中...',
      submit: '登录',
      success: '登录成功。',
      email: '邮箱',
      password: '密码',
      footer: '还没有账号？',
      registerLink: '立即注册',
    },
    register: {
      h2: '创建账号',
      p: '使用邮箱验证码完成注册。',
      submiting: '注册中...',
      submit: '创建账号',
      sendingCode: '发送中...',
      resendIn: '{seconds} 秒后重发',
      sendCode: '发送验证码',
      codeSent: '验证码已发送。',
      success: '注册成功。',
      name: '姓名',
      email: '邮箱',
      password: '密码',
      verificationCode: '验证码',
      footer: '已有账号？',
      loginLink: '立即登录',
    },
    workspace: {
      eyebrow: '工作台',
      h1: '你的项目已就绪。',
      p: '进入工作台继续编辑、审阅或导出你的剧本项目。',
      panel: { label: '当前会话', title: '服务端登录态', note: 'Cookie session 是唯一真实登录态。' },
      loading: '正在加载会话...',
      signedInAs: '当前登录用户',
      logout: '退出登录',
      cta: '进入工作台',
    },
    import: {
      h2: '开启新项目',
      p: '导入文档，开始 AI 驱动的导演流程。',
      upload: '拖拽脚本 / 文档',
      cta: '开始解析',
    },
    processing: { h2: '正在深度解析...', p: 'AI 正在对章节进行切分并分析叙事深度。' },
    editor: { outline: '脚本大纲', scenes: '场景', chapters: '章节' },
    archive: { tag: '电影作品', footer: '导演样片' },
    prompts: { tag: '系统指令', footer: '上下文：叙事' },
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
