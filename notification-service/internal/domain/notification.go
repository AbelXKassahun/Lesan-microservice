package domain

type Notification struct {
    UserID   string   `json:"user_id"`
    Message  string   `json:"message"`
    Channels []string `json:"channels"` // ["email"], later ["push"]
    Email    string   `json:"email,omitempty"`
    Token    string   `json:"token,omitempty"` // For push notifications in future
}
