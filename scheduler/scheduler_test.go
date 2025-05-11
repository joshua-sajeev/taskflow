package scheduler_test

import (
	"testing"
	"time"

	"taskflow/scheduler"
)

func TestSchedule(t *testing.T) {
	tests := []struct {
		name         string
		cronExpr     string
		jobTriggered chan bool
		expectError  bool
	}{
		{
			name:         "Valid: Job runs once within 2 seconds",
			cronExpr:     "@every 1s",
			jobTriggered: make(chan bool, 1),
			expectError:  false,
		},
		{
			name:         "Valid: Job runs twice within 3 seconds",
			cronExpr:     "@every 1s",
			jobTriggered: make(chan bool, 2),
			expectError:  false,
		},
		{
			name:         "Invalid: Incorrect cron expression",
			cronExpr:     "@every xyz",
			jobTriggered: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var runCount int
			var err error

			if !tt.expectError {
				err = scheduler.Schedule(tt.cronExpr, func() {
					runCount++
					if tt.jobTriggered != nil {
						tt.jobTriggered <- true
					}
				})
			} else {
				err = scheduler.Schedule(tt.cronExpr, func() {})
			}

			if tt.expectError {
				if err == nil {
					t.Error("Expected error for invalid cron expression, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to schedule job: %v", err)
			}

			timeout := time.After(3 * time.Second)
			expectedRuns := cap(tt.jobTriggered)

			for range expectedRuns {
				select {
				case <-tt.jobTriggered:
					// Job ran
				case <-timeout:
					t.Errorf("Job did not run expected %d times in time", expectedRuns)
					return
				}
			}
		})
	}
}
