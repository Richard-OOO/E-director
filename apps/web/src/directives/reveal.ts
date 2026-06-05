import type { Directive } from 'vue'

export const revealDirective: Directive<HTMLElement> = {
  mounted(el) {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            el.classList.add('active')
          }
        })
      },
      { threshold: 0.1 },
    )

    observer.observe(el)
  },
}
