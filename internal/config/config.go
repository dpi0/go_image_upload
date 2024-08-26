package config

import (
	"os"
)

const UploadDir = "./uploads"

func InitConfig() {
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		os.Mkdir(UploadDir, 0755)
	}
}
