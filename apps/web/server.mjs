import { createReadStream, existsSync } from 'node:fs'
import { extname, join, normalize, resolve } from 'node:path'
import { createServer } from 'node:http'

const port = Number(process.env.PORT || 8080)
const root = resolve(process.env.WEB_ROOT || './dist')

const contentTypes = {
  '.css': 'text/css; charset=utf-8',
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.ico': 'image/x-icon',
}

createServer((request, response) => {
  const url = new URL(request.url || '/', `http://${request.headers.host || 'localhost'}`)
  const requestedPath = normalize(decodeURIComponent(url.pathname)).replace(/^([/\\])+/, '')
  const filePath = resolve(join(root, requestedPath || 'index.html'))
  const safePath = filePath.startsWith(root) && existsSync(filePath) ? filePath : join(root, 'index.html')
  const contentType = contentTypes[extname(safePath)] || 'application/octet-stream'

  response.writeHead(200, { 'Content-Type': contentType })
  createReadStream(safePath).pipe(response)
}).listen(port, '0.0.0.0', () => {
  console.log(`Web app listening on http://0.0.0.0:${port}`)
})
