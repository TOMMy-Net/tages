package fileservice

import (
	"io"
	"os"
	"path/filepath"

	FS "github.com/TOMMy-Net/tages/pkg/protobuf/file_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (fs *FileServer) Download(req *FS.DownloadRequest, stream FS.FileService_DownloadServer) error {
	ctx := stream.Context()
	clientIP := getClientIP(ctx)
	limits := fs.limiter.GetLimits(clientIP)
	if !limits.DownloadAndUploadConn.Add() {
		return status.Error(codes.PermissionDenied, "many connections")
	}
	defer limits.DownloadAndUploadConn.Release()

	filename := filepath.Base(req.Filename)
	filePath := filepath.Join(fs.filesDir, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return status.Error(codes.NotFound, "file not found")
	}
	defer file.Close()

	buffer := make([]byte, 64*1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&FS.DownloadResponse{Chunk: buffer[:n]}); err != nil {
			return err
		}
	}

	return nil
}
