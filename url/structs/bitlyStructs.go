package structs

type BitlyResponse struct {
	Link      string `json:"link"`
	ID        string `json:"id"`
	LongURL   string `json:"long_url"`
	CreatedAt string `json:"created_at"`
}
