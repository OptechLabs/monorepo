package foundation

import (
	"context"
	"log"
	"sync"
)

type Processor interface {
	Start(ctx context.Context) (err error)
	Stop(wg *sync.WaitGroup) (err error)
}

func (f *Foundation) AddProcessor(p Processor) {
	f.processors = append(f.processors, p)
}

func (f *Foundation) StopProcessors(wg *sync.WaitGroup) (errs []error) {
	for _, p := range f.processors {
		wg.Add(1)
		err := p.Stop(wg)
		if err != nil {
			log.Printf("[foundation] ERROR: unable to stop processor: %s", err)
			errs = append(errs, err)
		}
	}
	return
}

func (f *Foundation) StartProcessors() (errs []error) {
	ctx := context.Background()
	for _, p := range f.processors {
		err := p.Start(ctx)
		if err != nil {
			log.Printf("[foundation] ERROR: unable to start processor: %s", err)
			errs = append(errs, err)
		}
	}
	return
}
