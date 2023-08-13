export function parseRedirect(query: any, fallback: any = { name: "console" }): Promise<any> {
  return new Promise((resolve, reject) => {
    if (query.redirect_uri != null) {
      window.location.href = query.redirect_uri
      setTimeout(() => {
        resolve(fallback)
      }, 10000)
    } else {
      resolve(fallback)
    }
  })
}
