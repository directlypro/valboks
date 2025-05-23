package dropbox

import (
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
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

	fileInfos = append(fileInfos, c.processEntries(result.Entries))

	for result.HasMore {
		continueArg := files.NewListFolderContinueArg(result.Cursor)
		result, err = c.filesClient.ListFolderContinue(continueArg)
		if err != nil {
			return nil, fmt.Errorf("failed to continue listing folder: %w", err)
		}

		fileInfos = append(fileInfos, c.processEntries(result.Entries))
	}

	return fileInfos, nil
}

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