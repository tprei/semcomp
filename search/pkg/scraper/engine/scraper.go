package engine

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tprei/semcomp/search/pkg/logger"
)

type (
	runFunc   func(ctx context.Context) (any, error)
	delayFunc func(attempt int) time.Duration

	Node struct {
		run   runFunc
		delay delayFunc

		currentAttempt int
		lastTried      time.Time
		err            error

		fields logrus.Fields
	}

	Scraper struct {
		fields logrus.Fields
		tasks  []*Node

		numWorkers    int
		maxAttempts   int
		timeout       int
		delayAttempts bool

		ctx     context.Context
		jobs    chan *Node
		failed  chan *Node
		results chan any
	}
)

func NewLoggingField(
	key string,
	value any,
) logrus.Fields {
	return logrus.Fields{
		key: value,
	}
}

func NewNode(
	runFunc runFunc,
	fields logrus.Fields,
) *Node {
	return &Node{
		run: func(ctx context.Context) (any, error) {
			return runFunc(ctx)
		},
		delay: QuadraticDelay, // default

		currentAttempt: 0,
		lastTried:      time.Now(),
		err:            nil,

		fields: fields,
	}
}

func (n *Node) WithDelay(delayFunc delayFunc) {
	n.delay = delayFunc
}

func NewScraper(
	ctx context.Context,
	nodes []*Node,
	fixedAttempts,
	delayAttempts bool,
	fields logrus.Fields,
) *Scraper {
	numWorkers := Config.NumWorkers
	if Config.FractionalNumWorkers > 0.0 && Config.FractionalNumWorkers <= 1.0 {
		numWorkers = int(math.Max(float64(1), math.Ceil(float64(len(nodes))*Config.FractionalNumWorkers)))
		numWorkers = int(math.Min(float64(numWorkers), float64(Config.MaxWorkers)))
	}

	maxAttempts := -1
	if fixedAttempts {
		maxAttempts = Config.MaxAttempts
	}

	fields["numWorkers"] = numWorkers
	return &Scraper{
		tasks:         nodes,
		numWorkers:    numWorkers,
		maxAttempts:   maxAttempts,
		delayAttempts: delayAttempts,
		timeout:       Config.Timeout,

		ctx:     ctx,
		jobs:    make(chan *Node, len(nodes)),
		results: make(chan any, len(nodes)),
		failed:  make(chan *Node, len(nodes)),

		fields: fields,
	}
}

func (sc *Scraper) Run() (results []any) {
	log := logger.GetLogger()
	log.WithFields(sc.fields)

	var ctx = sc.ctx
	var cancel context.CancelFunc

	if Config.Timeout != -1 {
		ctx, cancel = context.WithTimeout(
			ctx,
			time.Duration(sc.timeout*int(time.Second)),
		)
	} else {
		ctx, cancel = context.WithCancel(
			ctx,
		)
	}

	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("unexpected panic in scraper: %s", r))
			results = nil
			cancel()
			return
		}
	}()

	init := time.Now()

	defer close(sc.jobs)
	defer close(sc.results)
	defer close(sc.failed)

	// prepare workers
	log.Info("launching workers...")
	for i := 0; i < sc.numWorkers; i++ {
		go func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("unexpected panic in job goroutine: %s", r)
					results = nil
					cancel()
					return
				}
			}()
			for {
				select {
				case <-ctx.Done():
					return
				case job := <-sc.jobs:
					if job == nil {
						continue
					}

					log.Debug("starting job")

					// fail job if exceeded retries
					if sc.maxAttempts != -1 && job.currentAttempt > sc.maxAttempts {
						log.Warn("job exceeded max retries, failing...")
						sc.failed <- job
						break
					}

					var result any

					failed := false
					denied := false

					var failedErr error

					if sc.delayAttempts && job.delay == nil {
						panic("invalid configuration, delayAttempts is true but delayFn is nil")
					}

					// only run job if sufficient time has passed
					if !sc.delayAttempts || time.Since(job.lastTried) >= job.delay(job.currentAttempt) {
						var processingErr error

						log.Debug("executing process callback")
						result, processingErr = job.run(ctx)

						if processingErr != nil {
							failed = true // processing failure
							failedErr = fmt.Errorf("%s: %s", "processing error", processingErr.Error())
						}
					} else {
						denied = true // job denied by now
					}

					if denied {
						log.Warn("job denied by now...")
						sc.jobs <- job
					} else if failed {
						log.Error("job failed...")

						job.currentAttempt++
						job.lastTried = time.Now()
						job.err = failedErr

						sc.jobs <- job
					} else {
						log.Debug("job successful")
						sc.results <- result
					}
				}

			}
		}(ctx)
	}

	// send jobs
	log.Debug("sending jobs")
	go func() {
		for _, task := range sc.tasks {
			sc.jobs <- task
		}
	}()

	// receive job results
	results = make([]any, 0)
	failures := make([]*Node, 0)

	// wait for all jobs to fail, succeed or timeout
ConsumeLoop:
	for range sc.tasks {
		select {
		case result := <-sc.results:
			results = append(results, result)
		case failure := <-sc.failed:
			failures = append(failures, failure)
		case <-ctx.Done():
			cancel()
			break ConsumeLoop
		}
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		log.Info("work is done, stopping...")
	} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Error("timeout exceeded, stopping...")
	}

	if len(sc.failed) > 0 {
		log.Warnf("%d jobs failed", len(sc.failed))

		for _, failedJob := range failures {
			if failedJob.err != nil {
				log.Error("job failed..")
			}
		}
	}

	log.Infof("scraper finished, time elapsed: %v", time.Since(init))
	return results
}
