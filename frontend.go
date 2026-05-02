package main

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
<title>Vidarr</title>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --bg: #0f0f13;
  --bg-panel: #16161d;
  --bg-hover: #1e1e28;
  --bg-active: #262633;
  --border: #2a2a3a;
  --text: #e0e0e8;
  --text-dim: #888898;
  --accent: #6e8efb;
  --accent-hover: #8ba4fc;
  --accent-glow: rgba(110, 142, 251, 0.15);
  --danger: #e74c3c;
  --success: #2ecc71;
  --radius: 6px;
  --sidebar-width: 320px;
}

html, body {
  height: 100%;
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  font-size: 14px;
  overflow: hidden;
}

/* ─── Login ────────────────────────────────────────── */
#login-screen {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: linear-gradient(135deg, #0f0f13 0%, #1a1a2e 100%);
}

.login-box {
  background: var(--bg-panel);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 40px;
  width: 360px;
  text-align: center;
}

.login-box h1 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.login-box p {
  color: var(--text-dim);
  margin-bottom: 24px;
  font-size: 13px;
}

.login-box input {
  width: 100%;
  padding: 12px 16px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text);
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.login-box input:focus {
  border-color: var(--accent);
}

.login-box button {
  width: 100%;
  padding: 12px;
  margin-top: 16px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: var(--radius);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.login-box button:hover { background: var(--accent-hover); }

.login-error {
  color: var(--danger);
  margin-top: 12px;
  font-size: 13px;
  min-height: 20px;
}

/* ─── App Layout ───────────────────────────────────── */
#app {
  display: none;
  height: 100%;
}

#app.active { display: flex; }
#login-screen.hidden { display: none; }

.sidebar {
  width: var(--sidebar-width);
  min-width: var(--sidebar-width);
  background: var(--bg-panel);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.sidebar-header h2 {
  font-size: 16px;
  font-weight: 600;
}

.btn-logout {
  background: none;
  border: 1px solid var(--border);
  color: var(--text-dim);
  padding: 4px 10px;
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.btn-logout:hover {
  border-color: var(--danger);
  color: var(--danger);
}

.tree-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px 0;
}

.tree-container::-webkit-scrollbar { width: 6px; }
.tree-container::-webkit-scrollbar-track { background: transparent; }
.tree-container::-webkit-scrollbar-thumb { background: var(--border); border-radius: 3px; }

.tree-item {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
  transition: background 0.15s;
  font-size: 13px;
}

.tree-item:hover { background: var(--bg-hover); }
.tree-item.active { background: var(--accent-glow); color: var(--accent); }

.tree-arrow {
  width: 16px;
  height: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-right: 4px;
  flex-shrink: 0;
  font-size: 10px;
  color: var(--text-dim);
  transition: transform 0.2s;
}

.tree-arrow.open { transform: rotate(90deg); }
.tree-arrow.empty { visibility: hidden; }

.tree-icon {
  margin-right: 8px;
  flex-shrink: 0;
  font-size: 14px;
}

.tree-label {
  overflow: hidden;
  text-overflow: ellipsis;
}

.tree-children { display: none; }
.tree-children.open { display: block; }

/* ─── Main Content ─────────────────────────────────── */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  position: relative;
}

/* ─── File List ────────────────────────────────────── */
.file-list-bar {
  border-bottom: 1px solid var(--border);
  background: var(--bg-panel);
  flex-shrink: 0;
  max-height: 200px;
  overflow-y: auto;
}

.file-list-bar::-webkit-scrollbar { width: 4px; }
.file-list-bar::-webkit-scrollbar-thumb { background: var(--border); border-radius: 2px; }

.file-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  transition: background 0.15s;
  gap: 10px;
  font-size: 13px;
}

.file-item:hover { background: var(--bg-hover); }
.file-item.active { background: var(--accent-glow); }

.file-thumb {
  width: 48px;
  height: 32px;
  border-radius: 3px;
  background: var(--bg);
  object-fit: cover;
  flex-shrink: 0;
}

.file-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-item.active .file-name { color: var(--accent); font-weight: 500; }

.file-badge {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 3px;
  font-weight: 600;
  text-transform: uppercase;
  flex-shrink: 0;
}

