package utility

func ParseBranchIDFromContext(bid interface{}) uint64 {
	branchID, ok := bid.(uint64)
	if !ok {
		return 0
	}

	return branchID
}