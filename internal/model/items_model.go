package model

type ItemDTO struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	ImageURL       string   `json:"image_url"`
	Description    string   `json:"description"`
	Price          float32  `json:"price"`
	Rating         float32  `json:"rating"`
	Specifications []string `json:"specifications"`
}
