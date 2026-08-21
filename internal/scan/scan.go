package scan

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/jeonjw85/Ksecret/internal/config"
	"github.com/jeonjw85/Ksecret/internal/detect"
	"golang.org/x/sync/errgroup"
)

type Options struct {
	Root   string
	Staged bool
	Config config.Config
}

var skipDirs = map[string]struct{}{
	".git": {}, "node_modules": {}, "vendor": {}, "dist": {},
	"build": {}, ".idea": {}, ".venv": {}, "__pycache__": {},
}

func Run(ctx context.Context, opts Options) ([]detect.Finding, error) {
	files, err := listFiles(ctx, opts)
	if err != nil {
		return nil, err
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.GOMAXPROCS(0))
	var mu sync.Mutex
	var all []detect.Finding
	for _, path := range files {
		g.Go(func() error {
			found, err := scanFile(ctx, opts, path)
			if err != nil {
				return err
			}
			if len(found) == 0 {
				return nil
			}
			mu.Lock()
			all = append(all, found...)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	slices.SortFunc(all, cmpFinding)
	return all, nil
}

func cmpFinding(a, b detect.Finding) int {
	if a.File != b.File {
		return strings.Compare(a.File, b.File)
	}
	if a.Line != b.Line {
		return a.Line - b.Line
	}
	if a.Column != b.Column {
		return a.Column - b.Column
	}
	return strings.Compare(a.RuleID, b.RuleID)
}

func listFiles(ctx context.Context, opts Options) ([]string, error) {
	if opts.Staged {
		return stagedFiles(ctx, opts.Root)
	}
	var files []string
	err := filepath.WalkDir(opts.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		name := d.Name()
		if d.IsDir() {
			if _, skip := skipDirs[name]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		rel, relErr := filepath.Rel(opts.Root, path)
		if relErr != nil {
			rel = path
		}
		if opts.Config.PathExcluded(rel) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("디렉터리 탐색 %s: %w", opts.Root, err)
	}
	return files, nil
}

func stagedFiles(ctx context.Context, root string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "--cached", "--name-only", "-z", "--diff-filter=ACMRT")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("커밋 대기 파일 목록 : %w", err)
	}
	var files []string
	for _, rel := range strings.Split(string(bytes.TrimRight(out, "\x00")), "\x00") {
		if rel == "" {
			continue
		}
		files = append(files, filepath.Join(root, rel))
	}
	return files, nil
}

func scanFile(ctx context.Context, opts Options, path string) ([]detect.Finding, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("파일 정보 %s : %w", path, err)
	}
	if !info.Mode().IsRegular() || info.Size() > opts.Config.MaxFileBytes {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 %s : %w", path, err)
	}
	defer f.Close()

	head := make([]byte, 8192)
	n, err := io.ReadFull(f, head)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("파일 읽기 %s : %w", path, err)
	}
	head = head[:n]
	if bytes.IndexByte(head, 0) >= 0 {
		return nil, nil
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("파일 위치 이동 %s : %w", path, err)
	}

	rel, relErr := filepath.Rel(opts.Root, path)
	if relErr != nil {
		rel = path
	}
	rel = filepath.ToSlash(rel)

	var findings []detect.Finding
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := sc.Text()
		for _, m := range detect.ScanLine(line) {
			if opts.Config.Allowed(rel, m.RuleID, line) {
				continue
			}
			findings = append(findings, detect.Finding{
				RuleID:   m.RuleID,
				Severity: m.Severity,
				File:     rel,
				Line:     lineNo,
				Column:   m.Column,
				Secret:   m.Secret,
				Message:  m.Message,
			})
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("파일 스캔 %s : %w", path, err)
	}
	return findings, nil
}
