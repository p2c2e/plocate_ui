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

type IndexStatus struct {
	Name          string    `json:"name"`
	IsIndexing    bool      `json:"is_indexing"`
	LastIndexed   time.Time `json:"last_indexed"`
	LastError     string    `json:"last_error,omitempty"`
	IndexedPaths  []string  `json:"indexed_paths"`
	Enabled       bool      `json:"enabled"`
	DatabasePath  string    `json:"database_path"`
}

type Status struct {
	Indices       []IndexStatus `json:"indices"`
	NextScheduled time.Time     `json:"next_scheduled"`
}

type Indexer struct {
	mu             sync.RWMutex
	indexStatuses  map[string]*IndexStatus
	cron           *cron.Cron
	cancelFuncs    map[string]context.CancelFunc
	nextScheduled  time.Time
}

var Instance *Indexer

func Initialize() error {
	indexStatuses := make(map[string]*IndexStatus)

	// Initialize status for each configured index
	for _, indexCfg := range config.AppConfig.Plocate.Indices {
		indexStatuses[indexCfg.Name] = &IndexStatus{
			Name:         indexCfg.Name,
			IsIndexing:   false,
			IndexedPaths: indexCfg.IndexPaths,
			Enabled:      indexCfg.Enabled,
			DatabasePath: indexCfg.DatabasePath,
		}
	}

	Instance = &Indexer{
		indexStatuses: indexStatuses,
		cron:          cron.New(),
		cancelFuncs:   make(map[string]context.CancelFunc),
	}

	// Setup scheduled indexing (indexes all enabled indices)
	if config.AppConfig.Scheduler.Enabled {
		_, err := Instance.cron.AddFunc(config.AppConfig.Scheduler.Interval, func() {
			_ = Instance.StartIndexingAll()
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

	indices := make([]IndexStatus, 0, len(idx.indexStatuses))
	for _, status := range idx.indexStatuses {
		indices = append(indices, *status)
	}

	return Status{
		Indices:       indices,
		NextScheduled: idx.nextScheduled,
	}
}

func (idx *Indexer) GetIndexNames() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	names := make([]string, 0, len(idx.indexStatuses))
	for name := range idx.indexStatuses {
		names = append(names, name)
	}
	return names
}

func (idx *Indexer) StartIndexing(indexName string) error {
	idx.mu.Lock()
	status, exists := idx.indexStatuses[indexName]
	if !exists {
		idx.mu.Unlock()
		return fmt.Errorf("index '%s' not found", indexName)
	}

	if status.IsIndexing {
		idx.mu.Unlock()
		return fmt.Errorf("index '%s' is already being indexed", indexName)
	}

	status.IsIndexing = true
	status.LastError = ""
	idx.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	idx.mu.Lock()
	idx.cancelFuncs[indexName] = cancel
	idx.mu.Unlock()

	go func() {
		err := idx.runUpdatedb(ctx, indexName)

		idx.mu.Lock()
		status.IsIndexing = false
		if err != nil {
			status.LastError = err.Error()
		} else {
			status.LastIndexed = time.Now()
		}
		delete(idx.cancelFuncs, indexName)
		idx.mu.Unlock()

		idx.updateNextScheduled()
	}()

	return nil
}

func (idx *Indexer) StartIndexingAll() error {
	var errors []string

	for _, indexCfg := range config.AppConfig.Plocate.Indices {
		if indexCfg.Enabled {
			if err := idx.StartIndexing(indexCfg.Name); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", indexCfg.Name, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start some indices: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (idx *Indexer) StopIndexing(indexName string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	status, exists := idx.indexStatuses[indexName]
	if !exists {
		return fmt.Errorf("index '%s' not found", indexName)
	}

	if !status.IsIndexing {
		return fmt.Errorf("index '%s' is not being indexed", indexName)
	}

	if cancel, exists := idx.cancelFuncs[indexName]; exists {
		cancel()
	}

	return nil
}

func (idx *Indexer) StopIndexingAll() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for name, cancel := range idx.cancelFuncs {
		if cancel != nil {
			cancel()
		}
		if status, exists := idx.indexStatuses[name]; exists {
			status.IsIndexing = false
		}
	}

	idx.cancelFuncs = make(map[string]context.CancelFunc)
	return nil
}

func (idx *Indexer) EnableScheduler() {
	idx.cron.Start()
	idx.updateNextScheduled()
}

func (idx *Indexer) DisableScheduler() {
	idx.cron.Stop()
	idx.mu.Lock()
	idx.nextScheduled = time.Time{}
	idx.mu.Unlock()
}

// AddIndex registers a new index at runtime.
func (idx *Indexer) AddIndex(cfg config.IndexConfig) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.indexStatuses[cfg.Name] = &IndexStatus{
		Name:         cfg.Name,
		IsIndexing:   false,
		IndexedPaths: cfg.IndexPaths,
		Enabled:      cfg.Enabled,
		DatabasePath: cfg.DatabasePath,
	}
}

// RemoveIndex stops and deregisters an index at runtime.
func (idx *Indexer) RemoveIndex(name string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	status, exists := idx.indexStatuses[name]
	if !exists {
		return fmt.Errorf("index '%s' not found", name)
	}

	// Stop if currently indexing
	if status.IsIndexing {
		if cancel, ok := idx.cancelFuncs[name]; ok {
			cancel()
			delete(idx.cancelFuncs, name)
		}
	}

	delete(idx.indexStatuses, name)
	return nil
}

func (idx *Indexer) updateNextScheduled() {
	if len(idx.cron.Entries()) > 0 {
		idx.mu.Lock()
		idx.nextScheduled = idx.cron.Entries()[0].Next
		idx.mu.Unlock()
	}
}

func (idx *Indexer) runUpdatedb(ctx context.Context, indexName string) error {
	cfg := config.AppConfig.Plocate

	// Find the index configuration
	var indexCfg *config.IndexConfig
	for i := range cfg.Indices {
		if cfg.Indices[i].Name == indexName {
			indexCfg = &cfg.Indices[i]
			break
		}
	}

	if indexCfg == nil {
		return fmt.Errorf("index configuration for '%s' not found", indexName)
	}

	// Build updatedb command
	args := []string{
		"--output", indexCfg.DatabasePath,
		"--prunepaths", "",
	}

	// Add paths to index
	for _, path := range indexCfg.IndexPaths {
		args = append(args, "--database-root", path)
	}

	cmd := exec.CommandContext(ctx, cfg.UpdatedbBin, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("updatedb failed: %w - %s", err, string(output))
	}

	return nil
}

func (idx *Indexer) Search(query string, limit int, indexNames []string) ([]string, error) {
	cfg := config.AppConfig.Plocate

	// If no indices specified, search all enabled indices
	if len(indexNames) == 0 {
		for _, indexCfg := range cfg.Indices {
			if indexCfg.Enabled {
				indexNames = append(indexNames, indexCfg.Name)
			}
		}
	}

	if len(indexNames) == 0 {
		return []string{}, fmt.Errorf("no indices available to search")
	}

	// Collect database paths for the specified indices
	var dbPaths []string
	for _, indexName := range indexNames {
		found := false
		for _, indexCfg := range cfg.Indices {
			if indexCfg.Name == indexName {
				dbPaths = append(dbPaths, indexCfg.DatabasePath)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("index '%s' not found", indexName)
		}
	}

	// Build plocate command with multiple databases
	args := []string{}

	// Add all database paths (colon-separated as plocate expects)
	args = append(args, "--database", strings.Join(dbPaths, ":"))

	args = append(args, "--limit", fmt.Sprintf("%d", limit))
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
