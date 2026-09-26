package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Task represents a scheduled task
type Task struct {
	ID       string
	Name     string
	Schedule string
	Handler  func() error
	Enabled  bool
}

// Scheduler manages scheduled tasks
type Scheduler struct {
	cron     *cron.Cron
	tasks    map[string]*Task
	entries  map[string]cron.EntryID
	handlers map[string]func() error
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a new scheduler
func New() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Scheduler{
		cron:     cron.New(),
		tasks:    make(map[string]*Task),
		entries:  make(map[string]cron.EntryID),
		handlers: make(map[string]func() error),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Register adds a task to the scheduler
func (s *Scheduler) Register(task *Task) error {
	s.tasks[task.ID] = task
	s.handlers[task.ID] = task.Handler
	
	if task.Enabled {
		return s.Enable(task.ID)
	}
	
	return nil
}

// Enable activates a scheduled task
func (s *Scheduler) Enable(taskID string) error {
	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}
	
	// Remove existing entry if present
	if entryID, exists := s.entries[taskID]; exists {
		s.cron.Remove(entryID)
	}
	
	entryID, err := s.cron.AddFunc(task.Schedule, func() {
		s.executeTask(taskID)
	})
	
	if err != nil {
		return err
	}
	
	s.entries[taskID] = entryID
	task.Enabled = true
	
	log.Printf("[Scheduler] Enabled task: %s (%s)", task.Name, task.Schedule)
	return nil
}

// Disable deactivates a scheduled task
func (s *Scheduler) Disable(taskID string) error {
	if entryID, exists := s.entries[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, taskID)
	}
	
	if task, exists := s.tasks[taskID]; exists {
		task.Enabled = false
	}
	
	log.Printf("[Scheduler] Disabled task: %s", taskID)
	return nil
}

// executeTask runs a task with error handling
func (s *Scheduler) executeTask(taskID string) {
	task, exists := s.tasks[taskID]
	if !exists {
		return
	}
	
	log.Printf("[Scheduler] Executing task: %s", task.Name)
	
	if err := task.Handler(); err != nil {
		log.Printf("[Scheduler] Task %s failed: %v", task.Name, err)
	} else {
		log.Printf("[Scheduler] Task %s completed", task.Name)
	}
}

// Start begins the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("[Scheduler] Started")
}

// Stop halts the scheduler
func (s *Scheduler) Stop() {
	s.cancel()
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("[Scheduler] Stopped")
}

// List returns all registered tasks
func (s *Scheduler) List() []*Task {
	var list []*Task
	for _, task := range s.tasks {
		list = append(list, task)
	}
	return list
}

// Get returns a specific task
func (s *Scheduler) Get(taskID string) (*Task, error) {
	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	return task, nil
}

// UpdateSchedule changes a task's schedule
func (s *Scheduler) UpdateSchedule(taskID, newSchedule string) error {
	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}
	
	// Validate schedule
	if _, err := cron.ParseStandard(newSchedule); err != nil {
		return fmt.Errorf("invalid cron schedule: %w", err)
	}
	
	task.Schedule = newSchedule
	
	// Re-enable if currently enabled
	if task.Enabled {
		return s.Enable(taskID)
	}
	
	return nil
}

// RunNow executes a task immediately (not scheduled)
func (s *Scheduler) RunNow(taskID string) error {
	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}
	
	go s.executeTask(taskID)
	
	log.Printf("[Scheduler] Manually triggered task: %s", task.Name)
	return nil
}

// GetNextRun returns the next scheduled run time for a task
func (s *Scheduler) GetNextRun(taskID string) (*time.Time, error) {
	entryID, exists := s.entries[taskID]
	if !exists {
		return nil, fmt.Errorf("task not scheduled: %s", taskID)
	}
	
	entry := s.cron.Entry(entryID)
	next := entry.Next
	return &next, nil
}

// DefaultTaskIDs for the content pipeline
const (
	TaskIntelCollector  = "intel_collector"
	TaskScriptWriter    = "script_writer"
	TaskXPostGenerator  = "x_post_generator"
	TaskThumbnailForge  = "thumbnail_forge"
	TaskBlogWriter      = "blog_writer"
	TaskOutreachEngine  = "outreach_engine"
	TaskNewsletter      = "newsletter"
	TaskClipFactory     = "clip_factory"
	TaskPerformanceEval = "performance_eval"
	TaskWeeklyDigest    = "weekly_digest"
)

// GetDefaultTasks returns the standard content pipeline tasks
func GetDefaultTasks() []*Task {
	return []*Task{
		{
			ID:       TaskIntelCollector,
			Name:     "Intel Collector",
			Schedule: "0 5 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskScriptWriter,
			Name:     "Script Writer",
			Schedule: "0 6 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskXPostGenerator,
			Name:     "X Post Generator",
			Schedule: "0 6 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskThumbnailForge,
			Name:     "Thumbnail Forge",
			Schedule: "0 7 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskBlogWriter,
			Name:     "Blog Writer",
			Schedule: "0 7 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskOutreachEngine,
			Name:     "Outreach Engine",
			Schedule: "0 7 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskNewsletter,
			Name:     "Newsletter Composer",
			Schedule: "0 8 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskClipFactory,
			Name:     "Clip Factory",
			Schedule: "0 */4 * * *",
			Enabled:  true,
		},
		{
			ID:       TaskPerformanceEval,
			Name:     "Performance Evaluator",
			Schedule: "0 2 * * 0",
			Enabled:  true,
		},
		{
			ID:       TaskWeeklyDigest,
			Name:     "Weekly Digest",
			Schedule: "0 9 * * 1",
			Enabled:  true,
		},
	}
}
