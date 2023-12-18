package utils

const HASH_VALUE = 3581

func ToDJBHash(str string) int {

	hash := HASH_VALUE

	for _, ch := range str {
		hash = (hash << 5) + hash + int(ch)
	}

	return hash
}
