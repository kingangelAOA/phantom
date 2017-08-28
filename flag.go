package phantom

import "flag"

//InitFlag init flags
func InitFlag() *string {
	config := flag.String("config", "./", "file config path")
	flag.Parse()
	return config
}
