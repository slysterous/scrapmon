package scrapmon

import (
	"context"
	customNumber "github.com/slysterous/custom-number"
	"time"
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

// ConcurrentCodeAuthority is responsible for code creation and handling.
type ConcurrentCodeAuthority struct {}

// Produce produces codes in a channel while handling feedback.
func (cca ConcurrentCodeAuthority) Produce(
	ctx context.Context,
	index customNumber.Number,
	iterations int,
	channelSize int,
) (<-chan string, chan struct{}){
	produceMoreCodes := make(chan struct{}, iterations+1)
	codes := make(chan string, iterations+1)
	iterationsCounter := 0
	//fmt.Printf("PRODUCING CODES")
	go func() {
		defer close(codes)
		defer close(produceMoreCodes)

		for {
			if iterationsCounter < iterations {
				produceMoreCodes <- struct{}{}
				iterationsCounter++
			}
			//fmt.Println("IM HERE")

			select {
			case <-produceMoreCodes:
				//fmt.Printf("iterationsCounter: %d  iterations: %d\n", iterationsCounter, iterations)
				//fmt.Printf("PRODUCING CODE: %s \n", index.SmartString())
			codesFor:
				for {
					select {
					case codes <- index.String():
						index.Increment()
						break codesFor
					case <-time.After(1 * time.Second):
						//fmt.Println("DEADLOCK TIMEOUT PRODUCE CODES")
						break
					}
				}
			case <-ctx.Done():
				//fmt.Printf("CONTEXT DONE on produce codes")
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
	
				for1:
					for {
						select {
						case produceMoreCodes <- struct{}{}:
							break for1
						case <-time.After(1 * time.Second):
							//fmt.Println("DEADLOCK TIMEOUT PRODUCE MORE CODES")
							break
						}
					}
					//fmt.Printf("Image %s already exists, asking for another code.\n", code)
					continue
				}
				//fmt.Printf("Image %s does not exist, image will be downloaded.\n", code)
	
				select {
				case filteredCodes <- code:
				case <-ctx.Done():
					//fmt.Printf("CONTEXT DONE")
					return
				}
			}
		}()
		return filteredCodes, errc
	}
}