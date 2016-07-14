package controllers

import (
	//"bytes"
	"fmt"
	//"image"
	//_ "image/jpeg"
	//_ "image/png"

	"compress/gzip"
	"io"
  "log"
  "os"
	//"net/http"
	//"bufio"

	"path/filepath"

	"math/rand"
	"strconv"
	"html/template"

	//"strings"
	// "github.com/revel/samples/upload/app/routes"

	"github.com/revel/revel"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var global_code string

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

type Single struct {
		*revel.Controller
}

func (c Single) Upload() revel.Result {
	return c.Render()
}

func (c Single) Download() revel.Result  {
	return c.Render()
}

func (c Single) HandleUpload(avatar []byte) revel.Result {
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
	// var path string = "src/myapp/upload/" + name

	//Check extension of file
	var ext string = filepath.Ext(name)
	if (ext != ".pdf") {
		log.Printf("File must have a .pdf extension.")
		return c.Redirect(Single.Upload)
	}
	log.Printf("extension is " + ext)

	var r [8]int
	for i := 0; i < 8; i++ {
		r[i] = rand.Intn(9)
	}

	fmt.Println("%v", r[0])
	var strr string = ""
	for j := 0; j < 8; j++ {
		strr += strconv.Itoa(r[j])
	}
	fmt.Println("code: " + strr)

	global_code = strr

	var newName string = strr + name
	var newPath string = "src/myapp/upload/" + newName

	dst, err := os.Create("src/myapp/upload/" + newName)
	defer dst.Close()
	if err != nil {
		log.Printf("error in creating dst")
	}

	log.Printf("new file name is " + newName)
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
	UploadFile(newName, newPath)

	err = os.Remove(newPath)
	if err != nil	{
		log.Printf("error removing temp file")
	}

	//Returns json info of uploaded file, needs to be changed
	return c.RenderJson(FileInfo{
		ContentType: c.Params.Files["avatar"][0].Header.Get("Content-Type"),
		Filename:    c.Params.Files["avatar"][0].Filename,
		//RealFormat:  format,
		//Resolution:  fmt.Sprintf("%dx%d", conf.Width, conf.Height),
		//Size:        len(avatar),
		Status:      "Successfully uploaded",
	})
}

func UploadFile(key, path string) {
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

		//Creates new uploader and uploads file passed path location
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

func returnCode() {
	tmpl, err := template.New("").Parse("{{.returnCode}}")
	if err != nil {
		log.Fatalf("Parse %v", err)
	}
	tmpl.Execute(os.Stdout, global_code)
}
