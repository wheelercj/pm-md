package main

type Request struct {
	Method string `json:"method"`
	Header []any  `json:"header"`
	Body   struct {
		Mode    string `json:"mode"`
		Raw     string `json:"raw"`
		Options struct {
			Raw struct {
				Language string `json:"language"`
			} `json:"raw"`
		} `json:"options"`
	} `json:"body"`
	Url struct {
		Raw  string   `json:"raw"`
		Host []string `json:"host"`
		Path []string `json:"path"`
	} `json:"url"`
}

type Routes []struct {
	Name                    string `json:"name"`
	ProtocolProfileBehavior struct {
		DisableBodyPruning bool `json:"disableBodyPruning"`
	} `json:"protocolProfileBehavior"`
	Request   Request `json:"request"`
	Responses []struct {
		Name            string  `json:"name"`
		OriginalRequest Request `json:"originalRequest"`
		Status          string  `json:"status"`
		Code            int     `json:"code"`
		Language        string  `json:"_postman_previewlanguage"`
		Headers         []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"header"`
		Cookies []any  `json:"cookie"`
		Body    string `json:"body"`
	} `json:"response"`
}

type Collection struct {
	Info struct {
		PostmanId  string `json:"_postman_id"`
		Name       string `json:"name"`
		Schema     string `json:"schema"`
		ExporterId string `json:"_exporter_id"`
	} `json:"info"`
	Routes Routes `json:"item"`
	Events []struct {
		Listen string `json:"listen"`
		Script struct {
			Type string `json:"type"`
			Exec []string
		} `json:"script"`
	} `json:"event"`
	Variables []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"variable"`
}
