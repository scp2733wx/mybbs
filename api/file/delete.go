package file

import "os"

func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
