package goworkerpool

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

	wp := NewWorkerPool(3)

	total := 0

	for _, job := range testJobs {
		wp.RunJob(func() {
			r := job.Run()
			total = total + r
		})
	}

	wp.StopAndWait()

	assert.Equal(t, 410, total)

}

// func main() {

// testing.T

// works := []Work{
// {
//  Name:    "w4",
//  Amount:  200,
//  Timeout: time.Millisecond * 150,
// },
// {
//  Name:    "w5",
//  Amount:  100,
//  Timeout: time.Millisecond * 100,
// },
// {
//  Name:    "w6",
//  Amount:  5,
//  Timeout: time.Millisecond * 50,
// },
// {
//  Name:    "w7",
//  Amount:  6,
//  Timeout: time.Millisecond * 150,
// },
// {
//  Name:    "w8",
//  Amount:  10,
//  Timeout: time.Millisecond * 100,
// },
// {
//  Name:    "w9",
//  Amount:  5,
//  Timeout: time.Millisecond * 50,
// },
// {
//  Name:    "w10",
//  Amount:  6,
//  Timeout: time.Millisecond * 150,
// },
// {
//  Name:    "w11",
//  Amount:  5,
//  Timeout: time.Millisecond * 50,
// },
// {
//  Name:    "w12",
//  Amount:  6,
//  Timeout: time.Millisecond * 150,
// },
// {
//  Name:    "w13",
//  Amount:  100,
//  Timeout: time.Millisecond * 50,
// },
// {
//  Name:    "w14",
//  Amount:  20,
//  Timeout: time.Millisecond * 300,
// },
// {
//  Name:    "w15",
//  Amount:  150,
//  Timeout: time.Millisecond * 100,
// },
// {
//  Name:    "w16",
//  Amount:  6,
//  Timeout: time.Millisecond * 150,
// },
// }

// withTimeMetric("ParPool", func() {
// 	fmt.Println(ParallelPool(works))
// })

// withTimeMetric("Seq", func() {
// 	fmt.Println(Seq(works))

// })
// }

// type Work struct {
// 	Name    string
// 	Amount  int
// 	Timeout time.Duration
// }

// func Seq(works []Work) any {
// 	res := 0
// 	for _, w := range works {
// 		r := someWork(w)
// 		res += r
// 	}

// 	return res
// }

// func ParallelPool(works []Work) any {
// 	p := NewPool(5)
// 	res := 0
// 	for i := 0; i < len(works); i++ {
// 		w := works[i]
// 		p.RunWork(func() {
// 			r := someWork(w)
// 			res += r
// 		})
// 	}

// 	p.Wait()

// 	return res
// }

// func someWork(w Work) int {
// 	fmt.Println("start work: ", w.Name, w.Amount)
// 	for i := 0; i < w.Amount; i++ {
// 		time.Sleep(w.Timeout)
// 	}

// 	return w.Amount * 2
// }

// func withTimeMetric(comment string, f func()) {

// 	fmt.Println("Start: ", comment)

// 	timeStart := time.Now()
// 	f()
// 	dur := time.Since(timeStart)

// 	fmt.Printf("[%s]: elapsed time: %s \n", comment, dur)
// }
