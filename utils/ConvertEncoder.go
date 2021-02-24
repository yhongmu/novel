package utils

import "golang.org/x/text/encoding/simplifiedchinese"

func ConvertGBKToUTF8(str string) (string, error) {
	return simplifiedchinese.GBK.NewDecoder().String(str)
}
