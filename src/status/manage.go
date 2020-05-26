package status

import (
	"tea/src/manage"
)

func GetClientCount() int {
	return manage.LocalManage.GetClientCount()
}
