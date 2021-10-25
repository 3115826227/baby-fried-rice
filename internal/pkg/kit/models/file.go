package models

import "encoding/json"

type FileValueType int

const (
	FileId      FileValueType = 0
	FileDownUrl               = 1
)

type DeleteFileMessageQueueInfo struct {
	FileValueType FileValueType `json:"file_value_type"`
	FileValue     string        `json:"file_value"`
}

func (info DeleteFileMessageQueueInfo) ToString() string {
	data, _ := json.Marshal(info)
	return string(data)
}
