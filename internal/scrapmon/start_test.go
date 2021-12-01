package scrapmon_test

import (
	"testing"
)

func TestConcurrentCommandManagerStartCommand(t *testing.T) {
	//t.Run("Success",func(t *testing.T){
	//	mockCtrl := gomock.NewController(t)
	//	defer mockCtrl.Finish()
	//
	//	mockFm := scrapmon_mock.NewMockFileManager(mockCtrl)
	//	mockDm := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
	//
	//	mockStorage := scrapmon.Storage{
	//		Fm: mockFm,
	//		Dm: mockDm,
	//	}
	//	mockLogger := scrapmonmock.NewLogger()
	//	mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)
	//	commandManager := scrapmon.ConcurrentCommandManager{
	//		Storage: mockStorage,
	//		CodeAuthority: scrapmon.ConcurrentCodeAuthority{
	//			Logger:   mockLogger,
	//			Scrapper: mockScrapper,
	//		},
	//	}
	//
	//	mockDm.EXPECT().CodeAlreadyExists("0").Return(false,nil).Times(1)
	//	mockDm.EXPECT().CodeAlreadyExists("1").Return(false,nil).Times(1)
	//	mockDm.EXPECT().CodeAlreadyExists("2").Return(false,nil).Times(1)
	//	commandManager.StartCommand("0",3,3)
	//})
}
