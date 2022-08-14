package uploader

const (
	// AllSuccess 全部成功
	AllSuccess = iota
	// PartSuccess 部分成功
	PartSuccess
	// Fail 保存操作失败
	Fail

	fileDateLayout   = "2006-01"
	savePath         = "/home/d/public/dav/%s"
	savePathWhitDate = "/home/d/public/dav/%s/%s"
	save2Path        = "%s/%s"
	savePathPrefix   = "/home/d"
	visURL           = "https://d.evlic.cn/public/dav/%s/%s"
)