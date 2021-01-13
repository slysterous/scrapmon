package scrapmon

import (
	"context"
	customNumber "github.com/slysterous/custom-number"
	"log"
)

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

// ConcurrentCodeAuthority
type ConcurrentCodeAuthority struct {}

func (cca ConcurrentCodeAuthority) Produce(
	ctx context.Context,
	index customNumber.Number,
	iterations int,
	channelSize int,
) (<-chan string, chan struct{}){
	produceMoreCodes := make(chan struct{}, iterations+1)
	codes := make(chan string, iterations+1)
	iterationsCounter := 0
	fmt.Printf("PRODUCING CODES")
	go func() {
		defer close(codes)
		defer close(produceMoreCodes)

		for {
			if iterationsCounter < iterations {
				produceMoreCodes <- struct{}{}
				iterationsCounter++
			}
			fmt.Println("IM HERE")

			select {
			case <-produceMoreCodes:
				fmt.Printf("iterationsCounter: %d  iterations: %d\n", iterationsCounter, iterations)
				//fmt.Printf("PRODUCING CODE: %s \n", index.SmartString())
			codesFor:
				for {
					select {
					case codes <- index.String():
						index.Increment()
						break codesFor
					case <-time.After(1 * time.Second):
						fmt.Println("DEADLOCK TIMEOUT PRODUCE CODES")
						break
					}
				}
			case <-ctx.Done():
				fmt.Printf("CONTEXT DONE on produce codes")
				return
			}
		}
	}()
	return codes, produceMoreCodes
}

func (c CodeAuthority) Filter(
	ctx context.Context,
	storage Storage,
	codes <-chan string,
	produceMoreCodes chan<- struct{},
	channelSize int,
) (<-chan string, <-chan error) {

}