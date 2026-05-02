package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the vidarr.yaml configuration.
type Config struct {
	Sources []Source `yaml:"sources"`
}

// Source represents a named source with an access token and directories.
type Source struct {
	Name        string   `yaml:"name"`
	AccessToken string   `yaml:"access_token"`
	Directories []string `yaml:"directories"`
}

// Session represents an authenticated user session.
type Session struct {
	Token     string
	Source    *Source
	ExpiresAt time.Time
}

// ConversionJob tracks a running ffmpeg conversion.
type ConversionJob struct {
	ID         string  `json:"id"`
	SourcePath string  `json:"source_path"`
	OutputPath string  `json:"output_path"`
	Progress   float64 `json:"progress"`
	Done       bool    `json:"done"`
	Error      string  `json:"error,omitempty"`
	Duration   float64 `json:"-"`
	mu         sync.Mutex
}

// DirEntry represents a file or directory in the tree.
type DirEntry struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"is_dir"`
	Children []DirEntry `json:"children,omitempty"`
}

var (
	config       Config
	sessions     = make(map[string]*Session)
	sessionsMu   sync.RWMutex
	conversions  = make(map[string]*ConversionJob)
	conversionMu sync.RWMutex
	thumbnailDir string
	convertDir   string
)

var videoExtensions = map[string]bool{
	".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
	".webm": true, ".m4v": true, ".wmv": true, ".flv": true,
	".ts": true, ".mts": true, ".m2ts": true, ".ogv": true,
}

// nativeFormats are formats that browsers can play directly.
var nativeFormats = map[string]bool{
	".mp4": true, ".webm": true, ".m4v": true, ".ogv": true,
}

// remuxFormats can be quickly remuxed to MP4 (same codec, different container).
var remuxFormats = map[string]bool{
	".mkv": true, ".ts": true, ".mts": true, ".m2ts": true,
}

func main() {
	configPath := "vidarr.yaml"
	port := "8080"

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--config":
			if i+1 < len(os.Args) {
				configPath = os.Args[i+1]
				i++
			}
		case "--port":
			if i+1 < len(os.Args) {
				port = os.Args[i+1]
				i++
			}
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config %s: %v", configPath, err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	if len(config.Sources) == 0 {
		log.Fatal("No sources defined in config")
	}

	// Create data directories.
	thumbnailDir = filepath.Join(os.TempDir(), "vidarr-thumbnails")
	convertDir = filepath.Join(os.TempDir(), "vidarr-converted")
	os.MkdirAll(thumbnailDir, 0755)
	os.MkdirAll(convertDir, 0755)

	mux := http.NewServeMux()

	// API routes.
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/tree", authMiddleware(handleTree))
	mux.HandleFunc("/api/files", authMiddleware(handleFiles))
	mux.HandleFunc("/api/video", authMiddleware(handleVideo))
	mux.HandleFunc("/api/thumbnail", authMiddleware(handleThumbnail))
	mux.HandleFunc("/api/convert", authMiddleware(handleConvert))
	mux.HandleFunc("/api/convert/status", authMiddleware(handleConvertStatus))
	mux.HandleFunc("/api/convert/download", authMiddleware(handleConvertDownload))
	mux.HandleFunc("/api/download", authMiddleware(handleDownload))
	mux.HandleFunc("/api/logout", handleLogout)

	// Static files.
	mux.HandleFunc("/", handleIndex)

	log.Printf("Vidarr starting on port %s with %d source(s)", port, len(config.Sources))
	for _, s := range config.Sources {
		log.Printf("  Source %q: %d directories", s.Name, len(s.Directories))
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  0, // No timeout for video streaming.
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Find matching source using constant-time comparison.
	var matched *Source
	for i := range config.Sources {
		tokenHash := sha256.Sum256([]byte(config.Sources[i].AccessToken))
		inputHash := sha256.Sum256([]byte(req.AccessToken))
		if subtle.ConstantTimeCompare(tokenHash[:], inputHash[:]) == 1 {
			matched = &config.Sources[i]
			break
		}
	}

	if matched == nil {
		// Small delay to slow down brute-force attempts.
		time.Sleep(500 * time.Millisecond)
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	sessionToken := generateSessionToken()
	sessionsMu.Lock()
	sessions[sessionToken] = &Session{
		Token:     sessionToken,
		Source:    matched,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	sessionsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"session_token": sessionToken,
		"source_name":   matched.Name,
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token != "" {
		sessionsMu.Lock()
		delete(sessions, token)
		sessionsMu.Unlock()
	}
	w.WriteHeader(http.StatusOK)
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return r.URL.Query().Get("token")
}

func getSession(r *http.Request) *Session {
	token := extractToken(r)
	if token == "" {
		return nil
	}
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	s, ok := sessions[token]
	if !ok || time.Now().After(s.ExpiresAt) {
		return nil
	}
	return s
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if getSession(r) == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// isPathAllowed checks that the requested path is within one of the session's directories.
func isPathAllowed(session *Session, requestedPath string) bool {
	cleaned := filepath.Clean(requestedPath)
	for _, dir := range session.Source.Directories {
		cleanDir := filepath.Clean(dir)
		if cleaned == cleanDir || strings.HasPrefix(cleaned, cleanDir+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func handleTree(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	var roots []DirEntry

	for _, dir := range session.Source.Directories {
		entry := buildTree(dir, dir, 0)
		roots = append(roots, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roots)
}

func buildTree(basePath, fullPath string, depth int) DirEntry {
	name := filepath.Base(fullPath)
	if fullPath == basePath {
		name = filepath.Base(basePath)
	}

	entry := DirEntry{
		Name:  name,
		Path:  fullPath,
		IsDir: true,
	}

	if depth > 20 {
		return entry
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return entry
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		childPath := filepath.Join(fullPath, e.Name())
		if e.IsDir() {
			child := buildTree(basePath, childPath, depth+1)
			if hasVideoFiles(child) {
				entry.Children = append(entry.Children, child)
			}
		}
	}

	// Sort children alphabetically.
	sort.Slice(entry.Children, func(i, j int) bool {
		return strings.ToLower(entry.Children[i].Name) < strings.ToLower(entry.Children[j].Name)
	})

	return entry
}

// hasVideoFiles returns true if a directory tree contains any video files.
func hasVideoFiles(entry DirEntry) bool {
	// Check if the directory itself has videos.
	entries, err := os.ReadDir(entry.Path)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && isVideoFile(e.Name()) {
				return true
			}
		}
	}
	// Check children.
	for _, child := range entry.Children {
		if hasVideoFiles(child) {
			return true
		}
	}
	return false
}

func isVideoFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return videoExtensions[ext]
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		http.Error(w, "Missing dir parameter", http.StatusBadRequest)
		return
	}

	if !isPathAllowed(session, dir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	type FileInfo struct {
		Name         string `json:"name"`
		Path         string `json:"path"`
		Size         int64  `json:"size"`
		NativePlay   bool   `json:"native_play"`
		NeedsRemux   bool   `json:"needs_remux"`
		NeedsConvert bool   `json:"needs_convert"`
	}

	var files []FileInfo
	for _, e := range entries {
		if e.IsDir() || !isVideoFile(e.Name()) || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		info, _ := e.Info()
		size := int64(0)
		if info != nil {
			size = info.Size()
		}

		fi := FileInfo{
			Name:       e.Name(),
			Path:       filepath.Join(dir, e.Name()),
			Size:       size,
			NativePlay: nativeFormats[ext],
		}
		if !fi.NativePlay {
			if remuxFormats[ext] {
				fi.NeedsRemux = true
			} else {
				fi.NeedsConvert = true
			}
		}
		files = append(files, fi)
	}

	// Sort by name.
	sort.Slice(files, func(i, j int) bool {
		return naturalLess(files[i].Name, files[j].Name)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// naturalLess compares strings using natural sort order (so "ep2" < "ep10").
func naturalLess(a, b string) bool {
	re := regexp.MustCompile(`(\d+|\D+)`)
	partsA := re.FindAllString(strings.ToLower(a), -1)
	partsB := re.FindAllString(strings.ToLower(b), -1)

	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		na, errA := strconv.Atoi(partsA[i])
		nb, errB := strconv.Atoi(partsB[i])
		if errA == nil && errB == nil {
			if na != nb {
				return na < nb
			}
		} else {
			if partsA[i] != partsB[i] {
				return partsA[i] < partsB[i]
			}
		}
	}
	return len(partsA) < len(partsB)
}

func handleVideo(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	videoPath := r.URL.Query().Get("path")
	if videoPath == "" {
		http.Error(w, "Missing path", http.StatusBadRequest)
		return
	}

	if !isPathAllowed(session, videoPath) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	fi, err := os.Stat(videoPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(videoPath))
	contentType := "video/mp4"
	switch ext {
	case ".webm":
		contentType = "video/webm"
	case ".ogv":
		contentType = "video/ogg"
	case ".mkv":
		contentType = "video/x-matroska"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

func handleThumbnail(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	videoPath := r.URL.Query().Get("path")
	if videoPath == "" {
		http.Error(w, "Missing path", http.StatusBadRequest)
		return
	}

	if !isPathAllowed(session, videoPath) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Generate a deterministic thumbnail path.
	hash := sha256.Sum256([]byte(videoPath))
	thumbFile := filepath.Join(thumbnailDir, hex.EncodeToString(hash[:])+".jpg")

	// Check if already generated.
	if _, err := os.Stat(thumbFile); err == nil {
		http.ServeFile(w, r, thumbFile)
		return
	}

	// Generate thumbnail with ffmpeg.
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:05",
		"-vframes", "1",
		"-vf", "scale=320:-1",
		"-y",
		thumbFile,
	)
	if err := cmd.Run(); err != nil {
		// Try at 1 second if 5 seconds fails (short video).
		cmd2 := exec.Command("ffmpeg",
			"-i", videoPath,
			"-ss", "00:00:01",
			"-vframes", "1",
			"-vf", "scale=320:-1",
			"-y",
			thumbFile,
		)
		if err2 := cmd2.Run(); err2 != nil {
			http.Error(w, "Failed to generate thumbnail", http.StatusInternalServerError)
			return
		}
	}

	http.ServeFile(w, r, thumbFile)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := getSession(r)
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !isPathAllowed(session, req.Path) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	ext := strings.ToLower(filepath.Ext(req.Path))

	// Determine output filename.
	hash := sha256.Sum256([]byte(req.Path))
	outputName := hex.EncodeToString(hash[:]) + ".mp4"
	outputPath := filepath.Join(convertDir, outputName)

	// Check if already converted.
	if _, err := os.Stat(outputPath); err == nil {
		conversionMu.Lock()
		jobID := hex.EncodeToString(hash[:])[:16]
		conversions[jobID] = &ConversionJob{
			ID:         jobID,
			SourcePath: req.Path,
			OutputPath: outputPath,
			Progress:   100,
			Done:       true,
		}
		conversionMu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
		return
	}

	// Check if already in progress.
	jobID := hex.EncodeToString(hash[:])[:16]
	conversionMu.RLock()
	existing, exists := conversions[jobID]
	conversionMu.RUnlock()
	if exists && !existing.Done {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
		return
	}

	// Start conversion.
	job := &ConversionJob{
		ID:         jobID,
		SourcePath: req.Path,
		OutputPath: outputPath,
	}
	conversionMu.Lock()
	conversions[jobID] = job
	conversionMu.Unlock()

	isRemux := remuxFormats[ext]

	go func() {
		// Get duration first.
		durCmd := exec.Command("ffprobe",
			"-v", "quiet",
			"-show_entries", "format=duration",
			"-of", "default=noprint_wrappers=1:nokey=1",
			req.Path,
		)
		durOut, err := durCmd.Output()
		if err == nil {
			dur, _ := strconv.ParseFloat(strings.TrimSpace(string(durOut)), 64)
			job.mu.Lock()
			job.Duration = dur
			job.mu.Unlock()
		}

		var cmd *exec.Cmd
		if isRemux {
			cmd = exec.Command("ffmpeg",
				"-i", req.Path,
				"-c", "copy",
				"-movflags", "+faststart",
				"-y",
				"-progress", "pipe:1",
				outputPath,
			)
		} else {
			cmd = exec.Command("ffmpeg",
				"-i", req.Path,
				"-c:v", "libx264",
				"-preset", "fast",
				"-crf", "22",
				"-c:a", "aac",
				"-b:a", "192k",
				"-movflags", "+faststart",
				"-y",
				"-progress", "pipe:1",
				outputPath,
			)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			job.mu.Lock()
			job.Error = err.Error()
			job.Done = true
			job.mu.Unlock()
			return
		}

		if err := cmd.Start(); err != nil {
			job.mu.Lock()
			job.Error = err.Error()
			job.Done = true
			job.mu.Unlock()
			return
		}

		// Parse progress.
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				lines := strings.Split(string(buf[:n]), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "out_time_us=") {
						usStr := strings.TrimPrefix(line, "out_time_us=")
						us, e := strconv.ParseFloat(usStr, 64)
						if e == nil {
							job.mu.Lock()
							if job.Duration > 0 {
								job.Progress = (us / 1e6 / job.Duration) * 100
								if job.Progress > 99 {
									job.Progress = 99
								}
							}
							job.mu.Unlock()
						}
					}
				}
			}
			if err != nil {
				break
			}
		}

		if err := cmd.Wait(); err != nil {
			job.mu.Lock()
			job.Error = err.Error()
			job.Done = true
			job.mu.Unlock()
			os.Remove(outputPath)
			return
		}

		job.mu.Lock()
		job.Progress = 100
		job.Done = true
		job.mu.Unlock()
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func handleConvertStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	conversionMu.RLock()
	job, ok := conversions[jobID]
	conversionMu.RUnlock()

	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	job.mu.Lock()
	resp := map[string]interface{}{
		"id":       job.ID,
		"progress": job.Progress,
		"done":     job.Done,
	}
	if job.Error != "" {
		resp["error"] = job.Error
	}
	job.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleConvertDownload(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	conversionMu.RLock()
	job, ok := conversions[jobID]
	conversionMu.RUnlock()

	if !ok || !job.Done || job.Error != "" {
		http.Error(w, "Conversion not ready", http.StatusBadRequest)
		return
	}

	if !isPathAllowed(session, job.SourcePath) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(job.OutputPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, _ := f.Stat()
	baseName := strings.TrimSuffix(filepath.Base(job.SourcePath), filepath.Ext(job.SourcePath)) + ".mp4"
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", baseName))
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, baseName, fi.ModTime(), f)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	videoPath := r.URL.Query().Get("path")
	if videoPath == "" {
		http.Error(w, "Missing path", http.StatusBadRequest)
		return
	}

	if !isPathAllowed(session, videoPath) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, _ := f.Stat()
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(videoPath)))
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	io.WriteString(w, indexHTML)
}
