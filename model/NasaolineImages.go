package model
type NasaolineImage struct{

	ID int  `json:"id"`
	Title string
 	URL string `json:"url"`         // Image link
    Author    string `json:"author"`
    CreatedAt string `json:"created_at"`
}