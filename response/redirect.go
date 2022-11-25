package response

// Redirect can be used as handler response to redirect
type Redirect struct {
	URL        string
	StatusCode int
}
