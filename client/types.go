package client

type HistoryEntry struct {
	Query     string `json:"query"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updated_at"`
}

type historyRequest struct {
	URL    string `json:"url"`
	Title  string `json:"title,omitempty"`
	Query  string `json:"query"`
	Delete bool   `json:"delete,omitempty"`
}

type RulesResponse struct {
	Skip     []string          `json:"skip"`
	Priority []string          `json:"priority"`
	Aliases  map[string]string `json:"aliases"`
}
