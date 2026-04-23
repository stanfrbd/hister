# Changelog

## v0.13.0

### New Features

#### Semantic Search

Full vector search support via sentence-embedding models. Documents are chunked
and embedded at index time; search queries are embedded at query time and ranked
by cosine similarity. Two storage backends are supported:

- **SQLite** (default, via bundled `sqlite-vec`) zero extra infrastructure required
- **PostgreSQL** with `pgvector` auto-selected when the database is Postgres

Configure the embedding API endpoint, model, dimensions, and chunking parameters
in the new `semantic` config section. Semantic search is opt-in and off by default.
Relevance scores are shown alongside results when semantic search is active.

#### OAuth / SSO Authentication

OAuth 2.0 and OpenID Connect (OIDC) providers can now be configured as login
methods. Add one or more entries to the new `server.oauth` config section with
`client_id`, `client_secret`, `configuration_url` (for OIDC auto-discovery), or
manual `auth_url` / `token_url`, and optional `scopes`. Multiple providers can be
active at the same time alongside the built-in username/password login.

#### MCP Server

Hister now exposes a [Model Context Protocol](https://modelcontextprotocol.io/)
endpoint at `/api/mcp`, enabling LLM agents and MCP-compatible tools to search
the index directly.

#### Persistent Crawler State Management

Recursive crawl jobs (`hister index -r`) are now stored in the database and
survive interruptions. Each job gets a unique ID (auto-generated or set via
`--job-id`). Pass `--job-id <id>` without `--recursive` to resume an
interrupted crawl from exactly where it left off, including original validator
rules and visited-URL counts.

#### New Extractors

- **Wikipedia** extracts the article body and infobox, rewrites relative links, and sanitizes the output
- **GitHub project** extracts repository descriptions and README content from GitHub project pages
- **Lobste.rs** dedicated extractor for Lobste.rs story and comment pages
- **yt-dlp** extracts video metadata (title, description, channel) from video pages via yt-dlp
- **JSON-LD** surfaces structured metadata (`@type`, `headline`, description) from pages that embed JSON-LD

All extractors now expose a `Description()` method, and an extractor information
page is available at `/extractors` in the web UI.

#### OpenSearch Suggestions

The server now serves an OpenSearch suggestions endpoint (`/api/suggest`),
allowing browsers to display search-as-you-type completions when Hister is
configured as a search engine.

### Enhancements

#### Crawler Backend for All Index Operations

The `--backend` flag (and `--backend-option`) is now available on both
`hister index` (plain and `--recursive`) and `hister import-browser`, allowing
a headless Chrome/Chromium backend for JavaScript-heavy pages without running a
full recursive crawl:

```bash
hister index --backend chromedp https://example.com
hister import-browser --backend chromedp --backend-option exec_path=/usr/bin/chromium
```

Headers and cookies can also be injected per-invocation:

```bash
hister index --header "Accept-Language=en" --cookie "session=abc; Domain=example.com" https://example.com
```

Cookies use standard `Set-Cookie` format with a required `Domain` attribute.

#### CLI Search Improvements

- `--limit N` flag caps the number of results returned
- `--fields` flag selects which document fields to include in output
- `--html` flag includes raw HTML content in the output
- Paging support added to both CLI search and `list-urls`
- `list-urls` now fetches results from the server by default; `--offline` connects directly to the index without a running server

#### Quoted Field Queries

Field-qualified queries now support quoted values, enabling correct deletion and
lookup of URLs that contain spaces (common on Windows file paths):

```
url:"file:///C:/Users/My Documents/notes.txt"
```

#### Preview Panel Polish

- Preview title is now clickable (opens the result URL)
- Preview panel maximises available content width
- JSON-LD metadata surfaced inside the preview panel
- Dark theme font colors fixed in preview popup

#### NixOS / Nix Module

- `systemd` and `launchd` hardening applied to the Hister service units
- New `services.hister.environmentFile` option for secrets injection
- `openFirewall` now requires explicit opt-in
- `services.hister.config` renamed to `services.hister.settings`

#### Other

- Executable size reduced ~70 MB by switching to a trimmed `lingua-go` fork
- Sensitive content rejection errors surfaced in the browser extension
- `--verbose` flag on `hister delete` lists matched URLs before deleting
- Priority result deduplication now copies body text from the original result
- `/suggest` endpoint protected by auth middleware and `Sec-Fetch-Site` header check
- Version information included in the MCP endpoint response
- Timezone data bundled into the binary for environments without a system `tzdata`

### Bug Fixes

- File URLs (`file://`) now handled correctly in the UI for both opening and deletion (#362)
- Browser extension authentication documentation corrected (#366)
- URLs no longer lowercased during query building, preventing mismatches on case-sensitive paths
- History view correctly filtered per-user in multi-user mode (#314)
- Token authentication middleware now respects `NoAuth` flag (#348)
- Documents with no HTML content no longer attempt HTML extraction (#351)
- Extension no longer resubmits documents after a `406 Not Acceptable` response
- Priority results correctly deduplicated against standard results
- File indexing fixed on Windows
- Wide tables no longer overflow the preview panel
- Score field populated correctly in search responses
- `aws_access_key` sensitive content pattern tightened to reduce false positives
- Home-manager service units correctly gated on host platform in Nix module

## v0.12.0

### New Features

#### Web Crawler

New `hister index -r <url>` command crawls sites recursively using BFS traversal.
Supports an HTTP backend and a headless Chrome backend (chromedp).
Configurable depth, link count, allowed/excluded domains, and URL patterns.

#### PostgreSQL Backend

Full PostgreSQL support as an alternative to SQLite, including pgvector for semantic search.
Configure via a `postgres://` connection string in `server.database`.

#### Extractor Pipeline Overhaul

Extractors are now configurable, have explicit states (continue/done), and expose
a `Preview()` method used by the readability panel. New extractors included:

- Custom `pkg.go.dev` extractor for Go documentation pages
- Basic Stack Overflow extractor

#### Desktop Readability Panel

Focused search results load automatically in a split-pane reader on the right side
on screens wider than 1280 px. The panel is togglable and its open/closed state persists.

### Enhancements

- HTML sanitizer (bluemonday) applied to all extracted content
- `metadata` field added to documents for arbitrary key/value data
- `search` input type attribute on search fields for better mobile UX
- Build commit ID shown in the version string
- Admin users can create global indexes or indexes on behalf of other users
- `hister index` skips already-indexed URLs by default; pass `--force` to reindex them
- URL and domain wildcard matching automatically anchors to start and end
- Table of contents added to the API docs page
- Document indexed date shown in the preview panel
- Search query reflected in the browser tab title
- WebSocket communication optimised to reduce redundant round-trips
- Automatic redirect on zero results is now optional (configurable)
- `import` command renamed to `import-browser` to free `import` for index import/export

### Bug Fixes

- Browser history database opened read-only to avoid lock conflicts (#304)
- History entries now deleted when their associated document is deleted (#303)
- Crawler user-agent correctly applied after redirect handling (#302)
- Fixed field-specific alternation parts in query parser (#274)
- Negated query terms no longer trimmed twice
- HTML field no longer leaks into search results (#268)
- Expanded query hint only shown when the expansion is longer than the original query
- URL changes after HTTP redirects now resolved correctly
- Crawler no longer stops on HTTP errors
- Crawler timeout now applied during browser history import (#278)
- Pinned result titles no longer truncated on narrow screens
- Dark mode handled correctly in the preview panel
- Mobile layout no longer introduces unwanted line breaks
