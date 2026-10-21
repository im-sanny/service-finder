package repository

func CalculateOffset(page, limit int) (int, int) {
	if page < 0 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return (page - 1) * limit, limit
}
