# gitmap help-dashboard

Serve the interactive documentation site locally in your browser.

**Alias:** `hd`

## Usage

```
gitmap help-dashboard [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--port <number>` | Port to serve on (default: 5173) |

## How It Works

1. Locates the `docs-site/` directory relative to the gitmap binary
2. If `docs-site/` is missing but `docs-site.zip` exists, extracts it automatically
3. If a pre-built `dist/` folder exists, serves it with a built-in HTTP server
4. If no `dist/` found, falls back to `npm install && npm run dev`
5. Opens the dashboard in your default browser automatically

## Prerequisites

- **Static mode**: No dependencies — serves pre-built files directly
- **Auto-extract mode**: `docs-site.zip` is downloaded by the installer and extracted on first run
- **Dev mode (fallback)**: Requires Node.js and npm on PATH

## Examples

    $ gitmap help-dashboard

    Serving docs from /usr/local/bin/docs-site/dist on http://localhost:5173
    Opening http://localhost:5173 in browser...

    $ gitmap hd --port 8080

    No pre-built dist/ found, falling back to npm run dev
    Running npm install...
    Starting dev server from /usr/local/bin/docs-site...
    Opening http://localhost:8080 in browser...

Press Ctrl+C to stop the server.

## See Also

- docs — Open the hosted documentation website
- dashboard — Generate an HTML analytics dashboard for a repo
