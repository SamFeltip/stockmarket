package websockets

type RequestHeaders struct {
	HXRequest     string `json:"HX-Request"`
	HXTrigger     string `json:"HX-Trigger"`
	HXTriggerName string `json:"HX-Trigger-Name"`
	HXTarget      string `json:"HX-Target"`
	HXCurrentURL  string `json:"HX-Current-URL"`
}

type HTMXRequest struct {
	Template string         `json:"template"`
	Message  string         `json:"message"`
	Headers  RequestHeaders `json:"HEADERS"`
}
