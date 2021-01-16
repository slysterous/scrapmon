package http_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	shttp "github.com/slysterous/scrapmon/internal/http"
	httpMock "github.com/slysterous/scrapmon/internal/http/mock"
)

func TestScrapeByCode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDownloader := httpMock.NewMockDownloader(mockCtrl)
		mockReader := httpMock.NewMockReader(mockCtrl)

		scrapper := shttp.NewClient("baseball/", mockReader, mockDownloader)

		mockDownloader.EXPECT().Get("baseball/file.png").Return(&http.Response{
			Status:           "",
			StatusCode:       200,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           nil,
			Body:             http.NoBody,
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            false,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          nil,
			TLS:              nil,
		}, nil).Times(1)
		_, err := scrapper.ScrapeByCode("file", "png")
		if err != nil {
			t.Errorf("unexpected error occured, err:%v", err)
		}
	})
	t.Run("Get Error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDownloader := httpMock.NewMockDownloader(mockCtrl)
		mockReader := httpMock.NewMockReader(mockCtrl)

		scrapper := shttp.NewClient("baseball/", mockReader, mockDownloader)

		mockDownloader.EXPECT().Get("baseball/file.png").Return(nil, errors.New("test error")).Times(1)
		_, err := scrapper.ScrapeByCode("file", "png")
		if err == nil {
			t.Errorf("expected error got nil")
		}

	})
	t.Run("Not Found or Removed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDownloader := httpMock.NewMockDownloader(mockCtrl)
		mockReader := httpMock.NewMockReader(mockCtrl)

		scrapper := shttp.NewClient("baseball/", mockReader, mockDownloader)

		mockDownloader.EXPECT().Get("baseball/file.png").Return(&http.Response{
			Status:           "",
			StatusCode:       404,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           nil,
			Body:             http.NoBody,
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            false,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          nil,
			TLS:              nil,
		}, nil).Times(1)
		file, err := scrapper.ScrapeByCode("file", "png")
		if err != nil {
			t.Errorf("unexpected error occured, err:%v", err)
		}
		if file.Code != "" {
			t.Errorf("expected code to be empty string got: %s", file.Code)
		}
	})
}
