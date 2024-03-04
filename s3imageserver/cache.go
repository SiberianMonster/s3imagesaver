package s3imageserver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (i *Image) getFromCache(r *http.Request) (err error) {
	newFileName := i.getCachedFileName(r)
	info, err := os.Stat(newFileName)
	if err != nil {
		return err
	}
	if (time.Duration(i.CacheTime))*time.Second > time.Since(info.ModTime()) {
		f, err := os.Open(newFileName)
		if err != nil {
			return err
		}
		file, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return err
		}
		i.Image = file
		if i.Debug {
			fmt.Println("from cache")
		}
		return nil
	}
	go removeExpiredImage(newFileName)
	return errors.New("The file has expired")
}

func (i *Image) writeCache(r *http.Request) {
	err := ioutil.WriteFile(i.getCachedFileName(r), i.Image, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func removeExpiredImage(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		fmt.Println(err)
	}
}

func (i *Image) getCachedFileName(r *http.Request) (fileName string) {

	return fmt.Sprintf("%v/%v.png", i.CachePath, i.FileName)
}

// TODO: add garbage colection
// TODO: add documentation
