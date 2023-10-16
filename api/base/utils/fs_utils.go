package utils

import (
	"fmt"
	"os"
	"strings"
)

func CreateDirectory(name string) error {
	return os.Mkdir("resources/"+name, 0755)
}
func RemoveImage(categoryId int, image string) bool {
	err := os.Remove(fmt.Sprintf("resources/category_%d/%s", categoryId, image))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
func RemoveDir(dir string) {
	os.Remove(fmt.Sprintf("resources/%s", strings.ToLower(dir)))
}
