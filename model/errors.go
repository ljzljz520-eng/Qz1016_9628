package model

import "errors"

var errInvalidPage = errors.New("page bounds must not be negative")
var errInvalidPriority = errors.New("unknown priority")
var errInvalidStatus = errors.New("unknown status filter")

func IsValidationError(err error) bool {
	return err != nil
}

func StatusLabel(status string) string {
	switch status {
	case "new":
		return "待校验"
	case "validated":
		return "已校验"
	case "processing":
		return "处理中"
	case "blocked":
		return "已阻塞"
	case "resolved":
		return "已解决"
	case "closed":
		return "已关闭"
	case "archived":
		return "已归档"
	case "cancelled":
		return "已取消"
	default:
		return "未知"
	}
}
