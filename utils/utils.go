package utils

func InArray(need string, needArr []string) bool {
	for _,v := range needArr{
		if need == v{
			return true
		}
	}
	return false
}