package scrapmon

import (
	"context"
	"time"

	customNumber "github.com/slysterous/custom-number"
)

//go:generate mockgen -destination mock/codeproducer.go -package scrapmon_mock . ConcurrentCodeProducer

// ConcurrentCodeProducer describes the actions of a code producer.
type ConcurrentCodeProducer interface {
	Produce(
		ctx context.Context,
		index customNumber.Number,
		iterations int,
		channelSize int,
	) (<-chan string, chan struct{})
	Filter(
		ctx context.Context,
		storage Storage,
		codes <-chan string,
		produceMoreCodes chan<- struct{},
		channelSize int,
	) (<-chan string, <-chan error)
}

// ConcurrentCodeAuthority is responsible for code creation and handling.
type ConcurrentCodeAuthority struct {
	Logger   Logger
	Scrapper Scrapper
}

// Produce produces codes in a channel while handling feedback.
func (cca ConcurrentCodeAuthority) Produce(
	ctx context.Context,
	index customNumber.Number,
	iterations int,
	channelSize int,
) (<-chan string, chan struct{}) {
	produceMoreCodes := make(chan struct{}, iterations+1)
	codes := make(chan string, iterations+1)
	iterationsCounter := 0
	cca.Logger.Info("Initializing Code Production...! \n")
	go func() {
		defer close(codes)
		defer close(produceMoreCodes)

		for {
			if iterationsCounter < iterations {
				produceMoreCodes <- struct{}{}
				iterationsCounter++
			}

			select {
			case <-produceMoreCodes:
				cca.Logger.Debugf("Iterations Counter: -%d-, Desired Iterations: %d \n", iterationsCounter, iterations)
				cca.Logger.Debugf("Producing Code: %s \n", index.String())
			codesFor:
				for {
					select {
					case codes <- index.String():
						index.Increment()
						break codesFor
					case <-time.After(1 * time.Second):
						cca.Logger.Warnf("Deadlock Timeout on PRODUCE_CODES!")
						break
					}
				}
			case <-ctx.Done():
				cca.Logger.Debugf("Finished producing Codes!\n")
				return
			}
		}
	}()
	return codes, produceMoreCodes
}

// Filter filters codes depending if they exist in Storage or not.
func (cca ConcurrentCodeAuthority) Filter(
	ctx context.Context,
	storage Storage,
	codes <-chan string,
	produceMoreCodes chan<- struct{},
	channelSize int,
) (<-chan string, <-chan error) {
	filteredCodes := make(chan string, channelSize)
	errc := make(chan error, 1)

	go func() {
		defer close(filteredCodes)
		defer close(errc)

		for code := range codes {
			exists, err := storage.Dm.CodeAlreadyExists(code)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			if exists {
				//if code exists then
			moreCodes:
				for { //keep trying to add an item in produceMoreCodes
					select {
					case produceMoreCodes <- struct{}{}:
						cca.Logger.Debugf("File with code %s already exists, asking for another code\n", code)
						break moreCodes
					case <-time.After(1 * time.Second):
						// If it takes more than 1 second then keep trying
						cca.Logger.Warnf("Deadlock Timeout on FILTER_CODES!")
						break
					}
				}
				continue
			}
			cca.Logger.Debugf("File with code %s does not exist, will be downloaded\n", code)

			select {
			case filteredCodes <- code:
			case <-ctx.Done():
				cca.Logger.Debugf("Finished Filtering Codes!\n")
				return
			}
		}
	}()
	return filteredCodes, errc
}
