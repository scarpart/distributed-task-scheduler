package lbutils

import "net/http"

func HeaderToMap(header http.Header) map[string]string {
	result := make(map[string]string)
	for key, val := range header {
		if len(val) > 0 {
			result[key] = val[0]
		}
	}
	return result
}
