// response external api
package extresponse

type TranslateExtRes struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Original    string `json:"original"`
}
