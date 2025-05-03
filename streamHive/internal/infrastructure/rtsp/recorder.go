package rtsp

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

type Recorder struct {
	procs sync.Map // outputPath -> *exec.Cmd
}

func New() *Recorder {
	return &Recorder{}
}

func (r *Recorder) Record(ctx context.Context, streamURL, outputPath string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", streamURL, "-t", "600", "-c", "copy", outputPath)
	r.procs.Store(outputPath, cmd)
	err := cmd.Run()
	r.procs.Delete(outputPath)
	return err
}

func (r *Recorder) Stop(jobID string) error {
	path := fmt.Sprintf("/tmp/%s.mp4", jobID)
	if proc, ok := r.procs.Load(path); ok {
		cmd := proc.(*exec.Cmd)
		_ = cmd.Process.Kill()
		r.procs.Delete(path)
	}
	return nil
}
