package errors

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ProxyError struct {
	Err Error `json:"proxyError"`
}

var (
	InvalidJSON = &ProxyError{
		Err: Error{
			"Not a JSON encoded string",
			0x000001,
		},
	}

	Hostname = &ProxyError{
		Err: Error{
			"Cannot resolve hostname",
			0x000002,
		},
	}

	ConnectionFailed = &ProxyError{
		Err: Error{
			"Connection to the host failed",
			0x000003,
		},
	}
)
