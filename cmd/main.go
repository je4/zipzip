package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

func createZip() string {
	// step 1: create a zip file
	zipFP, err := os.CreateTemp("", "example.zip")
	if err != nil {
		panic(err)
	}
	defer zipFP.Close()
	fmt.Printf("%s created\n", zipFP.Name())

	// step 2: create a zip writer
	zipWriter := zip.NewWriter(zipFP)
	defer zipWriter.Close()

	// step 3: add an uncompressed file to the zip
	fileWriter, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:    "example.zip",
		Comment: "zip inside a zip",
		Method:  zip.Store,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s added to %s\n", "example.zip", zipFP.Name())

	// step 4: create a zip inside the zip
	zipzipFP := zip.NewWriter(fileWriter)
	defer zipzipFP.Close()

	// step 5: add a file to the zip inside the zip
	zipFileWriter, err := zipzipFP.Create("example.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s added to %s\n", "example.txt", "example.zip")

	content := "[This is the content of a file inside a zip file]"
	_, err = zipFileWriter.Write([]byte(content))
	if err != nil {
		panic(err)
	}
	fmt.Printf("'%s' written to %s\n", content, "example.txt")

	return zipFP.Name()
}

func printZipZipContent(zipName string) {

	// step 1: open the zip file
	zipFP, err := os.Open(zipName)
	if err != nil {
		panic(err)
	}
	defer zipFP.Close()
	fmt.Printf("%s opened\n", zipFP.Name())

	// step 2: create a zip reader
	fi, err := zipFP.Stat()
	if err != nil {
		panic(err)
	}
	zipReader, err := zip.NewReader(zipFP, fi.Size())
	if err != nil {
		panic(err)
	}

	// step 3: read the zip file inside the zip
	file := zipReader.File[0]
	zipzipFP, err := file.OpenRaw()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s opened raw inside %s\n", file.Name, zipFP.Name())

	zipzipFPReaderAt, ok := zipzipFP.(io.ReaderAt)
	if !ok {
		panic("zipzipFP does not implement io.ReaderAt")
	}
	zipzipReader, err := zip.NewReader(zipzipFPReaderAt, int64(file.CompressedSize64))
	if err != nil {
		panic(err)
	}

	file2 := zipzipReader.File[0]
	zipzipContentFP, err := file2.Open()
	if err != nil {
		panic(err)
	}
	defer zipzipContentFP.Close()
	fmt.Printf("%s opened inside %s\n", file2.Name, file.Name)

	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, zipzipContentFP)

	fmt.Printf("content of %s: '%s'\n", file2.Name, buffer.String())
}

func main() {
	zipName := createZip()
	defer os.Remove(zipName)

	printZipZipContent(zipName)

}
