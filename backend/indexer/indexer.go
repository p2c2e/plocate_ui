package indexer

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"plocate-ui/config"

	"github.com/robfig/cron/v3"
)

type Status struct {
	IsIndexing    bool      `json:"is_indexing"`
	LastIndexed   time.Time `json:"last_indexed"`
	LastError     string    `json:"last_error,omitempty"`
	NextScheduled time.Time `json:"next_scheduled"`
	IndexedPaths  []string  `json:"indexed_paths"`
}

type Indexer struct {
	mu            sync.RWMutex
	status        Status
	cron          *cron.Cron
	cancelCurrent context.CancelFunc
}

var Instance *Indexer

func Initialize() error {
	Instance = &Indexer{
		status: Status{
			IndexedPaths: config.AppConfig.Plocate.IndexPaths,
		},
		cron: cron.New(),
	}

	// Setup scheduled indexing
	if config.AppConfig.Scheduler.Enabled {
		_, err := Instance.cron.AddFunc(config.AppConfig.Scheduler.Interval, func() {
			_ = Instance.StartIndexing()
		})
		if err != nil {
			return fmt.Errorf("failed to setup cron: %w", err)
		}

		Instance.cron.Start()
		Instance.updateNextScheduled()
	}

	return nil
}

func (idx *Indexer) GetStatus() Status {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.status
}

func (idx *Indexer) StartIndexing() error {
	idx.mu.Lock()
	if idx.status.IsIndexing {
		idx.mu.Unlock()
		return fmt.Errorf("indexing already in progress")
	}
	idx.status.IsIndexing = true
	idx.status.LastError = ""
	idx.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	idx.mu.Lock()
	idx.cancelCurrent = cancel
	idx.mu.Unlock()

	go func() {
		err := idx.runUpdatedb(ctx)

		idx.mu.Lock()
		idx.status.IsIndexing = false
		if err != nil {
			idx.status.LastError = err.Error()
		} else {
			idx.status.LastIndexed = time.Now()
		}
		idx.cancelCurrent = nil
		idx.mu.Unlock()

		idx.updateNextScheduled()
	}()

	return nil
}

func (idx *Indexer) StopIndexing() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if !idx.status.IsIndexing {
		return fmt.Errorf("no indexing in progress")
	}

	if idx.cancelCurrent != nil {
		idx.cancelCurrent()
	}

	return nil
}

func (idx *Indexer) EnableScheduler() {
	idx.cron.Start()
	idx.updateNextScheduled()
}

func (idx *Indexer) DisableScheduler() {
	idx.cron.Stop()
	idx.mu.Lock()
	idx.status.NextScheduled = time.Time{}
	idx.mu.Unlock()
}

func (idx *Indexer) updateNextScheduled() {
	if len(idx.cron.Entries()) > 0 {
		idx.mu.Lock()
		idx.status.NextScheduled = idx.cron.Entries()[0].Next
		idx.mu.Unlock()
	}
}

func (idx *Indexer) runUpdatedb(ctx context.Context) error {
	cfg := config.AppConfig.Plocate

	// Build updatedb command
	args := []string{
		"--output", cfg.DatabasePath,
		"--prunepaths", "",
	}

	// Add paths to index
	for _, path := range cfg.IndexPaths {
		args = append(args, "--database-root", path)
	}

	cmd := exec.CommandContext(ctx, cfg.UpdatedbBin, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("updatedb failed: %w - %s", err, string(output))
	}

	return nil
}

func (idx *Indexer) Search(query string, limit int) ([]string, error) {
	cfg := config.AppConfig.Plocate

	args := []string{
		"--database", cfg.DatabasePath,
		"--limit", fmt.Sprintf("%d", limit),
	}

	// Add ignore case flag
	args = append(args, "--ignore-case", query)

	cmd := exec.Command(cfg.PlocateBin, args...)

	output, err := cmd.Output()
	if err != nil {
		// plocate returns exit code 1 when no results found
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("plocate search failed: %w", err)
	}

	// Parse results
	var results []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			results = append(results, line)
		}
	}

	return results, nil
}
