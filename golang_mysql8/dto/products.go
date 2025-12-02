package dto

type Products struct {
	Id             string `json:"id"`
	Category       string `json:"category"`
	Descriptions   string `json:"descriptions"`
	Qty            string `json:"qty"`
	Unit           string `json:"unit"`
	Costprice      string `json:"costprice"`
	Sellprice      string `json:"sellprice"`
	Saleprice      string `json:"saleprice"`
	Productpicture string `json:"productpicture"`
	Alertstocks    string `json:"alertstocks"`
	Criticalstocks string `json:"criticalstocks"`
}
