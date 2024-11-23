package model

type Status struct {
	Code    int    `json:"code"`
	Message string `message:"message"`
}

type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

type PagedResponse struct {
	Status Status        `json:"status"`
	Data   []interface{} `json:"data"`
	Paging Paging        `json:"paging"`
}

type Paging struct {
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
	TotalRows   int `json:"totalRows"`
	TotalPages  int `json:"totalPages"`
}