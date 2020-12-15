package a

type Resource struct {
        ID   int64    `json:"id"    xml:"id"`
        Data []string `json:"data,omitempty"           xml:"data"`
}
