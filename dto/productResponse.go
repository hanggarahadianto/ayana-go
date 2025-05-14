package dto

type HomeByClusterResponse struct {
	ID         string          `json:"id"`
	Title      string          `json:"title"`
	Type       string          `json:"type"`
	Content    string          `json:"content"`
	Maps       string          `json:"maps"`
	Price      float64         `json:"price"`
	Status     string          `json:"status"`
	Square     float64         `json:"square"`
	Bathroom   float64         `json:"bathroom"`
	Bedroom    float64         `json:"bedroom"`
	StartPrice float64         `json:"start_price"`
	Cluster    ClusterResponse `json:"cluster"`
	NearBies   []NearBy        `json:"near_bies"`
}

type ClusterResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Maps     string `json:"maps"`
}

type NearBy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Distance string `json:"distance"`
}
