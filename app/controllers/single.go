package controllers

import (
	//"bytes"
	//"fmt"
	//"image"
	//_ "image/jpeg"
	//_ "image/png"

	"compress/gzip"
	"io"
  "log"
  "os"
	//"net/http"
	//"bufio"

	//"strings"
	// "github.com/revel/samples/upload/app/routes"

	"github.com/revel/revel"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

type Single struct {
		App
}

func (c *Single) Upload() revel.Result {
	return c.Render()
}

func (c *Single) HandleUpload(avatar []byte) revel.Result {
	// Validation rules.
	c.Validation.Required(avatar)
	c.Validation.MinSize(avatar, 1*KB).
		Message("Minimum a file size of 1KB expected")
	c.Validation.MaxSize(avatar, 20*MB).
		Message("File cannot be larger than 20MB")

	// // Check format of the file.
	// conf, format, err := image.DecodeConfig(bytes.NewReader(avatar))
	// c.Validation.Required(err == nil).Key("avatar").
	// 	Message("Incorrect file format")
	// c.Validation.Required(format == "pdf" || format == "png").Key("avatar").
	// 	Message("JPEG or PNG file format is expected")

	var name string = c.Params.Files["avatar"][0].Filename
	var path string = "src/myapp/upload/" + name

	dst, err := os.Create("src/myapp/upload/" + name)
	defer dst.Close()
	if err != nil {
		log.Printf("error in creating dst")
	}

	//Writes resume to created destination
	n2, err := dst.Write(avatar)
	defer dst.Close()
	if err != nil {
		log.Printf("error in writing to dst %d", n2)
	}

	// // Handle errors.
	// if c.Validation.HasErrors() {
	// 	c.Validation.Keep()
	// 	c.FlashParams()
	// 	return c.Redirect(App.Upload)
	// }
	ListBuckets(name, path)

	err = os.Remove(path)
	if err != nil	{
		log.Printf("error removing temp file")
	}

	return c.RenderJson(FileInfo{
		ContentType: c.Params.Files["avatar"][0].Header.Get("Content-Type"),
		Filename:    c.Params.Files["avatar"][0].Filename,
		//RealFormat:  format,
		//Resolution:  fmt.Sprintf("%dx%d", conf.Width, conf.Height),
		//Size:        len(avatar),
		Status:      "Successfully uploaded",
	})
}

func ListBuckets(key, path string) {
	file, err := os.Open(path)
    if err != nil {
        log.Fatal("Failed to open file", err)
    }

    // Not required, but you could zip the file before uploading it
    // using io.Pipe read/writer to stream gzip'd file contents.
    reader, writer := io.Pipe()
    go func() {
        gw := gzip.NewWriter(writer)
        io.Copy(gw, file)

        file.Close()
        gw.Close()
        writer.Close()
    }()
    uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-west-2")}))
    result, err := uploader.Upload(&s3manager.UploadInput{
        Body:   reader,
        Bucket: aws.String("trurecruit"),
        Key:    aws.String(key),
    })
    if err != nil {
        log.Fatalln("Failed to upload", err)
    }

    log.Println("Successfully uploaded to", result.Location)
}


// func (c *Single) ListBuckets() revel.Result {
// 	bucket := "trurecruit"
// 	key := "test.txt"
//
// 	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))
// 	// result, err := svc.CreateBucket(&s3.CreateBucketInput{
// 	//     Bucket: &bucket,
// 	// })
// 	// if err != nil {
// 	//     log.Println("Failed to create bucket", err)
// 	//     return nil
// 	// }
// 	//
// 	// if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {
// 	//     log.Printf("Failed to wait for bucket to exist %s, %s\n", bucket, err)
// 	//     return nil
// 	// }
//
// 	uploadResult, err := svc.PutObject(&s3.PutObjectInput{
// 	    Body:   strings.NewReader("Hello World!"),
// 	    Bucket: &bucket,
// 	    Key:    &key,
// 	})
// 	if err != nil {
// 	    log.Printf("Failed to upload data to %s/%s, %s\n", bucket, key, err)
// 	    return nil
// 	}
// 	// if result != nil {
// 	// 	log.Printf("Successfully created bucket %s and uploaded data with key %s\n", bucket, key)
// 	// 	return nil
// 	// }
// 	if uploadResult != nil {
// 		log.Printf("wa")
// 	}
// 	return c.Redirect(App.Upload)
// }
