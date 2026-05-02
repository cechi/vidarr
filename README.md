# Vidarr 🎬

Token-based video streaming web application with directory tree navigation and built-in player.

## Features

- **Token Authentication** — Simple access token login, no registration or passwords needed
- **Directory Tree Navigation** — Multi-root tree view of your video directories (left panel)
- **Video Player** — Browser-based player with auto-play next video in directory (right panel)
- **Swipe Gestures** — Upper half: next/previous video. Lower half: ±10s seek
- **Fullscreen Support** — Immersive viewing experience
- **Video Thumbnails** — Auto-generated via ffmpeg for easy navigation
- **Resume Playback** — Remembers your position per video (localStorage)
- **Format Conversion** — One-click remux/transcode via ffmpeg for unsupported formats
- **Keyboard Shortcuts** — Arrow keys, spacebar, escape for full control

## Configuration

Create a `vidarr.yaml` config file:

```yaml
sources:
  - name: main
    access_token: "123"
    directories:
      - /mnt/data1
      - /mnt/data2
  - name: secondary
    access_token: "456"
    directories:
      - /mnt/data3
      - /mnt/data4
```

Each access token maps to one or more source directories. Users enter their token at login and see only the directories assigned to their token.

## Architecture

- **Backend:** Go (static file serving, directory scanning, ffmpeg remux/transcode, thumbnail generation)
- **Frontend:** Vanilla JS + CSS (no framework, fast and lightweight)
- **Config:** YAML-based, no database required

## Format Support

| Format | Browser Support | Action |
|--------|----------------|--------|
| MP4 (H.264) | ✅ Native | Direct play |
| WebM (VP9) | ✅ Native | Direct play |
| MKV (H.264) | ⚠️ Remux needed | Fast remux to MP4 (seconds) |
| AVI, MOV, other | ❌ Transcode needed | Full transcode via ffmpeg |

## Usage

```bash
# Install
go build -o vidarr .

# Run
./vidarr --config vidarr.yaml --port 8080

# Open
http://localhost:8080
```

## License

MIT
