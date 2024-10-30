package server

import "net/http"

type modFileSystem struct {
	fs  http.FileSystem
	app application
}

func (mfs modFileSystem) Open(path string) (http.File, error) {
	f, err := mfs.fs.Open(path)
	if err != nil {
		mfs.app.errorLog.Println(err)
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		mfs.app.errorLog.Println(err)
		return nil, err
	}

	if s.IsDir() {
		mfs.app.infoLog.Println(405, "Method Not Allowed", "/assets"+path)
		if _, err := mfs.fs.Open("/index.html"); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				mfs.app.errorLog.Println(err)
				return nil, closeErr
			}
			mfs.app.errorLog.Println(err)
			return nil, err
		}
	} else {
		mfs.app.infoLog.Println(200, "OK", "/assets"+path)
	}

	return f, nil
}
