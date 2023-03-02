package client

import "time"

type CreateFlipsideQueryRequest struct {
	Id         string `json:"id"`
	Sql        string `json:"sql"`
	TtlMinutes int    `json:"ttl_minutes"`
	Cache      bool   `json:"cache"`
	Params     struct {
		AdditionalProp1 string `json:"additionalProp1"`
		AdditionalProp2 string `json:"additionalProp2"`
		AdditionalProp3 string `json:"additionalProp3"`
	} `json:"params"`
}

func NewCreateFlipsideQueryResult(id string, sql string) *CreateFlipsideQueryRequest {
	return &CreateFlipsideQueryRequest{
		Id:         id,
		Sql:        sql,
		TtlMinutes: 15,
		Cache:      true,
	}
}

type CreateFlipsideQuerySuccessResponse struct {
	Token  string `json:"token"`
	Cached bool   `json:"cached"`
}

type CommonFlipsideQueryErrorResponse struct {
	Errors struct {
		AdditionalProp1 struct {
		} `json:"additionalProp1"`
	} `json:"errors"`
}

type GetFlipsideQueryResultSuccessResponse struct {
	Results      [][]interface{} `json:"results"`
	ColumnLabels []string        `json:"columnLabels"`
	ColumnTypes  []string        `json:"columnTypes"`
	StartedAt    time.Time       `json:"startedAt"`
	EndedAt      time.Time       `json:"endedAt"`
	PageNumber   int             `json:"pageNumber"`
	PageSize     int             `json:"pageSize"`
	Status       string          `json:"status"`
}

type GetFlipsideQueryResultRequest struct {
	Token string `json:"token"`
}

func NewGetFlipsideQueryResultRequest(token string) *GetFlipsideQueryResultRequest {
	return &GetFlipsideQueryResultRequest{
		Token: token,
	}
}
