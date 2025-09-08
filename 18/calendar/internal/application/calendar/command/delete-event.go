package command

type DeleteEvent struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}
