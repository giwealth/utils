package mail

import "testing"

func TestMail(t *testing.T)  {
	mail := New("admin@gmail.com", "password", "smtp.gmail.com", 25)
	if err := mail.Send("txt testing", "Hello world", nil, []string{"abc@gmail.com"}); err != nil {
		panic(err)
	}

	if err := mail.Send("attach testing", "Hello world", []string{"a.txt"}, []string{"abc@gmail.com"}); err != nil {
		panic(err)
	}

	if err := mail.Send("html testing", "Hello <b>Bob</b> and <i>Cora</i>!", nil, []string{"abc@gmail.com"}); err != nil {
		panic(err)
	}
}