.file-badge.remux { background: #3d3510; color: #f1c40f; }
.file-badge.convert { background: #3d1010; color: #e74c3c; }

.file-size {
  color: var(--text-dim);
  font-size: 12px;
  flex-shrink: 0;
}

/* ─── Player ───────────────────────────────────────── */
.player-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000;
  position: relative;
  overflow: hidden;
}

.player-wrapper video {
  max-width: 100%;
  max-height: 100%;
  outline: none;
}

.player-placeholder {
  color: var(--text-dim);
  text-align: center;
  font-size: 16px;
  user-select: none;
}

.player-placeholder span { font-size: 48px; display: block; margin-bottom: 16px; opacity: 0.3; }

/* ─── Player Controls ──────────────────────────────── */
.player-controls {
  display: none;
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.85));
  padding: 40px 16px 12px;
  flex-direction: column;
  gap: 8px;
  z-index: 10;
}

.player-controls.visible { display: flex; }

.progress-bar-container {
  width: 100%;
  height: 4px;
  background: rgba(255,255,255,0.2);
  border-radius: 2px;
  cursor: pointer;
  position: relative;
}

.progress-bar-container:hover { height: 6px; }

.progress-bar-fill {
  height: 100%;
  background: var(--accent);
  border-radius: 2px;
  width: 0%;
  position: relative;
}

.progress-bar-fill::after {
  content: '';
  position: absolute;
  right: -5px;
  top: 50%;
  transform: translateY(-50%);
  width: 12px;
  height: 12px;
  background: var(--accent);
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.15s;
}

.progress-bar-container:hover .progress-bar-fill::after { opacity: 1; }

.controls-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ctrl-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  font-size: 18px;
  padding: 4px;
  opacity: 0.9;
  transition: opacity 0.15s;
  display: flex;
  align-items: center;
}

.ctrl-btn:hover { opacity: 1; }

.time-display {
  color: rgba(255,255,255,0.7);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  user-select: none;
}

.volume-slider {
  width: 80px;
  height: 3px;
  -webkit-appearance: none;
  appearance: none;
  background: rgba(255,255,255,0.3);
  border-radius: 2px;
  outline: none;
  cursor: pointer;
}

.volume-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #fff;
}

.controls-spacer { flex: 1; }

/* ─── Conversion Overlay ───────────────────────────── */
.convert-overlay {
  display: none;
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.9);
  z-index: 20;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.convert-overlay.visible { display: flex; }

.convert-overlay h3 { font-size: 18px; font-weight: 500; }

.convert-progress-bar {
  width: 300px;
  max-width: 80%;
  height: 6px;
  background: rgba(255,255,255,0.15);
  border-radius: 3px;
  overflow: hidden;
}

.convert-progress-fill {
  height: 100%;
  background: var(--accent);
  border-radius: 3px;
  width: 0%;
  transition: width 0.3s;
}

.convert-text { color: var(--text-dim); font-size: 13px; }

.convert-actions { display: flex; gap: 12px; margin-top: 8px; }

.convert-btn {
  padding: 8px 20px;
  border: 1px solid var(--border);
  background: none;
  color: var(--text);
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.convert-btn:hover { border-color: var(--accent); color: var(--accent); }

.convert-btn.primary {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}

.convert-btn.primary:hover { background: var(--accent-hover); }

/* ─── Touch Hint ───────────────────────────────────── */
.touch-hint {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0,0,0,0.7);
  color: #fff;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 15;
}

.touch-hint.show { opacity: 1; }

/* ─── Mobile ───────────────────────────────────────── */
.sidebar-toggle {
  display: none;
  position: fixed;
  top: 12px;
  left: 12px;
  z-index: 100;
  background: var(--bg-panel);
  border: 1px solid var(--border);
  color: var(--text);
  width: 36px;
  height: 36px;
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 18px;
  align-items: center;
  justify-content: center;
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 50;
    transform: translateX(-100%);
    transition: transform 0.3s;
    width: 280px;
    min-width: 280px;
  }

  .sidebar.open { transform: translateX(0); }

  .sidebar-toggle { display: flex; }

  .sidebar-backdrop {
    display: none;
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.5);
    z-index: 49;
  }

  .sidebar-backdrop.open { display: block; }

  .file-list-bar { max-height: 150px; }
}

