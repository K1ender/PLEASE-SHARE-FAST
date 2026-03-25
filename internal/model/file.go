package model

import "time"

type File struct {
	Data []byte
	FileMetadata
}

type FileMetadata struct {
	Filename  string
	CreatedAt time.Time
}
