package phantom

import "flag"

func InitFlag() *string {
	filePath := flag.String("filePath", "./", "file config path")
	flag.Parse()
	return filePath
}
