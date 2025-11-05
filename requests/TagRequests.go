package requests

type TagBodyRequest struct {
	Name        string `json:"name" example:"Risky"`
	Description string `json:"description" example:"PickupLine that could send you to jail"`
}

type TagIdRequest struct {
	ID string `json:"id" example:"01976595-4044-7426-b6fa-f64173211b94"`
}
