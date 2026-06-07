<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <main>
    <section class="hero container">
      <h1>E-Director <em>Studio</em></h1>
      <p>把小说章节拆解成可编辑的场景结构、镜头提示和 YAML 生产稿，让长文本创作进入可控的导演工作流。</p>
      <RouterLink class="btn-submit hero-cta" to="/login">Enter Workspace</RouterLink>

      <div class="collage" aria-label="E-Director workflow preview">
        <div class="collage-item">
          <span>Chapter Intake</span>
        </div>
        <div class="collage-item">
          <span>Scene Mapping</span>
        </div>
        <div class="collage-item">
          <span>YAML Editor</span>
        </div>
      </div>
    </section>

    <section class="intro-section container">
      <div class="intro-layout">
        <h2>From prose to production structure.</h2>
        <p>登录后进入主工作台，可以创建项目、查看历史、管理提示词，并在生成完成后进入编辑器同步查看章节、场景和 YAML 解释。</p>
      </div>
    </section>
  </main>
</template>
