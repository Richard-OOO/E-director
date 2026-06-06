const { createApp, ref, computed, onMounted, onUnmounted } = Vue;

const app = createApp({
    setup() {
        const lang = ref('en');
        const view = ref('new'); // 'new', 'archive', 'prompts'
        const currentStage = ref('import'); // 'import', 'processing', 'editor', 'export'
        const sidebarWidth = ref(240);
        const activeScene = ref('s1');
        const projectName = ref('Untitled Directing Project');

        const translations = {
            en: {
                nav: { new: 'New Project', archive: 'Archive', prompts: 'Prompts' },
                sidebar: { new: 'New Project', archive: 'Archive', prompts: 'Prompts' },
                viewTitles: { new: 'Project Workshop', archive: 'Historical Archive', prompts: 'AI Prompt Library' },
                import: { h2: 'Start a New Project', p: 'Import your document to begin the AI-driven directing process.', upload: 'Drop Script / Document', cta: 'Begin Analysis' },
                processing: { h2: 'Deep Reading Scenario...', p: 'Our AI is currently segmenting chapters and analyzing narrative depth.' },
                editor: { outline: 'Script Outline', scenes: 'Scenes', chapters: 'Chapters' },
                archive: { tag: 'Cinematic Work', footer: 'Directorial Proof' },
                prompts: { tag: 'System Directive', footer: 'Context: Narrative' }
            },
            zh: {
                nav: { new: '新建项目', archive: '历史归档', prompts: '提示词管理' },
                sidebar: { new: '新建项目', archive: '历史归档', prompts: '提示词管理' },
                viewTitles: { new: '项目工作台', archive: '历史作品集', prompts: 'AI 提示词库' },
                import: { h2: '开启新导演项目', p: '导入您的文档，开启 AI 助力导演流程。', upload: '拖拽脚本或文档', cta: '开始解析' },
                processing: { h2: '正在深度解析脚本...', p: 'AI 正在对章节进行切分，并分析叙事结构与深度。' },
                editor: { outline: '脚本大纲', scenes: '场景', chapters: '章节' },
                archive: { tag: '电影作品', footer: '导演样片' },
                prompts: { tag: '系统指令', footer: '上下文: 叙事' }
            }
        };

        const t = computed(() => translations[lang.value] || translations.en);
        const currentViewTitle = computed(() => t.value.viewTitles[view.value]);

        const chapters = ref([
            { id: 'c1', num: 'CH. 01', title: 'The Silent Threshold', open: true, scenes: [
                { id: 's1', num: '01', title: 'Arrival at the Gate', intent: 'Establish Isolation', comment: 'Use a tight 50mm lens. The AI will prioritize deep shadows to emphasize character solitude.' },
                { id: 's2', num: '02', title: 'The First Encounter', intent: 'Create Tension', comment: 'Slow dolly zoom during the reveal. Maintain a cold color grade.' }
            ]},
            { id: 'c2', num: 'CH. 02', title: 'Echoes of the Past', open: false, scenes: [
                { id: 's3', num: '03', title: 'Memory Corridor', intent: 'Surreal Flashback', comment: 'Soft focus edges. The AI suggests a higher frame rate (60fps).' }
            ]}
        ]);

        const archiveItems = [
            { id: 1, title: 'Obsidian Dreams', date: '2026.05.12', pattern: 'pattern-dots' },
            { id: 2, title: 'Neon Monolith', date: '2026.04.08', pattern: 'pattern-grid' },
            { id: 3, title: 'The Paper Forest', date: '2026.03.22', pattern: 'pattern-lines' },
            { id: 4, title: 'Aperture Sky', date: '2026.02.15', pattern: 'pattern-dots' }
        ];

        const prompts = [
            { id: 1, num: '01', text: 'Act as a professional cinematographer focusing on the Golden Hour aesthetics. Use high dynamic range.', tags: ['Lighting', 'Visuals'] },
            { id: 2, num: '02', text: 'Analyze character motivation through micro-expressions. Direct the AI to capture subtle eye movements.', tags: ['Performance', 'AI'] },
            { id: 3, num: '03', text: 'Implement a non-linear narrative structure. Flashbacks should be triggered by cues.', tags: ['Story', 'Logic'] }
        ];

        const toggleLang = () => lang.value = lang.value === 'en' ? 'zh' : 'en';

        const startProcessing = () => {
            currentStage.value = 'processing';
            setTimeout(() => currentStage.value = 'editor', 2000);
        };

        const scrollToScene = (id) => {
            activeScene.value = id;
            const el = document.getElementById('scene-' + id);
            if (el) el.scrollIntoView({ behavior: 'smooth', block: 'center' });
        };

        // Resizable Sidebar logic
        let isResizing = false;
        const startResize = (e) => {
            isResizing = true;
            document.body.style.cursor = 'col-resize';
            e.preventDefault();
        };
        const stopResize = () => {
            isResizing = false;
            document.body.style.cursor = 'default';
        };
        const onResize = (e) => {
            if (!isResizing) return;
            const newWidth = e.clientX;
            if (newWidth > 60 && newWidth < 500) sidebarWidth.value = newWidth;
        };

        onMounted(() => {
            window.addEventListener('mousemove', onResize);
            window.addEventListener('mouseup', stopResize);
        });

        onUnmounted(() => {
            window.removeEventListener('mousemove', onResize);
            window.removeEventListener('mouseup', stopResize);
        });

        return {
            lang, view, currentStage, sidebarWidth, activeScene, projectName,
            t, currentViewTitle, chapters, archiveItems, prompts,
            toggleLang, startProcessing, scrollToScene, startResize
        };
    }
});

app.mount('#app');
