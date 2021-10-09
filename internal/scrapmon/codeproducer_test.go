package scrapmon_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	customNumber "github.com/slysterous/custom-number"

	"github.com/slysterous/scrapmon/internal/scrapmon"
	scrapmonmock "github.com/slysterous/scrapmon/internal/scrapmon/mock"
)

var CustomNumberDigitValues = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func TestProduceCodes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		counter := 0
		mockLogger := scrapmonmock.NewMockLogger(mockCtrl)
		cca := scrapmon.ConcurrentCodeAuthority{
			Logger: mockLogger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		number := customNumber.NewNumber(CustomNumberDigitValues, "aa")

		mockLogger.EXPECT().Infof("Initializing Code Production...! \n").Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 1, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 2, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 3, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 4, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 5, 5).Times(1)

		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "aa").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ab").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ac").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ad").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ae").Times(1)
		mockLogger.EXPECT().Debugf("Finished producing Codes!\n").Times(1)

		codes, _ := cca.Produce(ctx, number, 5, 5)
		var producedCodes []string
		wantCodes := []string{"aa", "ab", "ac", "ad", "ae"}

		for code := range codes {
			counter++
			producedCodes = append(producedCodes, code)
			if counter == 5 {
				cancel()
			}
		}
		if !reflect.DeepEqual(producedCodes, wantCodes) {
			t.Errorf("expected: %v, got: %v", wantCodes, producedCodes)
		}
	})
	t.Run("Success with +1", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		counter := 0
		mockLogger := scrapmonmock.NewMockLogger(mockCtrl)
		cca := scrapmon.ConcurrentCodeAuthority{
			Logger: mockLogger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		number := customNumber.NewNumber(CustomNumberDigitValues, "aa")

		mockLogger.EXPECT().Infof("Initializing Code Production...!").Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 1, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 2, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 3, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 4, 5).Times(1)
		mockLogger.EXPECT().Debugf("Iterations Counter: %d, Desired Iterations: %d\n", 5, 5).Times(2)

		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "aa").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ab").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ac").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ad").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "ae").Times(1)
		mockLogger.EXPECT().Debugf("Producing Code: %s\n", "af").Times(1)
		mockLogger.EXPECT().Debugf("Finished producing Codes!\n").Times(1)

		codes, produceMoreCodes := cca.Produce(ctx, number, 5, 5)
		var producedCodes []string
		wantCodes := []string{"aa", "ab", "ac", "ad", "ae", "af"}

		for code := range codes {
			counter++
			producedCodes = append(producedCodes, code)
			if counter == 3 {
				produceMoreCodes <- struct{}{}
			}
			if counter == 6 {
				cancel()
			}
		}
		if !reflect.DeepEqual(producedCodes, wantCodes) {
			t.Errorf("expected: %v, got: %v", wantCodes, producedCodes)
		}
	})
}

