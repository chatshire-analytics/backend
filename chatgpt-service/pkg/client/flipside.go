package client

type CreateFlipsideQueryRequest struct {
	Sql        string `json:"sql"`
	TtlMinutes int    `json:"ttl_minutes"`
	Cache      bool   `json:"cache"`
	Params     struct {
		AdditionalProp1 string `json:"additionalProp1"`
		AdditionalProp2 string `json:"additionalProp2"`
		AdditionalProp3 string `json:"additionalProp3"`
	} `json:"params"`
}

type CreateFlipsideQuerySuccessResponse struct {
	Token  string `json:"token"`
	Cached bool   `json:"cached"`
}

type CreateFlipsideQueryErrorResponse struct {
	Errors struct {
		AdditionalProp1 struct {
		} `json:"additionalProp1"`
	} `json:"errors"`
}

type GetFlipsideQueryResultSuccessResponse struct {
	Token  string `json:"token"`
	Cached bool   `json:"cached"`
}

type GetFlipsideQueryResultErrorResponse struct {
	Errors struct {
		AdditionalProp1 struct {
		} `json:"additionalProp1"`
	} `json:"errors"`
}
