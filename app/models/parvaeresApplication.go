package models

type ParvaeresApplicationData struct {
	Email      string
	Folder     string
	Repository string
}

type ParvaeresAPIResponse struct {
	Message string
	Items   []ParvaeresResponseItem
	Error   bool
}

type ParvaeresResponseItem struct {
	UUID   string
	Status string
}
