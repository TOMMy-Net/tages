package fileservice

import (
	"io"
	"os"
	"path/filepath"
	"time"

	FS "github.com/TOMMy-Net/tages/pkg/protobuf/file_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (fs *FileServer) Upload(stream FS.FileService_UploadServer) error {
	ctx := stream.Context()
	clientIP := getClientIP(ctx)
	limits := fs.limiter.GetLimits(clientIP)

	if !limits.DownloadAndUploadConn.Add() {
		return status.Error(codes.PermissionDenied, "many connections")
	}
	defer limits.DownloadAndUploadConn.Release()

	var filename string
	var fileData []byte

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filename == "" {
			filename = filepath.Base(req.GetFilename())
			if filename == "" {
				return status.Error(codes.InvalidArgument, "invalid filename")
			}
		}

		fileData = append(fileData, req.Chunk...)
	}

	filePath := filepath.Join(fs.filesDir, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	fs.metaDataMutex.Lock()
	defer fs.metaDataMutex.Unlock()

	now := time.Now()
	if meta, exists := fs.metaData[filename]; exists {
		meta.UpdatedAt = now
		fs.metaData[filename] = meta
	} else {
		fs.metaData[filename] = FileMeta{
			CreatedAt: now,
			UpdatedAt: now,
			FileName:  filename,
		}
	}

	if err := fs.SaveMetaData(); err != nil {
		return status.Error(codes.Internal, "metadata save failed")
	}

	return stream.SendAndClose(&FS.UploadResponse{
		Filename: filename,
		Size:     uint64(len(fileData)),
	})
}