/* ─── Fullscreen ───────────────────────────────────── */
.player-wrapper:fullscreen { background: #000; }
.player-wrapper:fullscreen video { max-width: 100vw; max-height: 100vh; }
.player-wrapper:-webkit-full-screen { background: #000; }
</style>
</head>
<body>

<!-- Login Screen -->
<div id="login-screen">
  <div class="login-box">
    <h1>Vidarr</h1>
    <p>Enter your access token to continue</p>
    <input type="password" id="token-input" placeholder="Access token" autocomplete="off" />
    <button id="login-btn">Sign In</button>
    <div class="login-error" id="login-error"></div>
  </div>
</div>

<!-- Main App -->
<div id="app">
  <button class="sidebar-toggle" id="sidebar-toggle">&#9776;</button>
  <div class="sidebar-backdrop" id="sidebar-backdrop"></div>

  <aside class="sidebar" id="sidebar">
    <div class="sidebar-header">
      <h2>Vidarr</h2>
      <button class="btn-logout" id="logout-btn">Logout</button>
    </div>
    <div class="tree-container" id="tree-container"></div>
  </aside>

  <div class="main-content">
    <div class="file-list-bar" id="file-list"></div>

    <div class="player-wrapper" id="player-wrapper">
      <div class="player-placeholder" id="player-placeholder">
        <span>&#9654;</span>
        Select a directory and video to begin
      </div>

      <video id="video-player" style="display:none" playsinline></video>

      <div class="player-controls" id="player-controls">
        <div class="progress-bar-container" id="progress-bar">
          <div class="progress-bar-fill" id="progress-fill"></div>
        </div>
        <div class="controls-row">
          <button class="ctrl-btn" id="btn-prev" title="Previous">&#9198;</button>
          <button class="ctrl-btn" id="btn-play" title="Play/Pause">&#9654;</button>
          <button class="ctrl-btn" id="btn-next" title="Next">&#9197;</button>
          <span class="time-display" id="time-display">0:00 / 0:00</span>
          <button class="ctrl-btn" id="btn-mute" title="Mute">&#128264;</button>
          <input type="range" class="volume-slider" id="volume-slider" min="0" max="1" step="0.05" value="1" />
          <div class="controls-spacer"></div>
          <button class="ctrl-btn" id="btn-fullscreen" title="Fullscreen">&#x26F6;</button>
        </div>
      </div>

      <div class="touch-hint" id="touch-hint"></div>

      <div class="convert-overlay" id="convert-overlay">
        <h3 id="convert-title">Converting...</h3>
        <div class="convert-progress-bar">
          <div class="convert-progress-fill" id="convert-progress-fill"></div>
        </div>
        <div class="convert-text" id="convert-text">0%</div>
        <div class="convert-actions" id="convert-actions">
          <button class="convert-btn" id="convert-download-original">Download Original</button>
          <button class="convert-btn primary" id="convert-play" style="display:none">Play Converted</button>
        </div>
      </div>
    </div>
  </div>
</div>

<script>
(function() {
  'use strict';

  // ─── State ──────────────────────────────────────────
  let sessionToken = localStorage.getItem('vidarr_session') || '';
  let currentFiles = [];
  let currentFileIndex = -1;
  let controlsTimer = null;
  let activeDir = null;
  let conversionPollTimer = null;

  // ─── Elements ───────────────────────────────────────
  const $ = id => document.getElementById(id);
  const loginScreen = $('login-screen');
  const app = $('app');
  const tokenInput = $('token-input');
  const loginBtn = $('login-btn');
  const loginError = $('login-error');
  const logoutBtn = $('logout-btn');
  const treeContainer = $('tree-container');
  const fileList = $('file-list');
  const playerWrapper = $('player-wrapper');
  const video = $('video-player');
  const placeholder = $('player-placeholder');
  const controls = $('player-controls');
  const progressBar = $('progress-bar');
  const progressFill = $('progress-fill');
  const btnPlay = $('btn-play');
  const btnPrev = $('btn-prev');
  const btnNext = $('btn-next');
  const btnMute = $('btn-mute');
  const btnFullscreen = $('btn-fullscreen');
  const volumeSlider = $('volume-slider');
  const timeDisplay = $('time-display');
  const touchHint = $('touch-hint');
  const sidebarToggle = $('sidebar-toggle');
  const sidebar = $('sidebar');
  const sidebarBackdrop = $('sidebar-backdrop');
  const convertOverlay = $('convert-overlay');
  const convertProgressFill = $('convert-progress-fill');
  const convertText = $('convert-text');
  const convertTitle = $('convert-title');
  const convertActions = $('convert-actions');
  const convertDownloadOriginal = $('convert-download-original');
  const convertPlay = $('convert-play');

  // ─── API Helpers ────────────────────────────────────
  function api(path, opts = {}) {
    const headers = opts.headers || {};
    if (sessionToken) headers['Authorization'] = 'Bearer ' + sessionToken;
    return fetch(path, { ...opts, headers }).then(r => {
      if (r.status === 401) { doLogout(); throw new Error('Unauthorized'); }
      return r;
    });
  }

  function apiJSON(path, opts) {
    return api(path, opts).then(r => {
      if (!r.ok) throw new Error(r.statusText);
      return r.json();
    });
  }

  // ─── Auth ───────────────────────────────────────────
  function doLogin() {
    const token = tokenInput.value.trim();
    if (!token) return;
    loginBtn.disabled = true;
    loginError.textContent = '';

    fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ access_token: token })
    })
    .then(r => {
      if (!r.ok) throw new Error('Invalid token');
      return r.json();
    })
    .then(data => {
      sessionToken = data.session_token;
      localStorage.setItem('vidarr_session', sessionToken);
      showApp();
    })
    .catch(e => {
      loginError.textContent = e.message;
    })
    .finally(() => { loginBtn.disabled = false; });
  }

  function doLogout() {
    api('/api/logout', { method: 'POST' }).catch(() => {});
    sessionToken = '';
    localStorage.removeItem('vidarr_session');
    loginScreen.classList.remove('hidden');
    app.classList.remove('active');
    tokenInput.value = '';
    video.pause();
    video.removeAttribute('src');
  }

  function showApp() {
    loginScreen.classList.add('hidden');
    app.classList.add('active');
    loadTree();
  }

  // ─── Tree ──────────────────────────────────────────
  function loadTree() {
    apiJSON('/api/tree').then(roots => {
      treeContainer.innerHTML = '';
      if (!roots || roots.length === 0) {
        treeContainer.innerHTML = '<div style="padding:16px;color:var(--text-dim)">No directories available</div>';
        return;
      }
      roots.forEach(root => {
        treeContainer.appendChild(buildTreeNode(root, 0));
      });
    }).catch(() => {});
  }

  function buildTreeNode(entry, depth) {
    const div = document.createElement('div');

    const item = document.createElement('div');
    item.className = 'tree-item';
    item.style.paddingLeft = (12 + depth * 16) + 'px';

    const arrow = document.createElement('span');
    arrow.className = 'tree-arrow' + (entry.children && entry.children.length > 0 ? '' : ' empty');
    arrow.innerHTML = '&#9654;';

    const icon = document.createElement('span');
    icon.className = 'tree-icon';
    icon.textContent = '\uD83D\uDCC1';

    const label = document.createElement('span');
    label.className = 'tree-label';
    label.textContent = entry.name;
    label.title = entry.path;

    item.appendChild(arrow);
    item.appendChild(icon);
    item.appendChild(label);
    div.appendChild(item);

    let childrenContainer = null;
    if (entry.children && entry.children.length > 0) {
      childrenContainer = document.createElement('div');
      childrenContainer.className = 'tree-children';
      entry.children.forEach(child => {
        childrenContainer.appendChild(buildTreeNode(child, depth + 1));
      });
      div.appendChild(childrenContainer);
    }

    item.addEventListener('click', (e) => {
      e.stopPropagation();
      // Toggle children
      if (childrenContainer) {
        const isOpen = childrenContainer.classList.toggle('open');
        arrow.classList.toggle('open', isOpen);
      }
      // Load files
      loadFiles(entry.path);
      // Highlight
      document.querySelectorAll('.tree-item.active').forEach(el => el.classList.remove('active'));
      item.classList.add('active');
      // Mobile: close sidebar
      closeSidebar();
    });

    return div;
  }

  // ─── File List ─────────────────────────────────────
  function loadFiles(dirPath) {
    activeDir = dirPath;
    apiJSON('/api/files?dir=' + encodeURIComponent(dirPath)).then(files => {
      currentFiles = files || [];
      renderFileList();
    }).catch(() => {});
  }

  function renderFileList() {
    fileList.innerHTML = '';
    if (currentFiles.length === 0) {
      fileList.innerHTML = '<div style="padding:12px 16px;color:var(--text-dim);font-size:13px">No videos in this directory</div>';
      return;
    }

    currentFiles.forEach((file, index) => {
      const item = document.createElement('div');
      item.className = 'file-item' + (index === currentFileIndex ? ' active' : '');

      const thumb = document.createElement('img');
      thumb.className = 'file-thumb';
      thumb.loading = 'lazy';
      thumb.src = '/api/thumbnail?path=' + encodeURIComponent(file.path) + '&token=' + sessionToken;
      thumb.onerror = function() { this.style.display = 'none'; };

      const name = document.createElement('span');
      name.className = 'file-name';
      name.textContent = file.name;
      name.title = file.name;

      item.appendChild(thumb);
      item.appendChild(name);

      if (file.needs_remux) {
        const badge = document.createElement('span');
        badge.className = 'file-badge remux';
        badge.textContent = 'remux';
        item.appendChild(badge);
      } else if (file.needs_convert) {
        const badge = document.createElement('span');
        badge.className = 'file-badge convert';
        badge.textContent = 'convert';
        item.appendChild(badge);
      }

      const size = document.createElement('span');
      size.className = 'file-size';
      size.textContent = formatSize(file.size);
      item.appendChild(size);

      item.addEventListener('click', () => playFile(index));
      fileList.appendChild(item);
    });
  }

  // ─── Video Playback ────────────────────────────────
  function playFile(index) {
    if (index < 0 || index >= currentFiles.length) return;
    currentFileIndex = index;
    const file = currentFiles[index];

    // Stop any conversion poll
    if (conversionPollTimer) { clearInterval(conversionPollTimer); conversionPollTimer = null; }
    convertOverlay.classList.remove('visible');

    if (file.needs_remux || file.needs_convert) {
      startConversion(file);
      return;
    }

    loadVideo(file.path);
    renderFileList();
  }

  function loadVideo(path) {
    placeholder.style.display = 'none';
    video.style.display = 'block';
    controls.classList.add('visible');

    const src = '/api/video?path=' + encodeURIComponent(path) + '&token=' + sessionToken;
    video.src = src;

    // Restore position
    const savedPos = localStorage.getItem('vidarr_pos_' + path);
    if (savedPos) {
      video.addEventListener('loadedmetadata', function onMeta() {
        video.removeEventListener('loadedmetadata', onMeta);
        const pos = parseFloat(savedPos);
        if (pos > 0 && pos < video.duration - 5) {
          video.currentTime = pos;
        }
      });
    }

    video.play().catch(() => {});
    showControls();
  }

  // ─── Conversion ────────────────────────────────────
  function startConversion(file) {
    const isRemux = file.needs_remux;
    convertTitle.textContent = isRemux ? 'Remuxing to MP4...' : 'Transcoding to MP4...';
    convertProgressFill.style.width = '0%';
    convertText.textContent = '0%';
    convertPlay.style.display = 'none';
    convertOverlay.classList.add('visible');

    // Download original button
    convertDownloadOriginal.onclick = () => {
      window.open('/api/download?path=' + encodeURIComponent(file.path) + '&token=' + sessionToken, '_blank');
    };

    api('/api/convert', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: file.path })
    })
    .then(r => r.json())
    .then(data => {
      pollConversion(data.job_id, file);
    })
    .catch(e => {
      convertText.textContent = 'Error: ' + e.message;
    });
  }

  function pollConversion(jobId, file) {
    if (conversionPollTimer) clearInterval(conversionPollTimer);
    conversionPollTimer = setInterval(() => {
      apiJSON('/api/convert/status?id=' + jobId).then(status => {
        const pct = Math.round(status.progress);
        convertProgressFill.style.width = pct + '%';
        convertText.textContent = pct + '%';

        if (status.done) {
          clearInterval(conversionPollTimer);
          conversionPollTimer = null;

          if (status.error) {
            convertText.textContent = 'Error: ' + status.error;
          } else {
            convertText.textContent = 'Ready!';
            convertPlay.style.display = 'inline-block';
            convertPlay.onclick = () => {
              convertOverlay.classList.remove('visible');
              const convertedUrl = '/api/convert/download?id=' + jobId + '&token=' + sessionToken;
              placeholder.style.display = 'none';
              video.style.display = 'block';
              controls.classList.add('visible');
              video.src = convertedUrl;
              video.play().catch(() => {});
              showControls();
              renderFileList();
            };
          }
        }
      }).catch(() => {});
    }, 500);
  }

  // ─── Player Controls ───────────────────────────────
  function showControls() {
    controls.classList.add('visible');
    clearTimeout(controlsTimer);
    controlsTimer = setTimeout(() => {
      if (!video.paused) controls.classList.remove('visible');
    }, 3000);
  }

  playerWrapper.addEventListener('mousemove', showControls);
  playerWrapper.addEventListener('click', (e) => {
    if (e.target === video) togglePlay();
  });

  video.addEventListener('timeupdate', () => {
    if (video.duration) {
      const pct = (video.currentTime / video.duration) * 100;
      progressFill.style.width = pct + '%';
      timeDisplay.textContent = formatTime(video.currentTime) + ' / ' + formatTime(video.duration);

      // Save position
      if (currentFileIndex >= 0 && currentFiles[currentFileIndex]) {
        localStorage.setItem('vidarr_pos_' + currentFiles[currentFileIndex].path, video.currentTime.toString());
      }
    }
  });

  video.addEventListener('play', () => { btnPlay.innerHTML = '&#9646;&#9646;'; });
  video.addEventListener('pause', () => { btnPlay.innerHTML = '&#9654;'; showControls(); });

  video.addEventListener('ended', () => {
    // Clear saved position
    if (currentFileIndex >= 0 && currentFiles[currentFileIndex]) {
      localStorage.removeItem('vidarr_pos_' + currentFiles[currentFileIndex].path);
    }
    // Auto-play next
    if (currentFileIndex < currentFiles.length - 1) {
      playFile(currentFileIndex + 1);
    }
  });

  video.addEventListener('volumechange', () => {
    volumeSlider.value = video.muted ? 0 : video.volume;
    btnMute.innerHTML = video.muted || video.volume === 0 ? '&#128263;' : '&#128264;';
  });

  btnPlay.addEventListener('click', togglePlay);
  btnPrev.addEventListener('click', () => playFile(currentFileIndex - 1));
  btnNext.addEventListener('click', () => playFile(currentFileIndex + 1));
  btnMute.addEventListener('click', () => { video.muted = !video.muted; });
  volumeSlider.addEventListener('input', () => { video.volume = volumeSlider.value; video.muted = false; });

  btnFullscreen.addEventListener('click', () => {
    if (document.fullscreenElement) {
      document.exitFullscreen();
    } else {
      playerWrapper.requestFullscreen().catch(() => {});
    }
  });

  progressBar.addEventListener('click', (e) => {
    if (!video.duration) return;
    const rect = progressBar.getBoundingClientRect();
    const pct = (e.clientX - rect.left) / rect.width;
    video.currentTime = pct * video.duration;
  });

  function togglePlay() {
    if (video.paused) video.play().catch(() => {});
    else video.pause();
  }

  // ─── Touch Gestures ────────────────────────────────
  let touchStartX = 0, touchStartY = 0, touchStartTime = 0;

  playerWrapper.addEventListener('touchstart', (e) => {
    if (e.target.closest('.player-controls') || e.target.closest('.convert-overlay')) return;
    touchStartX = e.touches[0].clientX;
    touchStartY = e.touches[0].clientY;
    touchStartTime = Date.now();
  }, { passive: true });

  playerWrapper.addEventListener('touchend', (e) => {
    if (e.target.closest('.player-controls') || e.target.closest('.convert-overlay')) return;
    const dx = e.changedTouches[0].clientX - touchStartX;
    const dy = e.changedTouches[0].clientY - touchStartY;
    const dt = Date.now() - touchStartTime;
    const absDx = Math.abs(dx);
    const absDy = Math.abs(dy);

    if (dt > 500 || absDx < 40) {
      // Tap: toggle play
      if (absDx < 10 && absDy < 10 && dt < 300) {
        togglePlay();
        showControls();
      }
      return;
    }

    const rect = playerWrapper.getBoundingClientRect();
    const touchY = touchStartY - rect.top;
    const isUpperHalf = touchY < rect.height / 2;

    if (isUpperHalf) {
      // Upper half swipe: next/prev video
      if (dx > 0) {
        showHint('\u23EE Prev');
        playFile(currentFileIndex - 1);
      } else {
        showHint('Next \u23ED');
        playFile(currentFileIndex + 1);
      }
    } else {
      // Lower half swipe: seek
      if (dx > 0) {
        video.currentTime = Math.min(video.duration, video.currentTime + 10);
        showHint('+10s \u23E9');
      } else {
        video.currentTime = Math.max(0, video.currentTime - 10);
        showHint('\u23EA -10s');
      }
    }
  }, { passive: true });

  function showHint(text) {
    touchHint.textContent = text;
    touchHint.classList.add('show');
    setTimeout(() => touchHint.classList.remove('show'), 600);
  }

  // ─── Keyboard Shortcuts ────────────────────────────
  document.addEventListener('keydown', (e) => {
    // Don't intercept if typing in input
    if (e.target.tagName === 'INPUT') return;
    if (!video.src) return;

    switch (e.key) {
      case ' ':
        e.preventDefault();
        togglePlay();
        break;
      case 'ArrowRight':
        e.preventDefault();
        video.currentTime = Math.min(video.duration || 0, video.currentTime + 10);
        showControls();
        break;
      case 'ArrowLeft':
        e.preventDefault();
        video.currentTime = Math.max(0, video.currentTime - 10);
        showControls();
        break;
      case 'ArrowUp':
        e.preventDefault();
        video.volume = Math.min(1, video.volume + 0.1);
        break;
      case 'ArrowDown':
        e.preventDefault();
        video.volume = Math.max(0, video.volume - 0.1);
        break;
      case 'f':
        if (document.fullscreenElement) document.exitFullscreen();
        else playerWrapper.requestFullscreen().catch(() => {});
        break;
      case 'm':
        video.muted = !video.muted;
        break;
      case 'Escape':
        if (document.fullscreenElement) document.exitFullscreen();
        break;
      case 'n':
        playFile(currentFileIndex + 1);
        break;
      case 'p':
        playFile(currentFileIndex - 1);
        break;
    }
  });

  // ─── Sidebar Toggle (Mobile) ───────────────────────
  function closeSidebar() {
    sidebar.classList.remove('open');
    sidebarBackdrop.classList.remove('open');
  }

  sidebarToggle.addEventListener('click', () => {
    sidebar.classList.toggle('open');
    sidebarBackdrop.classList.toggle('open');
  });

  sidebarBackdrop.addEventListener('click', closeSidebar);

  // ─── Logout ────────────────────────────────────────
  logoutBtn.addEventListener('click', doLogout);

  // ─── Login Events ──────────────────────────────────
  loginBtn.addEventListener('click', doLogin);
  tokenInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') doLogin(); });

  // ─── Helpers ───────────────────────────────────────
  function formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return '0:00';
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);
    if (h > 0) return h + ':' + String(m).padStart(2, '0') + ':' + String(s).padStart(2, '0');
    return m + ':' + String(s).padStart(2, '0');
  }

  function formatSize(bytes) {
    if (!bytes) return '';
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let i = 0;
    let size = bytes;
    while (size >= 1024 && i < units.length - 1) { size /= 1024; i++; }
    return size.toFixed(i === 0 ? 0 : 1) + ' ' + units[i];
  }

  // ─── Init ──────────────────────────────────────────
  if (sessionToken) {
    // Verify existing session
    api('/api/tree').then(r => {
      if (r.ok) showApp();
      else doLogout();
    }).catch(() => doLogout());
  }
})();
</script>
</body>
</html>
`
