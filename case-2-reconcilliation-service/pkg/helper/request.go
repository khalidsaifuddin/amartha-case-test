package helper

func SetupListParameter(page, pageSize int32) (int32, int32, int32, int32) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	limit := pageSize
	start := (page - 1) * pageSize
	var totalData int32 = 0
	var totalPage int32 = 0

	return limit, start, totalPage, totalData
}
