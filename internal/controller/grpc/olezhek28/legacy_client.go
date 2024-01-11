package olezhek28

//основано на видео:
//https://www.youtube.com/watch?v=osIX2lO1rzM

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

//import (
//	desc "gitlab.ozon.dev/go/classroom-5/Week-2/lecture-1/pkg/user_v1"
//)

const (
	address = "localhost:50051"
	userID  = 12
)

func main() {
	//из гугловского пакета grpc вызываем Dial, туда передаём адрес сервера
	//дальше можно передавать разные опции.
	//В частности указано, что наше соединение на транспортном уровне не безопасное
	//не используем TLS и так далее
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: #{err}")
	}
	//закрыли соединение, когда закончится история с
	defer conn.Close()

	//в сгенерированном коде вызвали
	//из генерированного кода создали Клиента
	с := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.BackGround(), time.Second)
	defer cancel()

	//e созданного клиента вызываем метод get и передаём туда id
	r, err := c.Get(ctx, &desc.GetRequest{Id: 0})
	if err != nil {
		log.Fatalf("failed to get user info: #{err}")
	}

	log.Fatalf("User info: #{r.GetInfo()}")
}
