package mock

//type FileManager struct {
//	SaveFileFn func(src scrapmon.ScrapedFile) error
//	PurgeFn func() error
//	SaveFileCalls int
//	PurgeCalls int
//}
//
//func (fm FileManager) SaveFile(src scrapmon.ScrapedFile) error {
//	if fm.SaveFileFn != nil {
//		fm.SaveFileCalls++
//		return fm.SaveFileFn(src)
//	}
//	return nil
//}
//
//func (fm FileManager) Purge() error {
//	if fm.PurgeFn != nil {
//		fm.PurgeCalls++
//		return fm.PurgeFn()
//	}
//	return nil
//}