package models

type ParvaeresApplicationData struct {
	Email      string
	Path       string
	Repository string
}

type ParvaeresAPIResponse struct {
	Message string
	Items   []ParvaeresResponseItem
	Error   bool
}

type ParvaeresResponseItem struct {
	UUID    string
	RepoURL string
	Path    string
	Email   string
	Status  string
	LogsURL string
	Errors  []string
}
