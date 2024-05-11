package api_rest

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"testing"
)

// cd internal\usecase\api_rest
// go test -run GetUserLogin
func Test_GetUserLogin(t *testing.T) {

	type (
		ConfigC3po struct {
			C3poUrl string `env-required:"true"   env:"C3PO_URL"`
		}
	)
	cfg := &ConfigC3po{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Println("Ошибка получения данных из конфига")
		return //nil, err
	} else {
		log.Println(cfg.C3poUrl)
	}

	/*
		unifiC3po := &UnifiC3po{
			client: http.Client{
				Timeout: 10 * time.Second,
			},
			url: "",
		}*/
	c3po := NewC3po(cfg.C3poUrl)

	client1 := &entity.Client{
		Hostname: "NBCN-GUSIKHIN1",
	}
	client2 := &entity.Client{
		Hostname: "",
	}
	client3 := &entity.Client{
		Hostname: "mock",
	}
	//url1 := unifiC3po.url + "pc/" + client1.Hostname
	//t.Logf(url)

	//пустой логин пользователя
	err = c3po.GetUserLogin(client1)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf(client1.Hostname)
		if client1.UserLogin == "" {
			t.Logf("Hostname is empty")
		}
		if client1.UserLogin == " " {
			t.Logf("Hostname is one space")
		}
	}

	//Пустое имя компьютера
	err = c3po.GetUserLogin(client2)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("EMPTY : ")
	}

	//Не корректное имя компьютера
	err = c3po.GetUserLogin(client3)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("MOCK : ")
	}

}
