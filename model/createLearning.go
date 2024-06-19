

type Create_req struct {
	User_id     string `json:"user_id"`
	Word_number int    `json:"word_number"`
}

type Create_res struct {
	sessionId  int
	Data []struct {
		Answer string
	}
}


func CreateLearning() {

}