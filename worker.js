// Cloudflare Worker: serves the SPA from static assets and reverse-proxies
// /api/* to the external Go backend (Workers cannot run Go/Postgres/Redis).
// Set API_ORIGIN (e.g. https://riftstats-api.fly.dev) in the Worker vars.
export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname.startsWith('/api/')) {
      const origin = (env.API_ORIGIN || '').replace(/\/+$/, '');
      if (!origin) {
        return Response.json(
          { type: 'https://api.riftstats.bn.io/errors/bad-gateway', title: 'Bad Gateway', status: 502, detail: 'API_ORIGIN is not configured on this Worker.' },
          { status: 502 }
        );
      }
      const target = origin + url.pathname + url.search;
      const init = {
        method: request.method,
        headers: request.headers,
        redirect: 'manual',
      };
      if (request.method !== 'GET' && request.method !== 'HEAD') {
        init.body = request.body;
        init.duplex = 'half';
      }
      return fetch(target, init);
    }

    return env.ASSETS.fetch(request);
  },
};
