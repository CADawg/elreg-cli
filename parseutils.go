package main

func GetNLines(command string) int {
	length, ok := lengthMap[command]

	if !ok {
		return 1
	}

	return length
}
