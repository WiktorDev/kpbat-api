package utils

import "os"

func CreateDirectory(name string) error {
	return os.Mkdir("resources/"+name, 0755)
}
