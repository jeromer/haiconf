package stringutils

func RemoveDuplicates(strList []string) []string {
	l := len(strList)
	alreadyInserted := make(map[string]bool, l)
	noDups := make([]string, l)

	for i, str := range strList {
		_, present := alreadyInserted[str]
		if present {
			continue
		}

		noDups[i] = str
		alreadyInserted[str] = true
	}

	return noDups[:len(alreadyInserted)]
}
