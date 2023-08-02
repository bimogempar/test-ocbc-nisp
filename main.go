package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DataFile struct {
	ID         string
	FileName   string
	UploadTime time.Time
	Status     string
}

var InMemoryDB map[string]*DataFile

func init() {
	InMemoryDB = make(map[string]*DataFile)
}

func uploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileID := fmt.Sprintf("%d", time.Now().UnixNano())

	_ = saveCSVFile(file, fileID)

	dataFile := &DataFile{
		ID:         fileID,
		FileName:   "your_file_name.csv",
		UploadTime: time.Now(),
		Status:     "Processing",
	}

	InMemoryDB[fileID] = dataFile

	go importDataFromFile(dataFile)

	w.Write([]byte(fileID))
}

func importDataFromFile(dataFile *DataFile) {
	time.Sleep(5 * time.Second)

	dataFile.Status = "Completed"
}

func saveCSVFile(file io.Reader, fileID string) error {
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(record)
	}
	return nil
}

func main() {
	http.HandleFunc("/upload", uploadCSVHandler)
	http.ListenAndServe(":7000", nil)
	fmt.Println("running bro")
}
