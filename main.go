package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Downloader interface{
	Download(url string)(r io.Reader, err error)
}
var webID = 0

type web struct{
	id int
}
func (v web)NewDownloader(){
	webID++
}

func (v *web)Download(url string)(io.Reader,error){

	// get the data
	resp, err := http.Get(url)
	if err != nil{
		return resp.Body,err
	}

	//defer resp.Body.Close()

	// Create the file
	filepath := "f"
	out, err := os.Create(fmt.Sprintf(filepath+"%d",webID))

	if err != nil{
		fmt.Println("Error creating file")
		return resp.Body,err
	}

	//defer out.Close()

	// Write the body to file
	_,err = io.Copy(out,resp.Body)
	return resp.Body,err
}

type Archiver interface{
	Archive(names []string, readers ...io.Reader)(outR io.Reader,err error)
}

type zipp struct{

}

func(z *zipp) Archive(names []string, readers ...io.Reader)(outR io.Reader,err error){
	result, err := os.Create("result.zip")
	
	if err != nil {
        panic(err)
    }
	defer result.Close()
    zipW := zip.NewWriter(result)
	for i, f := range names{
		w1, err := zipW.Create(f)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(w1, readers[i])
		if err != nil {
			panic(err)
		}
	}
	zipW.Close()
	return result,nil

	// result, err := os.Create("result.zip")
	// if err != nil {
	// 	panic(err)
	// }
	// defer result.Close()
	// zipW := zip.NewWriter(result)
	// for _, f := range names {
	// 	file, err := os.Open(f)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()
	// 	w1, err := zipW.Create(f)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	_, err = io.Copy(w1, file)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// zipW.Close()
	// return result, nil
}



func main(){
	fileUrl := "https://www.jungleerummy.com/blog/assets/front/img/blog/bsdy4m1efq9i8md9m6gzelymija995khb9kwinei8kjhewbsjjmain.jpg"
	webFile := web{}
	webFile.NewDownloader()
	webFile.id = webID
	r1, err := webFile.Download(fileUrl)
	if err!= nil{
		fmt.Printf("got error while downloading file with id %d\n",webFile.id)
		log.Fatal(err)
	}
	fileUrl = "https://humornama.com/wp-content/uploads/2020/10/Miracle-Miracle-meme-template-of-welcome-movie-1024x576.jpg"
	webFile = web{}
	webFile.NewDownloader()
	webFile.id = webID
	r2, err := webFile.Download(fileUrl)
	if err != nil{
		fmt.Printf("got error while downloading file with id %d\n",webFile.id)
	}
	zipper := zipp{}
	_, err = zipper.Archive([]string{"f1","f2"}, r1, r2)
	if err != nil{
		panic(err)
	}
	
}

