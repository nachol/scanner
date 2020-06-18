package model

/*
H1response : Program Of H1
*/
type H1Pictureresponse struct {
	Data data `json:"data"`
}

type data struct {
	Resource resource `json:"resource"`
}

type resource struct {
	Profile_picture string `json:"profile_picture"`
}
