import type { ProjectSnapshot, ProjectSnapshotChapter, ProjectSnapshotScene, StreamedChapterPayload } from '@/auth'
import type { LandingChapter, LandingScene } from '@/components/landing/types'

const snapshotDebug = (label: string, payload: unknown) => {
  console.debug(`[e-director:snapshot] ${label}`, payload)
}

const summarizeSnapshot = (snapshot: ProjectSnapshot) => ({
  project_id: snapshot.project?.id,
  chapter_count: snapshot.chapters?.length ?? 0,
  scene_count: snapshot.scenes?.length ?? 0,
  top_level_keys: Object.keys(snapshot),
  first_chapter_keys: Object.keys(snapshot.chapters?.[0] ?? {}),
  first_scene_keys: Object.keys(snapshot.scenes?.[0] ?? {}),
  first_scene_yaml_lengths: {
    editable_yaml: snapshot.scenes?.[0]?.editable_yaml?.length ?? 0,
    generated_yaml: snapshot.scenes?.[0]?.generated_yaml?.length ?? 0,
    design_reason_yaml: snapshot.scenes?.[0]?.design_reason_yaml?.length ?? 0,
  },
})

const summarizeLandingChapters = (chapters: LandingChapter[]) => ({
  chapter_count: chapters.length,
  scene_count: chapters.reduce((sum, chapter) => sum + chapter.scenes.length, 0),
  yaml_scene_count: chapters.reduce((sum, chapter) => sum + chapter.scenes.filter((scene) => scene.yaml?.trim()).length, 0),
  first_chapter: chapters[0]
    ? {
        id: chapters[0].id,
        title: chapters[0].title,
        scene_count: chapters[0].scenes.length,
        first_scene_yaml_length: chapters[0].scenes[0]?.yaml?.length ?? 0,
      }
    : null,
})

const getString = (value: unknown, fallback = '') => (typeof value === 'string' && value.trim() ? value : fallback)
const getNumber = (value: unknown, fallback = 0) => (typeof value === 'number' && Number.isFinite(value) ? value : fallback)

const chapterId = (chapter: ProjectSnapshotChapter, index: number) => getString(chapter.chapter_id ?? chapter.id, `chapter-${index + 1}`)
const sceneId = (scene: ProjectSnapshotScene, index: number) => getString(scene.scene_id ?? scene.id, `scene-${index + 1}`)

export const snapshotToLandingChapters = (snapshot: ProjectSnapshot): LandingChapter[] => {
  snapshotDebug('raw project snapshot', summarizeSnapshot(snapshot))
  const chapters = snapshot.chapters ?? []
  const flatScenes = snapshot.scenes ?? []

  const landingChapters = chapters
    .map((chapter, chapterIndex) => {
      const id = chapterId(chapter, chapterIndex)
      const index = getNumber(chapter.chapter_index, chapterIndex + 1)
      const scenes = flatScenes
        .filter((scene) => getString(scene.chapter_id) === id)
        .sort((a, b) => getNumber(a.scene_index) - getNumber(b.scene_index))
        .map<LandingScene>((scene, sceneIndex) => ({
          id: sceneId(scene, sceneIndex),
          num: String(getNumber(scene.scene_index, sceneIndex + 1)).padStart(2, '0'),
          title: getString(scene.title),
          intent: getString(scene.summary ?? scene.status),
          yaml: getString(scene.editable_yaml ?? scene.generated_yaml),
        }))
        .filter((scene) => getString(scene.yaml))

      return {
        id,
        num: `CH. ${String(index).padStart(2, '0')}`,
        title: getString(chapter.title),
        open: chapterIndex === 0,
        scenes,
      }
    })
    .filter((chapter) => chapter.scenes.length > 0)

  snapshotDebug('mapped project snapshot', summarizeLandingChapters(landingChapters))
  return landingChapters
}

export const streamedChapterToLandingChapter = (chapter: StreamedChapterPayload, index: number): LandingChapter => {
  const chapterIndex = getNumber(chapter.chapter_index, index + 1)
  const scenes = (chapter.scenes ?? [])
    .map<LandingScene>((scene, sceneIndex) => ({
      id: getString(scene.scene_id, `streamed-scene-${sceneIndex + 1}`),
      num: String(getNumber(scene.scene_index, sceneIndex + 1)).padStart(2, '0'),
      title: getString(scene.title),
      intent: getString(scene.summary),
      yaml: getString(scene.yaml_content),
    }))
    .filter((scene) => getString(scene.yaml))

  return {
    id: getString(chapter.chapter_id, `streamed-chapter-${index + 1}`),
    num: `CH. ${String(chapterIndex).padStart(2, '0')}`,
    title: getString(chapter.chapter_title),
    open: index === 0,
    scenes,
  }
}

export const streamedChaptersToLandingChapters = (chapters: StreamedChapterPayload[]): LandingChapter[] => {
  snapshotDebug('raw streamed chapters', {
    chapter_count: chapters.length,
    chapters: chapters.map((chapter) => ({
      chapter_id: chapter.chapter_id,
      chapter_index: chapter.chapter_index,
      scene_count: chapter.scenes?.length ?? 0,
      yaml_lengths: chapter.scenes?.map((scene) => scene.yaml_content?.length ?? 0) ?? [],
    })),
  })

  const landingChapters = chapters
    .filter((chapter) => (chapter.scenes ?? []).some((scene) => getString(scene.yaml_content)))
    .sort((a, b) => getNumber(a.chapter_index) - getNumber(b.chapter_index))
    .map(streamedChapterToLandingChapter)

  snapshotDebug('mapped streamed chapters', summarizeLandingChapters(landingChapters))
  return landingChapters
}
