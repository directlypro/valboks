package dropbox

import (
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"io"
	"os"
	"strings"
)

type Client struct {
	filesClient files.Client
}

type FileInfo struct {
	Name string
	Path string
	IsFolder bool
	Size uint64
}

func NewClient(accessToken string) *Client {
	config := dropbox.Config{
		Token: accessToken,
		LogLevel: dropbox.LogOff,
	}

	return &Client{
		filesClient: files.New(config),
	}
}

func (c *Client) ListFolder(path string) ([]FileInfo, error) {

	path = normalizePath(path)

	listArg := files.NewListFolderArg(path)
	result, err := c.filesClient.ListFolder(listArg)
	if err != nil {
		return nil, fmt.Errorf("failed to list folder '%s': %w", path, err)
	}

	var fileInfos []FileInfo

	fileInfos = append(fileInfos, c.processEntries(result.Entries)...)

	for result.HasMore {
		continueArg := files.NewListFolderContinueArg(result.Cursor)
		result, err = c.filesClient.ListFolderContinue(continueArg)
		if err != nil {
			return nil, fmt.Errorf("failed to continue listing folder: %w", err)
		}

		fileInfos = append(fileInfos, c.processEntries(result.Entries)...)
	}

	return fileInfos, nil
}

//processEntries converts Dropbox API entries to FileInfo structs
func (c *Client) processEntries(entries []files.IsMetadata) []FileInfo {
	var fileInfos []FileInfo

	for _, entry := range entries {
		switch e := entry.(type) {
		case *files.FolderMetadata:
			fileInfos = append(fileInfos, FileInfo{
			Name: e.Name,
			Path: e.PathLower,
			IsFolder: true,
			Size: 0,
			})
		case *files.FileMetadata:
			fileInfos = append(fileInfos, FileInfo{
			Name: e.Name,
			Path: e.PathLower,
			IsFolder: false,
			Size: e.Size,
			})
		}
	}

	return fileInfos
}

//downloading file from dropbox to local path
func (c *Client) DownloadFile(dropboxPath, localPath string) error {

	dropboxPath = normalizePath(dropboxPath)

	downloadArg := files.NewDownloadArg(dropboxPath)
	_, content, err := c.filesClient.Download(downloadArg)
	if err != nil {
		return fmt.Errorf("failed to download file '%s': %w", dropboxPath, err)
	}
	defer content.Close()

	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file '%s': %w", localPath, err)
	}
	defer  outFile.Close()

	//Copy content to file
	_, err = io.Copy(outFile, content)
	if err != nil {
		return fmt.Errorf("failed to write file content: '%w'", err)
	}

	return nil
}


func (c *Client) UploadFile(localPath, dropboxPath string, overwrite bool) error {

	dropboxPath = normalizePath(dropboxPath)

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file '%s': %w", localPath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get the file info: %w", err)
	}

	commitInfo := files.NewCommitInfo(dropboxPath)
	if overwrite {
		commitInfo.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}}
	}

	uploadArg := files.NewUploadArg(commitInfo.Path)

	if fileInfo.Size() < 150*1024*1024 {
		_, err = c.filesClient.Upload(uploadArg, file)
		if err != nil {
			return fmt.Errorf("failed to upload file '%s': %w", localPath, err)
		}
	} else {
		return fmt.Errorf("files lager than 150MB are not supported in this version")
	}

	return nil
}

func (c * Client) DeletePath(path string) error {

	path = normalizePath(path)

	deleteArg := files.NewDeleteArg(path)
	_, err := c.filesClient.DeleteV2(deleteArg)
	if err != nil {
		return fmt.Errorf("failed to delete '%s': %w", path, err)
	}

	return nil
}

func (c *Client) CreateFolder(path string) error {

	path = normalizePath(path)

	createArg := files.NewCreateFolderArg(path)
	_, err := c.filesClient.CreateFolderV2(createArg)
	if err != nil {
		return fmt.Errorf("failed to create a folder '%s': %w", path, err)
	}

	return nil
}

func (c *Client) GetFileInfo(path string) (*FileInfo, error) {

	path = normalizePath(path)

	getMetadataArg := files.NewGetMetadataArg(path)
	metadata, err := c.filesClient.GetMetadata(getMetadataArg)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for '%s': %w", path, err)
	}

	switch m := metadata.(type) {
	case *files.FolderMetadata:
		return &FileInfo{
			m.Name,
			m.PathLower,
			true,
			0,
		}, nil
	case *files.FileMetadata:
		return &FileInfo{
			m.Name,
			m.PathLower,
			false,
			m.Size,
		}, nil
	default:
		return nil, fmt.Errorf("unknown metadata type for '%s'", path)
	}
}

func (c *Client) TestConnection() error {
	_, err := c.GetFileInfo("/")
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	return nil
}

func normalizePath(path string) string {
	if path == "" || path == "." {
		return ""
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if path == "/" {
		return ""
	}

	return path
}