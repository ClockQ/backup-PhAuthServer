package PhResource

// The Response struct implements api2go.Responder
type Response struct {
	Res        interface{}
	Code       int
	QueryRes   string
	TotalPage  int
	TotalCount int
}

// Metadata returns additional meta data
func (r Response) Metadata() map[string]interface{} {
	var meta = map[string]interface{}{
		"author":      "Pharbers Inc. AlfredYang",
		"license":     "pharbers.ntm",
		"license-url": "http://www.pharbers.com",
	}

	if r.QueryRes != "" && r.TotalPage >= 0 {
		meta["query-res"] = r.QueryRes
		meta["total-page"] = r.TotalPage
		meta["total-count"] = r.TotalCount
	}

	return meta
}

// Result returns the actual payload
func (r Response) Result() interface{} {
	return r.Res
}

// StatusCode sets the return status code
func (r Response) StatusCode() int {
	return r.Code
}
