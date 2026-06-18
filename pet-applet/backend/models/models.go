package models

type Pet struct {
	ID        string `json:"id"`
	Avatar    string `json:"avatar"`
	Name      string `json:"name"`
	Breed     string `json:"breed"`
	Birthday  string `json:"birthday"`
	Weight    string `json:"weight"`
	Notes     string `json:"notes"`
	CreatedAt int64  `json:"createdAt"`
}

type FeedingSchedule struct {
	ID       string `json:"id"`
	PetID    string `json:"petId"`
	Time     string `json:"time"`
	FoodType string `json:"foodType"`
	Amount   string `json:"amount"`
}

type FeedingRecord struct {
	ID         string `json:"id"`
	PetID      string `json:"petId"`
	ScheduleID string `json:"scheduleId"`
	Time       string `json:"time"`
	FoodType   string `json:"foodType"`
	Amount     string `json:"amount"`
	Notes      string `json:"notes"`
	CreatedAt  int64  `json:"createdAt"`
}
