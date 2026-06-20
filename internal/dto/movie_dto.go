package dto

type CreateMovieRequest struct {
	Title     string `json:"title"`
	Director  string `json:"director"`
	Genre     string `json:"genre"`
	Year      int    `json:"year"`
	Synopsis  string `json:"synopsis"`
	PosterURL string `json:"poster_url"`
}

func (r *CreateMovieRequest) Validate() error {
	if r.Title == "" {
		return errField("title is required")
	}
	if r.Director == "" {
		return errField("director is required")
	}
	if r.Year < 1888 || r.Year > 2100 {
		return errField("year is invalid")
	}
	return nil
}

type UpdateMovieRequest struct {
	Title     *string `json:"title"`
	Director  *string `json:"director"`
	Genre     *string `json:"genre"`
	Year      *int    `json:"year"`
	Synopsis  *string `json:"synopsis"`
	PosterURL *string `json:"poster_url"`
}
