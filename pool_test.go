package goconcurrentpool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testJob struct {
	Name     string
	Duration time.Duration // duration of work
	Result   int           // the result that the job should return
}

func (j *testJob) Run() int {
	time.Sleep(j.Duration)
	return j.Result
}

func emptyJob() (any, error) {
	return nil, nil
}

func TestPool(t *testing.T) {

	testJobs := []testJob{
		{
			Name:     "w1",
			Duration: time.Second,
			Result:   10,
		},
		{
			Name:     "w2",
			Duration: time.Millisecond * 100,
			Result:   200,
		},
		{
			Name:     "w3",
			Duration: time.Second * 2,
			Result:   200,
		},
	}

	cp := New(3)

	err := cp.RunJob(emptyJob)
	assert.Error(t, err)
	assert.ErrorIs(t, ErrPoolNotActive, err)

	cp.Run()

	total := 0

	for _, job := range testJobs {
		err := cp.RunJob(func() (any, error) {
			r := job.Run()
			total = total + r
			return r, nil
		})

		assert.NoError(t, err)
	}

	cp.WaitAndClose()

	assert.Equal(t, 410, total)

	err = cp.RunJob(emptyJob)

	assert.Error(t, err)
	assert.ErrorIs(t, ErrPoolNotActive, err)

}
