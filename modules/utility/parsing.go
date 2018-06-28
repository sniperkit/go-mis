package utility

import "strconv"

func ParseBranchIDFromContext(bid interface{}) uint64 {
	str, ok := bid.(string)
	if !ok {
		return 0
	}

	branchID, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}

	return branchID
}