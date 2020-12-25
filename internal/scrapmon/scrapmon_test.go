package domain_test

import (
	"github.com/slysterous/scrapmon/internal/mock"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	"testing"
	"errors"
)

func TestStoragePurge(t *testing.T) {
	t.Run("Success",func(t *testing.T) {
		mockPurger:=mock.Purger{
			PurgeFn: func() error {
				return nil
			},
		}

		mockFm:=mock.FileManager{
			Purger:        mockPurger,
		}

		mockDm:=mock.DatabaseManager{
			Purger:        mockPurger,
		}
		storage:= scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}
		err := storage.Purge()
		if err !=nil {
			t.Errorf("unexpected error occured, err: %v",err)
		}
	})
	t.Run("Failure on Database", func(t *testing.T) {
		dbPurger:=mock.Purger{
			PurgeFn: func() error {
				return errors.New("test DB error")
			},
		}
		fmPurger:=mock.Purger{
			PurgeFn: func() error {
				return errors.New("test FM error")
			},
		}

		mockFm:=mock.FileManager{
			Purger:        fmPurger,
		}

		mockDm:=mock.DatabaseManager{
			Purger:        dbPurger,
		}
		storage:= scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}
		err := storage.Purge()
		if err ==nil {
			t.Error("expected error got nil",err)
		}
		if err.Error()!="test DB error" {
			t.Errorf("wanted: test DB error got: %v",err)
		}

	})
	t.Run("Failure on FileStorage", func(t *testing.T) {
		dbPurger:=mock.Purger{
			PurgeFn: func() error {
				return nil
			},
		}
		fmPurger:=mock.Purger{
			PurgeFn: func() error {
				return errors.New("test FM error")
			},
		}

		mockFm:=mock.FileManager{
			Purger:        fmPurger,
		}

		mockDm:=mock.DatabaseManager{
			Purger:        dbPurger,
		}
		storage:= scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}
		err := storage.Purge()
		if err ==nil {
			t.Fatal("expected error got nil",err)
		}
		if err.Error()!="test FM error" {
			t.Errorf("wanted: test FM error got: %v",err)
		}

	})
}

func TestCommandManagerPurge(t *testing.T){
	t.Run("Success", func(t *testing.T){
		mockStorage:= scrapmon.Storage{
			Fm:         mock.FileManager{},
			Dm:         mock.DatabaseManager{},
		}
		mockScrapper:=mock.Scrapper{
			ScrapeByCodeFn: func(code string) (scrapmon.ScrapedFile, error) {
				return scrapmon.ScrapedFile{},nil
			},
			ScrapeByCodeCalls: 0,
		}
		cm:=scrapmon.CommandManager{
			Storage:  mockStorage,
			Scrapper: mockScrapper,
		}
		err:= cm.PurgeCommand()
		if err !=nil{
			t.Errorf("unexpected error occured, err:= %v",err)
		}
	})
	t.Run("Failure", func(t *testing.T){
		mockStorage := scrapmon.Storage{
			Fm: mock.FileManager{
				Purger:        mock.Purger{
					PurgeFn: func() error {
						return errors.New("test error from FM Purge")
					},
					PurgeCalls: 0,
				},
			},
			Dm: mock.DatabaseManager{},
		}
		mockScrapper := mock.Scrapper{
			ScrapeByCodeFn: func(code string) (scrapmon.ScrapedFile, error) {
				return scrapmon.ScrapedFile{}, nil
			},
			ScrapeByCodeCalls: 0,
		}
		cm := scrapmon.CommandManager{
			Storage:  mockStorage,
			Scrapper: mockScrapper,
		}
		err := cm.PurgeCommand()
		if err == nil {
			t.Fatal("expected error got nil")
		}
		if err.Error() !="could not purge storage, err: test error from FM Purge"{
			t.Errorf("wanted: could not purge storage, err: test error from FM Purge, got: %v",err)
		}
	})
}