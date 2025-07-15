package dto

type PerformerResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TotalBooking int    `json:"total_booking"`
	TotalAmount  int64  `json:"total_amount"`
}

type MarketerPerformanceResponse struct {
	TopPerformers   []PerformerResponse `json:"top_performers"`
	UnderPerformers []PerformerResponse `json:"under_performers"`
}
