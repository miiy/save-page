package file

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ReadAll(name string) ([]byte, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	b, err := ioutil.ReadAll(file)
	return b, err
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	} else {
		return false
	}
}

func Mkdir(name string) error {
	return os.Mkdir("." + string(os.PathSeparator) + name, 0777)
}

func SafeName(str string) string {
	f := [...]string{"/", "\\", ":", "*", "?", "\"", ">", "<", "|", "\n", "\t", "\r"}

	for _, v := range f {
		str = strings.Replace(str, v, " ", -1)
	}
	return str
}
