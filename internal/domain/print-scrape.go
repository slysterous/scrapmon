package domain


// Config represents the applications configuration parameters
type Config struct{
	Env               string
	DatabaseUser      string
	DatabasePassword  string
	DatabaseHost      string
	DatabasePort      string
	DatabaseName      string
	HTTPClientTimeout int
	MaxDBConnections  int
}

// type ScrapSaver interface {
// 	SaveImageToFileSystem()
// 	SaveImageRefToDB()
// 	SaveNonExistentImageRefToDB()
// }

// type ScrapGetter interface {
// 	GetImageByCode()
// }

// type ScrapImage struct {
// 	PrntCode   string
// 	StorageRef string
// }
