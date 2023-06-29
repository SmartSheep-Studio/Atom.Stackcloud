export const placeholders: { [id: string]: string } = {
  description: "No description yet",
  banner: "https://static.smartsheep.studio/d/Atom/Placeholders/default-banner.jpg",
  avatar: "https://static.smartsheep.studio/d/Atom/Icon.png",
  project: "mdi-sitemap",
}

export function usePlaceholder(id: string, val?: string): string {
  if (val == null || val.length <= 0) {
    return placeholders[id]
  } else {
    return val
  }
}
