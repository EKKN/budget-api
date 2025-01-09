package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type JobIDGenerator struct {
	lastMinute int
	sequence   int
	mutex      sync.Mutex
}

var (
	jobIDGenInstance *JobIDGenerator
	once             sync.Once
)

func GetJobIDGenerator() *JobIDGenerator {
	once.Do(func() {
		jobIDGenInstance = &JobIDGenerator{}
	})
	return jobIDGenInstance
}

func (gen *JobIDGenerator) GenerateJobID() string {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()

	now := time.Now()
	currentMinute := now.Minute()

	if gen.lastMinute != currentMinute {
		gen.sequence = 0
		gen.lastMinute = currentMinute
	}

	gen.sequence++
	timestamp := now.Format("20060102.150405.00000")
	jobID := strconv.Itoa(gen.sequence)

	return fmt.Sprintf("%s.%s", timestamp, jobID)
}

func JobID() string {
	return GetJobIDGenerator().GenerateJobID()
}
