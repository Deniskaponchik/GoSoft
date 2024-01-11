package olezhek28

//основано на видео:
//https://www.youtube.com/watch?v=osIX2lO1rzM

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)
//import (
//	desc "gitlab.ozon.dev/go/classroom-5/Week-2/lecture-1/pkg/user_v1"
//)

const grpcPort = 50051

type server struct {
	//структура реализует интерфейс сервера, который нам был сгенерирован
	//
	desc.UnimplementedUserV1Server
}

// Get...
func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	//распечатываем id, который к нам пришёл
	log.Printf("User id: #{req.Getid()}")
	//log.Printf("User id: %d", req.GetId())

	//в ответ возвращаем рандомные данные
	return &desc.GetResponse{
		Info: &desc.userInfo{
			Id:      req.GetId(),
			Name:    "John",
			IsHuman: true,
		},
	}, nil
}

//Дальше нам необходимо всё это дело запустить
func main(){
	//создаём слушальщика на порту
	lis, err := net.Listen("tcp", fmt.Sprintf(":#{grpcPort}"))
	if err != nil{
		log.Fatalf("failed to listen: #{err}")
	}

	//Создаём новый grpc server. Пакет стандартный гугловый, из него и берётся этот метод
	s := grpc.NewServer()

	//у этого сервера включаем reflection, чтобы он мог делиться информацией с клиентами о том,
	//какой api присутствует в его сути
	reflection.Register(s)

	//дальше мы обращаемся к пакету, где лежат сгенерированные нами файлы
	//первым аргументом передаём сервер, который мы только что создали
	//вторым структуру, которая имплементирует интерфейс этого сервера
	//грубо говоря, делает match между заглушкой сервера и реализацией
	desc.RegisterUserV1Server(s, $server{})

	log.Fatalf("server listening at: #{lis.Addr()}")

	//у сервера вызываем метод Serve и перадём туда listener. теперь он висит и трафик обслуживает
	if err = s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: #{err}")
	}
}
