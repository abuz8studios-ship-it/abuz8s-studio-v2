package scheduler

import (
	"testing"
	"time"
)

// Default tasks use 5-field cron specs; they must register without error.
func TestDefaultTasksRegister(t *testing.T) {
	s := New()
	defer s.Stop()

	for _, task := range GetDefaultTasks() {
		task.Handler = func() error { return nil }
		if err := s.Register(task); err != nil {
			t.Fatalf("register %s (%q) failed: %v", task.ID, task.Schedule, err)
		}
	}

	if got := len(s.List()); got != 10 {
		t.Fatalf("expected 10 tasks, got %d", got)
	}

	s.Start()
	defer s.Stop()

	// Every enabled task must have a real next run scheduled.
	for _, task := range s.List() {
		if !task.Enabled {
			t.Fatalf("task %s not enabled after register", task.ID)
		}
		next, err := s.GetNextRun(task.ID)
		if err != nil {
			t.Fatalf("GetNextRun(%s) failed: %v", task.ID, err)
		}
		if next.IsZero() || next.Before(time.Now()) {
			t.Fatalf("task %s has no future next run: %v", task.ID, next)
		}
	}
}

func TestUpdateScheduleRejectsInvalid(t *testing.T) {
	s := New()
	defer s.Stop()

	task := &Task{ID: "t1", Name: "t1", Schedule: "0 5 * * *", Handler: func() error { return nil }, Enabled: false}
	if err := s.Register(task); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := s.UpdateSchedule("t1", "not a cron"); err == nil {
		t.Fatal("expected error for invalid cron spec, got nil")
	}
	if err := s.UpdateSchedule("t1", "0 6 * * *"); err != nil {
		t.Fatalf("valid reschedule failed: %v", err)
	}
}
