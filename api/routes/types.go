package routes

type UserData struct {
	ID        string
	Email     string
	Name      string
	SessionID string
}
type UserDataRoute struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RoomData struct {
	ID        string
	Name      string
	CreatedBy string
	CreatedAt string
}

type RoomDataRoute struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

type PlayersData struct {
	ID       string
	RoomID   string
	PlayerID string
}
type PlayersDataRoute struct {
	RoomID   string `json:"room_id"`
	PlayerID string `json:"player_id"`
}

type RoundData struct {
	ID          string
	RoomID      string
	Calculation string
	AnsweredBy  string
	CreatedAt   string
}

type RoundDataRoute struct {
	RoomID      string `json:"room_id"`
	Calculation string `json:"calculation"`
	AnsweredBy  string `json:"answered_by"`
}
type SessionData struct {
	ID          string
	IP          string
	UserAgent   string
	UserID      string
	ConnectedAt string
}

type SessionDataRoute struct {
	UserID string `json:"user_id"`
}
