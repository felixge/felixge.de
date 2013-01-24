package main

import (
	"crypto/md5"
	"fmt"
	"github.com/felixge/felixge.de/generator"
	"io"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"mime"
	"net/http"
	"os"
	gopath "path"
	"strconv"
)

func main() {
	auth := aws.Auth{os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET")}
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket(os.Getenv("S3_BUCKET"))

	log.Printf("Listing s3 bucket ...")
	s3Files, err := listS3(bucket)
	if err != nil {
		panic(err)
	}

	log.Printf("Listing local fs ...")
	fs := generator.NewFs()
	localFiles, err := listFs(fs, "/")
	if err != nil {
		panic(err)
	}

	for localPath, localFile := range localFiles {
		put := false

		if s3File, ok := s3Files[localPath]; !ok {
			fmt.Printf("A %s\n", localPath)
			put = true
		} else if s3File.Md5 != localFile.Md5 {
			fmt.Printf("M: %s (%s vs %s)\n", localPath, localFile.Md5, s3File.Md5)
			put = true
		} else {
			fmt.Printf("U: %s\n", localPath)
		}

		if put {
			err := putS3(fs, localPath, bucket)
			if err != nil {
				panic(err)
			}
		}
	}

	for s3Path, _ := range s3Files {
		if _, ok := localFiles[s3Path]; !ok {
			err := delS3(s3Path, bucket)
			if err != nil {
				panic(err)
			}
		}
	}
}

func listS3(bucket *s3.Bucket) (map[string]fileInfo, error) {
	// BUG: This will only return up to 1000 files, which may not be all files
	res, err := bucket.List("", "", "", 1000)
	if err != nil {
		return nil, err
	}

	results := make(map[string]fileInfo, len(res.Contents))
	for _, key := range res.Contents {
		path := "/" + key.Key
		md5, err := strconv.Unquote(key.ETag)
		if err != nil {
			return nil, err
		}

		results[path] = fileInfo{Path: path, Md5: md5}
	}
	return results, nil
}

func listFs(fs http.FileSystem, path string) (map[string]fileInfo, error) {
	dir, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	stats, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	results := make(map[string]fileInfo, len(stats))

	for _, stat := range stats {
		path := gopath.Join(path, stat.Name())

		if stat.IsDir() {
			dirResults, err := listFs(fs, path)
			if err != nil {
				return nil, err
			}

			for path, result := range dirResults {
				results[path] = result
			}
			continue
		}

		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		hash := md5.New()
		if _, err = io.Copy(hash, file); err != nil {
			return nil, err
		}

		results[path] = fileInfo{
			Path: path,
			Md5:  fmt.Sprintf("%x", hash.Sum(nil)),
		}
	}
	return results, nil
}

func putS3(fs http.FileSystem, path string, bucket *s3.Bucket) error {
	log.Printf("Uploading %s ...", path)

	file, err := fs.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil
	}

	size := stat.Size()
	if size < 0 {
		return fmt.Errorf("could not determine size for: %s", path)
	}

	ext := gopath.Ext(path)
	mimeType := mime.TypeByExtension(ext)

	return bucket.PutReader(path[1:], file, size, mimeType, s3.PublicRead)
}

func delS3(path string, bucket *s3.Bucket) error {
	log.Printf("Deleting %s ...", path)

	return bucket.Del(path[1:])
}

type fileInfo struct {
	Path string
	Md5  string
}
