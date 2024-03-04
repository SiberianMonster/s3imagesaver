package s3imageserver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xiam/to"

)

type Image struct {
	Path            string
	FileName        string
	Bucket          string
	Crop            bool
	Debug           bool
	Height          int
	Width           int
	Image           []byte
	CacheTime       int
	CachePath       string
	ErrorImage      string
	OutputFormat    string
}

var allowedTypes = []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}

func NewImage(r *http.Request, config HandlerConfig, fileName string) (image *Image, err error) {
	maxDimension := 3064

	crop := false
	if r.URL.Query().Get("c") != "" {
		crop = to.Bool(r.URL.Query().Get("c"))
	}
	image = &Image{
		Path:            config.Timeweb.FilePath,
		Bucket:          config.Timeweb.BucketName,
		TimewebToken:	 config.Timeweb.TimewebToken,
		Height:          100,
		Crop:            crop,
		Width:           100,
		CacheTime:       604800, // cache time in seconds, set 0 to infinite and -1 for disabled
		CachePath:       config.CachePath,
		ErrorImage:      "",
		OutputFormat:    ".png",
	}
	if config.CacheTime != nil {
		image.CacheTime = *config.CacheTime
	}
	
	if image.FileName == "" {
		err = errors.New("File name cannot be an empty string")
	}
	if image.Bucket == "" {
		err = errors.New("Bucket cannot be an empty string")
	}
	return image, err
}

func (i *Image) getImage(w http.ResponseWriter, r *http.Request) {
	var err error
	if i.CacheTime > -1 {
		err = i.getFromCache(r)
	} else {
		err = errors.New("Caching disabled")
	}
	if err != nil {
		fmt.Println(err)
		err := i.getImageFromS3()
		if err != nil {
			fmt.Println(err)
			err = i.getErrorImage()
			w.WriteHeader(404)
		} else {
			go i.writeCache(r)
		}
	}
	i.write(w)
}


func (i *Image) write(w http.ResponseWriter) {
	w.Header().Set("Content-Length", strconv.Itoa(len(i.Image)))
	w.Write(i.Image)
}

func (i *Image) getErrorImage() (err error) {
	if i.ErrorImage != "" {
		i.Image, err = ioutil.ReadFile(i.ErrorImage)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Error image not specified")
}

func (i *Image) getImageFromS3() (err error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://s3.timeweb.com/%v%v%v", i.Bucket, i.Path, i.FileName), nil)
	req.Header.Set("Authorization", "Bearer " + i.TimewebToken)

	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		i.Image, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else if i.Debug {
			fmt.Println("Retrieved image from from S3")
		}
		return nil
	} else if resp.StatusCode != http.StatusOK {
		err = errors.New("Error while making request")
	}
	return err
}

