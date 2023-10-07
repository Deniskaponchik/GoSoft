package main

type CreateUserRequest struct {
	Email    string `xml:"Email,omitempty"`
	Password string `xml:"Password,omitempty"`
}

type CreateUserResponse struct {
	ID string `xml:"ID"`
}

/*
func main() {
	client := NewSOAPClient("https://soap.example.com/call", true, nil)

	req := &CreateUserRequest{
		Email:    "jdoe@example.com",
		Password: "1234567890",
	}
	res := &CreateUserResponse{}
	if err := client.Call("create_user", req, res); err != nil {
		panic(err)
	}

	// if everything went well res.ID should have its
	// value set with the one returned by the service.
	fmt.Println(res.ID)
}
*/