func TestFilterCodes(t *testing.T) {
	t.Run("Success - no codes found", func(t *testing.T) {
		//	t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		counter := 0
		mockLogger := scrapmonmock.NewMockLogger(mockCtrl)
		cca := scrapmon.ConcurrentCodeAuthority{
			Logger: mockLogger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockDM := scrapmonmock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmonmock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		unfilteredCodes := make(chan string, 5)
		produceMoreCodes := make(chan struct{}, 5)
		wantCodes := []string{"aa", "ab", "ac", "ad", "ae"}

		//feed codes
		for _, code := range wantCodes {
			mockDM.EXPECT().CodeAlreadyExists(code).Return(false, nil).Times(1)
			mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", code).Times(1)
			unfilteredCodes <- code
		}

		filteredCodes, _ := cca.Filter(ctx, mockStorage, unfilteredCodes, produceMoreCodes, 5)

		var producedCodes []string

		for code := range filteredCodes {
			counter++
			producedCodes = append(producedCodes, code)

			if counter == 5 {
				cancel()
				close(unfilteredCodes)
				close(produceMoreCodes)

			}
		}
		if !reflect.DeepEqual(producedCodes, wantCodes) {
			t.Errorf("expected: %v, got: %v", wantCodes, producedCodes)
		}
	})
	t.Run("Success - code found", func(t *testing.T) {
		//t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		counter := 0
		mockLogger := scrapmonmock.NewMockLogger(mockCtrl)
		cca := scrapmon.ConcurrentCodeAuthority{
			Logger: mockLogger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockDM := scrapmonmock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmonmock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		unfilteredCodes := make(chan string, 5)
		produceMoreCodes := make(chan struct{}, 5)
		wantCodes := []string{"aa", "ab", "ac", "ad"}

		mockDM.EXPECT().CodeAlreadyExists("aa").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "aa").Times(1)
		unfilteredCodes <- "aa"
		mockDM.EXPECT().CodeAlreadyExists("ab").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "ab").Times(1)
		unfilteredCodes <- "ab"
		mockDM.EXPECT().CodeAlreadyExists("ac").Return(true, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s already exists, asking for another code\n", "ac").Times(1)
		unfilteredCodes <- "ac"
		mockDM.EXPECT().CodeAlreadyExists("ad").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "ad").Times(1)
		unfilteredCodes <- "ad"
		mockDM.EXPECT().CodeAlreadyExists("ae").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "ae").Times(1)
		unfilteredCodes <- "ae"

		filteredCodes, _ := cca.Filter(ctx, mockStorage, unfilteredCodes, produceMoreCodes, 10)

		var producedCodes []string
		for code := range filteredCodes {
			counter++
			producedCodes = append(producedCodes, code)
			if counter == 4 {
				close(unfilteredCodes)
				cancel()
			}
		}

		if !reflect.DeepEqual(producedCodes, []string{"aa", "ab", "ad", "ae"}) {
			t.Errorf("expected: %v, got: %v", wantCodes, producedCodes)
		}
	})
	t.Run("Failed", func(t *testing.T) {
		//t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		counter := 0
		mockLogger := scrapmonmock.NewMockLogger(mockCtrl)
		cca := scrapmon.ConcurrentCodeAuthority{
			Logger: mockLogger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockDM := scrapmonmock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmonmock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		unfilteredCodes := make(chan string, 5)
		produceMoreCodes := make(chan struct{}, 5)
		wantCodes := []string{"aa", "ab", "ac", "ad"}

		mockDM.EXPECT().CodeAlreadyExists("aa").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "aa").Times(1)
		unfilteredCodes <- "aa"
		mockDM.EXPECT().CodeAlreadyExists("ab").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "ab").Times(1)
		unfilteredCodes <- "ab"
		mockDM.EXPECT().CodeAlreadyExists("ac").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "ac").Times(1)
		unfilteredCodes <- "ac"
		mockDM.EXPECT().CodeAlreadyExists("ad").Return(false, nil).Times(1)
		mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", "ad").Times(1)
		unfilteredCodes <- "ad"
		mockDM.EXPECT().CodeAlreadyExists("ae").Return(false, errors.New("test error")).Times(1)
		unfilteredCodes <- "ae"

		filteredCodes, errC := cca.Filter(ctx, mockStorage, unfilteredCodes, produceMoreCodes, 5)

		var producedCodes []string
		for code := range filteredCodes {
			counter++
			producedCodes = append(producedCodes, code)
			if counter == 4 {
				close(unfilteredCodes)
				close(produceMoreCodes)
				cancel()
			}
		}

		for err := range errC {
			if err == nil {
				t.Error("Expected error got nil")
			}
		}

		for range produceMoreCodes {
			fmt.Printf("more asked!")
		}

		if !reflect.DeepEqual(producedCodes, wantCodes) {
			t.Errorf("expected: %v, got: %v", wantCodes, producedCodes)
		}
	})
	t.Run("Context Cancelled", func(t *testing.T) {
		//t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockLogger := scrapmonmock.NewMockLogger(mockCtrl)
		cca := scrapmon.ConcurrentCodeAuthority{
			Logger: mockLogger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockDM := scrapmonmock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmonmock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		unfilteredCodes := make(chan string, 5)

		produceMoreCodes := make(chan struct{}, 5)
		defer close(produceMoreCodes)
		wantCodes := []string{"aa", "ab", "ac", "ad", "ae"}

		go func(wantCodes []string) {
			// async feed codes
			defer close(unfilteredCodes)
			for _, code := range wantCodes {
				mockDM.EXPECT().CodeAlreadyExists(code).Return(false, nil).Times(1)
				mockLogger.EXPECT().Debugf("File with code %s does not exist, will be downloaded\n", code).Times(1)

				unfilteredCodes <- code
				time.Sleep(time.Millisecond * 500)
			}
		}(wantCodes)

		mockLogger.EXPECT().Debugf("Finished Filtering Codes!\n").Times(1)

		filteredCodes, _ := cca.Filter(ctx, mockStorage, unfilteredCodes, produceMoreCodes, 1)

		var producedCodes []string

		for code := range filteredCodes {
			producedCodes = append(producedCodes, code)
			cancel()
		}

	})

}
