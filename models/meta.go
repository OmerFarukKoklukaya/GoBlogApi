package models

type ViewData struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	LastPage int `json:"lastPage"`
}

func (Pagination) Paginate(page int, dataLength int, limit int) Pagination {
	var lastPage = 0
	if dataLength%limit == 0 {
		lastPage = dataLength / limit
	} else {
		lastPage = dataLength/limit + 1
	}
	var paging = Pagination{
		Page:  page,
		Limit: limit,

		LastPage: lastPage,
	}

	return paging
}
