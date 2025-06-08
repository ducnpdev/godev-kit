// request external api
package extrequest

type TranslateExtReq struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Original    string `json:"original"`
}
