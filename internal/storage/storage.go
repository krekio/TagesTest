package storage

import (
	"errors"
	pb "grpc-file-service/proto"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	storagePath string
}

// NewFileStorage creates a new FileStorage
//
// It creates a directory if it's not already exists and returns
// a pointer to a new FileStorage instance.
func NewFileStorage(path string) (*FileStorage, error) {
	// Check if a directory already exists
	if _, err := os.Stat(path); err != nil {
		// Create a directory
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &FileStorage{storagePath: path}, nil
}

// Upload Uploading a file
//
// This function is used to upload a file to the storage. It receives a stream of
// UploadFileRequest and writes the data to the file. Once the stream is closed,
// the function returns a UploadFileResponse with a success message.
func (s *FileStorage) Upload(stream pb.FileService_UploadFileServer) error {
	// Open a file for writing
	var filename string
	var file *os.File

	// Read a file from the stream
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// If the stream is closed, break the loop
			break
		}
		if err != nil {
			// If an error occurs, return the error
			return err
		}

		// If the file is not opened yet, open it
		if file == nil {
			// Create a full path to the file
			filename = filepath.Join(s.storagePath, req.GetFilename())
			// Create the file
			file, err = os.Create(filename)
			if err != nil {
				// If an error occurs, return the error
				return err
			}
			// Close the file when the function is finished
			defer file.Close()
		}

		// Write the data to the file
		if _, err := file.Write(req.GetData()); err != nil {
			// If an error occurs, return the error
			return err
		}
	}

	// Send a response to the client
	return stream.SendAndClose(&pb.UploadFileResponse{Message: "Файл успешно загружен"})
}

// List Viewing a list of files
//
// This function is used to view a list of files. It reads a list of files
// from the storage directory and returns a ListFilesResponse with a list
// of FileInfo.
func (s *FileStorage) List() (*pb.ListFilesResponse, error) {
	files, err := os.ReadDir(s.storagePath)
	if err != nil {
		return nil, err
	}

	var fileInfos []*pb.FileInfo
	for _, f := range files {
		if f.IsDir() {
			// Skip directories
			continue
		}

		filePath := filepath.Join(s.storagePath, f.Name())
		fileStat, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, &pb.FileInfo{
			Filename:  f.Name(),
			CreatedAt: fileStat.ModTime().Format(time.RFC3339),
			UpdatedAt: fileStat.ModTime().Format(time.RFC3339),
		})
	}

	return &pb.ListFilesResponse{Files: fileInfos}, nil
}

// Download Downloading a file
//
// This function is used to download a file from the storage. It opens the file,
// reads its content and sends it to the client in chunks of 1024 bytes.
func (s *FileStorage) Download(filename string, stream pb.FileService_DownloadFileServer) error {
	// Create a full path to the file
	filePath := filepath.Join(s.storagePath, filename)
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// If an error occurs, return an error with a message
		return errors.New("файл не найден")
	}
	// Close the file when the function is finished
	defer file.Close()

	// Create a buffer for reading the file
	buf := make([]byte, 1024)
	for {
		// Read from the file
		n, err := file.Read(buf)
		// If the end of the file is reached, break the loop
		if err == io.EOF {
			break
		}
		// If an error occurs, return the error
		if err != nil {
			return err
		}

		// Send the data to the client
		if err := stream.Send(&pb.DownloadFileResponse{Data: buf[:n]}); err != nil {
			// If an error occurs, return the error
			return err
		}
	}
	return nil
}
