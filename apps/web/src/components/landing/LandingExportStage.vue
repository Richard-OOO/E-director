<script setup lang="ts">
const props = defineProps<{
  activeProjectId: string
  isExporting: boolean
  exportingMode: 'batch' | 'combined' | ''
  exportStatus: string
  canExportYaml: boolean
}>()

const emit = defineEmits<{
  (event: 'export-yaml', mode: 'batch' | 'combined'): void
}>()
</script>

<template>
  <section class="stage-view">
    <div class="export-panel export-panel-single">
      <span class="meta">Delivery</span>
      <h2>YAML Export</h2>
      <p class="export-lead">Choose whether to download scene YAML as separate files or one combined YAML document.</p>
      <div class="export-grid export-grid-dual">
        <div class="export-card export-card-primary">
          <span class="prompt-tag">Batch YAML</span>
          <h3>分批次导出</h3>
          <p>Download a zip that contains one YAML file per generated scene. Edited YAML is used first, then generated YAML as fallback.</p>
          <button class="btn-submit" type="button" :disabled="props.isExporting || !props.canExportYaml" @click="emit('export-yaml', 'batch')">
            {{ props.isExporting && props.exportingMode === 'batch' ? 'Exporting...' : 'Export YAML Zip' }}
          </button>
        </div>

        <div class="export-card export-card-primary">
          <span class="prompt-tag">Combined YAML</span>
          <h3>合并导出</h3>
          <p>Download one YAML file that concatenates all generated scene YAML blocks in chapter and scene order, separated by YAML document markers.</p>
          <button class="btn-submit" type="button" :disabled="props.isExporting || !props.canExportYaml" @click="emit('export-yaml', 'combined')">
            {{ props.isExporting && props.exportingMode === 'combined' ? 'Exporting...' : 'Export Combined YAML' }}
          </button>
        </div>
      </div>
      <p v-if="!props.activeProjectId" class="export-status muted">Open or generate a project before exporting.</p>
      <p v-else-if="!props.canExportYaml" class="export-status muted">No YAML scenes are ready to export yet.</p>
      <p v-if="props.exportStatus" class="export-status">{{ props.exportStatus }}</p>
    </div>
  </section>
</template>
