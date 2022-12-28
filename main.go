package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/mdp/qrterminal/v3"
)

func main() {
	port := randomport()
	qrcode(port)
	httplisten(port)
}

// ///////////////////////////////////////////////////////////
// http listener funcsiya
func httplisten(a string) {
	app := iris.New()

	app.Get("/", indexPage)
	app.Get("/download/{url:path}", downloadfile)
	app.Post("/upload", uploadfile)
	app.Listen(a)
}

func uploadfile(ctx iris.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.Writef(string("error"))
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./%s", fileHeader.Filename))
	if err != nil {
		ctx.Writef("error for upload file")
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		ctx.Writef("error for upload file")
		return
	}
	ctx.Redirect("/")
}

// download file function
func downloadfile(ctx iris.Context) {
	filename := ctx.Params().Get("url")

	// ctx.SendFile(filename, "./")  2-usuli fayl junatish uchun

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		ctx.JSON(iris.Map{"message": "error for open file"})
		return
	}

	attfile := fmt.Sprintf("attachment; filename=%s", filename)
	ctx.StatusCode(200)
	ctx.Header("Content-Disposition", attfile)
	ctx.Header("Context-Type", "application/octet-stream")
	ctx.Write(fileBytes)
	return
}

// index page code
func indexPage(ctx iris.Context) {
	index := `<!DOCTYPE html>
	<html>
	<head>
	<style>
	#customers {
	  font-family: Arial, Helvetica, sans-serif;
	  border-collapse: collapse;
	  width: 100%;
	}
	
	#customers td, #customers th {
	  border: 1px solid #ddd;
	  padding: 8px;
	}
	
	#customers tr:nth-child(even){background-color: #f2f2f2;}
	
	#customers tr:hover {background-color: #ddd;}
	
	#customers th {
	  padding-top: 12px;
	  padding-bottom: 12px;
	  text-align: left;
	  background-color: #04AA6D;
	  color: white;
	}
	</style>
	</head>
	<body>
	<form method='POST' action='/upload' enctype='multipart/form-data'>
	<h4>Select file for upload:</h4>
	<input type='file' name='file'><input type='submit' value='Upload'></input>
	</form><br>
	`
	bulim := `
	<table id="customers">
  <tr>
    <th>File Name</th>
    <th>size</th>
    <th>Option</th>
  </tr>
	`
	tugash := `
	</table>
	</body>
	</html>`
	ctx.HTML(index + bulim + fileslist() + tugash)
}

// ///////////////////////////////////////////////////////////
// files list function
func fileslist() string {
	var lists string
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.Size() == 4096 {
			continue
		}
		lists += fmt.Sprintf("<tr><td>%s</td><td>%v byte</td><td><a href='/download/%s'>Download</a></td></tr>", f.Name(), f.Size(), f.Name())
	}
	return lists
}

// qr kodni ekranga chiqaradigan funcsiya
func qrcode(a string) {
	myip := GetLocalIP()
	connect := fmt.Sprintf("http://%s%s", myip, a)
	fmt.Println("Scan to qr code for connect!")
	fmt.Println(connect)
	qrterminal.Generate(connect, qrterminal.L, os.Stdout)
}

// get my ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// random port tayorlab beradigan funcsiya
func randomport() string {
	randomport, _ := rand.Int(rand.Reader, big.NewInt(int64(50000)))
	port := fmt.Sprintf(":%v", randomport.Uint64())
	return port
}
