import type { ProjectSnapshot, ProjectSnapshotChapter, ProjectSnapshotScene, StreamedChapterPayload } from '@/auth'
import type { LandingChapter, LandingScene } from '@/components/landing/types'

const getString = (value: unknown, fallback = '') => (typeof value === 'string' && value.trim() ? value : fallback)
const getNumber = (value: unknown, fallback = 0) => (typeof value === 'number' && Number.isFinite(value) ? value : fallback)

const chapterId = (chapter: ProjectSnapshotChapter, index: number) => getString(chapter.chapter_id ?? chapter.id, `chapter-${index + 1}`)
const sceneId = (scene: ProjectSnapshotScene, index: number) => getString(scene.scene_id ?? scene.id, `scene-${index + 1}`)

export const snapshotToLandingChapters = (snapshot: ProjectSnapshot): LandingChapter[] => {
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

      return {
        id,
        num: `CH. ${String(index).padStart(2, '0')}`,
        title: getString(chapter.title),
        open: chapterIndex === 0,
        scenes,
      }
    })

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

  return {
    id: getString(chapter.chapter_id, `streamed-chapter-${index + 1}`),
    num: `CH. ${String(chapterIndex).padStart(2, '0')}`,
    title: getString(chapter.chapter_title),
    open: index === 0,
    scenes,
  }
}

export const streamedChaptersToLandingChapters = (chapters: StreamedChapterPayload[]): LandingChapter[] => {
  const landingChapters = chapters
    .sort((a, b) => getNumber(a.chapter_index) - getNumber(b.chapter_index))
    .map(streamedChapterToLandingChapter)

  return landingChapters
}
