package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"dramastudio/internal/postproduction/domain"
)

// FFmpegAdapter builds deterministic, seek-safe ffmpeg invocations from a
// structured timeline (§42). Arguments are always passed as an argv slice —
// never a shell string — so paths/URLs cannot inject flags.
type FFmpegAdapter struct {
	binaryPath string
}

func NewFFmpegAdapter(path string) *FFmpegAdapter {
	if path == "" {
		path = "ffmpeg"
	}
	return &FFmpegAdapter{binaryPath: path}
}

// ConcatOptions describes a deterministic video concat + audio mix render.
type ConcatOptions struct {
	VideoURLs   []string
	AudioTracks []domain.AudioTrack // voice/music/sfx mixed via amix
	SubtitleURL string
	Format      string // mp4, webm
	Resolution  string // e.g. 1080x1920
	OutputPath  string
}

// BuildConcatCommand returns (binary, argv) for a full render with a
// single -filter_complex graph: per-clip normalize → concat → optional
// subtitle burn → optional audio mix.
func (f *FFmpegAdapter) BuildConcatCommand(opts ConcatOptions) (string, []string) {
	resolution := opts.Resolution
	if resolution == "" {
		resolution = "1080x1920"
	}
	dim := strings.ReplaceAll(resolution, "x", ":")
	nV := len(opts.VideoURLs)
	nA := len(opts.AudioTracks)

	args := []string{"-y"}
	for _, u := range opts.VideoURLs {
		args = append(args, "-i", u)
	}
	for _, a := range opts.AudioTracks {
		args = append(args, "-i", a.URL)
	}

	var graph strings.Builder
	for i := 0; i < nV; i++ {
		fmt.Fprintf(&graph, "[%d:v]scale=%s:force_original_aspect_ratio=decrease,pad=%s:(ow-iw)/2:(oh-ih)/2[v%d];", i, dim, dim, i)
	}
	for i := 0; i < nV; i++ {
		fmt.Fprintf(&graph, "[v%d]", i)
	}
	if opts.SubtitleURL != "" {
		fmt.Fprintf(&graph, "concat=n=%d:v=1:a=0[vcat];[vcat]subtitles=%s[outv]", nV, escapeFilterPath(opts.SubtitleURL))
	} else {
		fmt.Fprintf(&graph, "concat=n=%d:v=1:a=0[outv]", nV)
	}

	args = append(args, "-filter_complex", graph.String())
	args = append(args, "-map", "[outv]")

	if nA > 0 {
		var mix strings.Builder
		for i, t := range opts.AudioTracks {
			vol := t.Volume
			if vol == 0 {
				vol = 1.0
			}
			fmt.Fprintf(&mix, ";[%d:a]volume=%.2f[a%d]", nV+i, vol, i)
		}
		for i := 0; i < nA; i++ {
			fmt.Fprintf(&mix, "[a%d]", i)
		}
		fmt.Fprintf(&mix, "amix=inputs=%d:duration=first[outa]", nA)
		// Rebuild args with the combined graph (single -filter_complex).
		args = []string{"-y"}
		for _, u := range opts.VideoURLs {
			args = append(args, "-i", u)
		}
		for _, a := range opts.AudioTracks {
			args = append(args, "-i", a.URL)
		}
		args = append(args, "-filter_complex", graph.String()+mix.String(),
			"-map", "[outv]", "-map", "[outa]")
	}

	args = append(args, opts.OutputPath)
	return f.binaryPath, args
}

// RenderTimeline executes the render. Errors are returned — callers must
// not fabricate output URLs on failure.
func (f *FFmpegAdapter) RenderTimeline(ctx context.Context, timeline *domain.Timeline, outputPath string) (string, error) {
	if timeline == nil || len(timeline.VideoTracks) == 0 {
		return "", fmt.Errorf("timeline is empty or missing video tracks")
	}
	if _, err := exec.LookPath(f.binaryPath); err != nil {
		return "", fmt.Errorf("ffmpeg binary %q not found: %w", f.binaryPath, err)
	}

	urls := make([]string, len(timeline.VideoTracks))
	for i, t := range timeline.VideoTracks {
		urls[i] = t.AssetURL
	}

	audioTracks := make([]domain.AudioTrack, 0, len(timeline.AudioTracks))
	for _, t := range timeline.AudioTracks {
		audioTracks = append(audioTracks, domain.AudioTrack{Type: "track", URL: t.AssetURL, Volume: 1.0})
	}

	bin, args := f.BuildConcatCommand(ConcatOptions{
		VideoURLs:   urls,
		AudioTracks: audioTracks,
		Resolution:  "1080x1920",
		OutputPath:  outputPath,
	})

	cmd := exec.CommandContext(ctx, bin, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg render failed: %w: %s", err, truncate(string(out), 500))
	}
	return outputPath, nil
}

// escapeFilterPath escapes a filesystem path for use inside a
// -filter_complex option value (drive-letter colons, backslashes,
// single-quotes must not be parsed as filter syntax).
func escapeFilterPath(p string) string {
	r := strings.NewReplacer(`\`, `\\`, `:`, `\:`, `'`, `\'`)
	return r.Replace(p)
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
