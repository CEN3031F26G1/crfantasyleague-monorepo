# Royale Fantasy

A thin React + TypeScript + Vite frontend for the CRL 2026 fantasy league. Radix primitives provide the profile menu, avatars, qualification select, and accessible player dialogs. Styling is plain CSS; fonts and player images are served locally. No component megaframework, database, or client-side token storage.

## Run locally

Requirements: Node.js **22.12+** (or a current supported version), npm, and the Go version specified in `../backend/go.mod`.

1. Start the existing backend in one terminal:

   ```sh
   cd backend  # from the repository root
   # Keep your existing .env and private.pem. If absent, copy .env.example
   # to .env, fill in the Discord secret, and generate an RSA private key:
   # openssl genrsa -out private.pem 2048
   go run .
   ```

2. Start the frontend in a second terminal:

   ```sh
   cd frontend  # from the repository root
   npm ci
   npm run dev
   ```

3. Visit **http://localhost:5173**. Use `localhost` consistently, not `127.0.0.1`, because the OAuth state and session cookies must remain on the same host.

Vite proxies `/api/*` and `/auth/*` to `http://localhost:8080`. Port 5173 is strict so a busy port does not silently break the configured OAuth return URL. No frontend `.env` is needed, and no Discord secret belongs in the browser.

### Discord application settings

Register **http://localhost:8080/auth/callback** as a redirect in the Discord Developer Portal for the configured client ID. `backend/.env.example` documents all settings. Existing `.env` values are not overwritten. A valid Discord application secret and matching registered redirect are required for actual sign-in.

Flow:

1. `/login` is the React sign-in page.
2. `GET /auth/login` starts the backend's Discord OAuth flow with a random, short-lived, HttpOnly state cookie and only the `identify` scope.
3. `GET /auth/callback` validates state, exchanges the code, fetches the Discord profile, and sets a **12-hour HttpOnly SameSite=Lax** session cookie before redirecting to the frontend. No JWT is placed in a URL or local/session storage. The original backend `GET /login` entry point remains an alias.
4. `GET /api/me` returns ID, display name, username, avatar URL, and expiration. Existing bearer-token API clients still work.
5. `POST /auth/logout` clears the cookie and rejects cross-origin browser requests. The profile menu becomes a sign-in button; the catalog stays public.

The frontend checks the session on load, window focus, and once per minute, and clears the local profile at session expiry. Failed session checks do not prevent public browsing. Sessions are stateless: logout removes this browser's cookie, but does not revoke a previously copied bearer token before its expiry or disconnect the application from Discord. There is no persistent user database yet.

## Routes and scope

- `/`: minimal responsive landing page, event facts, and public catalog CTAs.
- `/players`: all 16 sourced qualifiers, name/country search, qualification filters, and player profiles. Filters are shareable through `?q=Japan&route=season`.
- `/login`: Discord sign-in and friendly OAuth cancellation/failure handling.
- Other routes: a small 404 page.

Logging in changes only the account controls. Team drafting, scoring, payments, betting, and league management are intentionally **not** implemented; the UI labels future fantasy features as coming soon. This is a fantasy fan site, not a real-money wagering implementation.

## Player data and images

The catalog is a curated **October 30, 2026** snapshot, not a live API feed. All 16 qualification routes and countries come from Liquipedia's 2026 World Finals participants list; biographies use individual player pages where available. See [`SOURCES.md`](SOURCES.md) for links, provenance, and remaining image gaps.

Nine players have locally served profile images from linked public player channels/accounts. Some players use artwork or logos as their profile image rather than a headshot. Seven players use explicit monogram fallbacks because a verified image could not be retrieved. No stock people or invented player portraits are used. Sources and image credits are also accessible inside each player dialog. Image rights remain with the original creators; verify usage permissions before a public/commercial launch.

To update the catalog, edit `src/data/players.json`, put approved images in `public/players/`, and update the verification date in `src/components.tsx` and source documentation. Keep `image` and `imageSource` null when no verified image is available.

## Build and checks

```sh
npm run build        # TypeScript + production assets in dist/
npm test             # Vitest + Testing Library
npm audit

# From the repository root:
cd backend
go test -race ./...
go vet ./...
```

Frontend tests cover public access, the 16-player roster, search/reset, URL filters, accessible profile dialogs, OAuth entry/error handling, profile menu/logout, and backend failure fallback. Backend tests mock Discord to exercise the OAuth exchange, state validation, cancellation, session cookies, bearer compatibility, issuer/expiry/algorithm validation, and logout origin checks. A real Discord account authorization is still a manual integration check.

## Production

Build with `npm run build` and serve `dist/` with SPA fallback to `index.html`. Reverse-proxy `/api/*` and `/auth/*` to the Go service **before** the SPA fallback. Keep the browser-facing API on the same origin; permissive CORS is not required.

Set, for example:

```dotenv
FRONTEND_URL=https://fantasy.example.com
DISCORD_REDIRECT_URL=https://fantasy.example.com/auth/callback
```

Register the production callback with Discord. The backend automatically enables Secure cookies when the callback URL is HTTPS. Terminate TLS at your reverse proxy, keep the private key and client secret server-side, and preserve the request Origin header. Avoid cross-site frontend/API deployments without redesigning the cookie and CSRF policy. `npm run preview` is only a static build preview, not the configured authenticated development server.

Design references RoyaleTools' compact tool-oriented navigation and card layouts without copying its branding, source code, or art.
