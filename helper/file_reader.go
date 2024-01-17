package helper

import "os"

// ReadFile is a helper function to read a file.
// path parameter is where the file is located.
func ReadFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	return file, nil
}
