package dto

type StartRecommendationResponse struct {
	SessionID  int64  `json:"session_id"`
	Intro      string `json:"intro"`
	Question   string `json:"question"`
	Step       int    `json:"step"`        // 1-indexed, for display
	TotalSteps int    `json:"total_steps"`
}

type AnswerRecommendationRequest struct {
	Answer string `json:"answer"`
}

func (r *AnswerRecommendationRequest) Validate() error {
	if r.Answer == "" {
		return errField("answer is required")
	}
	return nil
}

type AnswerRecommendationResponse struct {
	SessionID      int64  `json:"session_id"`
	Status         string `json:"status"` // "in_progress" | "completed"
	Question       string `json:"question,omitempty"`
	Step           int    `json:"step,omitempty"`
	TotalSteps     int    `json:"total_steps,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}
