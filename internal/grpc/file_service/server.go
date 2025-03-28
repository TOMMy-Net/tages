package fileservice

import (
	"sync"
	"time"

	FS "github.com/TOMMy-Net/tages/pkg/protobuf/file_service/proto"
	"google.golang.org/grpc"
)

type FileMeta struct {
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileServer struct {
	FS.UnimplementedFileServiceServer
	filesDir        string
	serviceFilesDir string
	metaData        map[string]FileMeta
	metaDataMutex   sync.RWMutex
	limiter         *Limiter
}

func (fs *FileServer) SetFilesDir(dir string) {
	fs.filesDir = dir
}

func (fs *FileServer) GetFilesDir() string {
	return fs.filesDir
}

func (fs *FileServer) SetServiceFilesDir(dir string) {
	fs.serviceFilesDir = dir
}

func (fs *FileServer) GetServiceFilesDir() string {
	return fs.serviceFilesDir
}
func NewFileServer(filesDir string) *FileServer {
	fs := &FileServer{
		filesDir:        filesDir,
		serviceFilesDir: "./",
		metaData:        make(map[string]FileMeta, 100),
		limiter:         NewLimiter(),
	}
	return fs
}

func Register(gs *grpc.Server, fs *FileServer) error{
	err := fs.LoadMetaData()
	if err != nil {
		return err
	}
	FS.RegisterFileServiceServer(gs, fs)
	return nil
}
