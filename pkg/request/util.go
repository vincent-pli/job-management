package request

import (
	"fmt"
)

func GetJobKeyByReq(req *Request) string {
	return fmt.Sprintf("%s/%s", req.Namespace, req.JobName)
}