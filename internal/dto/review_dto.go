package dto

type CreateReviewRequest struct {
	Rating float64 `json:"rating"`
	Body   string  `json:"body"`
}

func (r *CreateReviewRequest) Validate() error {
	if r.Rating < 1.0 || r.Rating > 5.0 {
		return errField("rating must be between 1.0 and 5.0")
	}
	if r.Body == "" {
		return errField("review body is required")
	}
	return nil
}

type UpdateReviewRequest struct {
	Rating *float64 `json:"rating"`
	Body   *string  `json:"body"`
}

func (r *UpdateReviewRequest) Validate() error {
	if r.Rating != nil && (*r.Rating < 1.0 || *r.Rating > 5.0) {
		return errField("rating must be between 1.0 and 5.0")
	}
	return nil
}
