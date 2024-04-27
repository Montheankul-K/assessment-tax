package taxHandlers

import "github.com/labstack/echo/v4"

type IResponse interface {
	ResponseSuccess(statusCode int, data interface{}) error
	ResponseError(statusCode int, errMsg string) error
}

type Response struct {
	Context echo.Context
}

type Error struct {
	Message string `json:"message"`
}

type TaxResponse struct {
	Tax      float64            `json:"tax"`
	TaxLevel []TaxLevelResponse `json:"taxLevel"`
	TotalTax float64            `json:"totalTax"`
}

type TaxResponseWithRefund struct {
	TaxResponse
	TaxRefund float64 `json:"taxRefund"`
}

func NewResponse(c echo.Context) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) ResponseSuccess(statusCode int, data interface{}) error {
	return r.Context.JSON(statusCode, data)
}

func (r *Response) ResponseError(statusCode int, errMessage string) error {
	return r.Context.JSON(statusCode, &Error{
		Message: errMessage,
	})
}
