package models

type UserReaction struct {
	UserID   int      `json:"user_id"`
	Reaction Reaction `json:"reaction"`
	EventID  int      `json:"event_id"`
}
