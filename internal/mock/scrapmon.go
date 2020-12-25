package mock

import scrapmon"github.com/slysterous/scrapmon/internal/scrapmon"

type Purger struct {
	PurgeFn func() error
	PurgeCalls int
}

func (p Purger) Purge() error {
	if p.PurgeFn != nil {
		p.PurgeCalls++
		return p.PurgeFn()
	}
	return nil
}

type FileManager struct {
	SaveFileFn func (src scrapmon.ScrapedFile) error
	SaveFileCalls int
	Purger
}

func (fm FileManager) SaveFile(src scrapmon.ScrapedFile) error{
	if fm.SaveFileFn != nil {
		fm.SaveFileCalls ++
		return fm.SaveFileFn(src)
	}
	return  nil
}

type DatabaseManager struct {
	CreateScrapFn func(s scrapmon.Scrap) (int, error)
	UpdateScrapStatusByCodeFn func(code string, status scrapmon.ScrapStatus) error
	UpdateScrapByCodeFn func(s scrapmon.Scrap) error
	GetLatestCreatedScrapCodeFn func() (*string, error)
	CodeAlreadyExistsFn func(code string) (bool, error)
	CreateScrapCalls int
	UpdateScrapStatusByCodeCalls int
	UpdateScrapByCodeCalls int
	GetLatestCreatedScrapCodeCalls int
	CodeAlreadyExistsCalls int
	Purger
}

func (dm DatabaseManager) CreateScrap(s scrapmon.Scrap) (int,error){
	if dm.CreateScrapFn != nil {
		dm.CreateScrapCalls ++
		return dm.CreateScrapFn(s)
	}
	return  -1,nil
}

func (dm DatabaseManager) UpdateScrapStatusByCode(code string, status scrapmon.ScrapStatus) error {
	if dm.UpdateScrapStatusByCodeFn !=nil{
		dm.UpdateScrapStatusByCodeCalls++
		return dm.UpdateScrapStatusByCodeFn(code,status)
	}
	return nil
}

func (dm DatabaseManager) UpdateScrapByCode(s scrapmon.Scrap) error {
	if dm.UpdateScrapByCodeFn !=nil {
		dm.UpdateScrapByCodeCalls++
		return dm.UpdateScrapByCodeFn(s)
	}
	return nil
}

func (dm DatabaseManager) GetLatestCreatedScrapCode() (*string,error) {
	if dm.GetLatestCreatedScrapCodeFn !=nil{
		dm.GetLatestCreatedScrapCodeCalls++
		return dm.GetLatestCreatedScrapCodeFn()
	}
	return nil,nil
}

func (dm DatabaseManager) CodeAlreadyExists(code string) (bool, error) {
	if dm.CodeAlreadyExistsFn !=nil{
		dm.CodeAlreadyExistsCalls++
		return dm.CodeAlreadyExistsFn(code)
	}
	return false,nil
}

type Storage struct {
	Fm FileManager
	Dm DatabaseManager
	PurgeFn func()error
	PurgeCalls int
}

func (s Storage)Purge() error {
	if s.PurgeFn !=nil {
		s.PurgeCalls++
		return s.PurgeFn()
	}
	return nil
}

type Scrapper struct {
	ScrapeByCodeFn func(code string) (scrapmon.ScrapedFile,error)
	ScrapeByCodeCalls int
}

func (s Scrapper)ScrapeByCode(code string)(scrapmon.ScrapedFile,error){
	if s.ScrapeByCodeFn !=nil {
		s.ScrapeByCodeCalls++
		return s.ScrapeByCodeFn(code)
	}
	return scrapmon.ScrapedFile{},nil
}

type CommandManager struct {
	Storage Storage
	Scrapper Scrapper
}