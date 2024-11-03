package entity

//go:generate msgp -tests=false -io=false -unexported=true

type BucketMeta struct {
	Name string `json:"name,omitempty"` // 分块的文件名称
	Hash string `json:"hash,omitempty"` // 分块文件的哈希值
	Size int64  `json:"size,omitempty"` // 分块文件的大小
}

type FileMeta struct {
	Name    string       `json:"name,omitempty"`    // 文件名称
	Buckets []BucketMeta `json:"buckets,omitempty"` // 文件分块信息
	Hash    string       `json:"hash,omitempty"`    // 文件的哈希值
}
