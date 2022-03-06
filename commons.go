package kubeprobes

import "sync"

// ProbeFunction is a function that determines whether
// the given metric may be marked as correctly functioning.
// It not, the error should be returned.
type ProbeFunction func() error

type statusQuery struct {
	allGreen bool
	mux      sync.Mutex
	wg       sync.WaitGroup
}

func (sq *statusQuery) isAllGreen() bool {
	sq.wg.Wait()
	sq.mux.Lock()
	defer sq.mux.Unlock()
	return sq.allGreen
}

func newStatusQuery(probes []ProbeFunction) *statusQuery {
	sq := &statusQuery{
		allGreen: true,
		mux:      sync.Mutex{},
		wg:       sync.WaitGroup{},
	}

	sq.wg.Add(len(probes))
	for _, probe := range probes {
		probe := probe
		go func() {
			defer sq.wg.Done()
			if err := probe(); err != nil {
				sq.mux.Lock()
				sq.allGreen = false
				sq.mux.Unlock()
			}
		}()
	}

	return sq
}
