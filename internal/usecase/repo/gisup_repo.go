package repo

//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://github.com/evrone/go-clean-template/blob/master/internal/usecase/repo/translation_postgres.go
import (
	//"database/sql"
	"github.com/deniskaponchik/GoSoft/pkg/mysql1"
	"github.com/deniskaponchik/GoSoft/pkg/postgres"
	"github.com/deniskaponchik/GoSoft/pkg/redis1"
	_ "github.com/go-sql-driver/mysql" //для установки драйвера mysql. Сам пакет как бы не используется явно в коде
	//"log"
)

type Repo struct {
	repoGlpi  mysql1.SqlMy
	repoMySql mysql1.SqlMy
	//repoGisupMySqlProd mysql1.SqlMy
	//repoGisupMySqlTest mysql1.SqlMy
	repoPG    postgres.Postgres
	repoRedis redis1.Redis
}

// реализуем Инъекцию зависимостей DI. Используется в app
// repoGisupMySqlProd mysql1.SqlMy, repoGisupMySqlTest mysql1.SqlMy
func NewRepo(repoGlpi mysql1.SqlMy, repoMySql mysql1.SqlMy, repoPG postgres.Postgres, repoRedis redis1.Redis) (*Repo, error) {

	repo := &Repo{
		repoGlpi:  repoGlpi,
		repoMySql: repoMySql,
		//repoGisupMySqlProd: repoGisupMySqlProd,
		//repoGisupMySqlTest: repoGisupMySqlTest,
		repoPG:    repoPG,
		repoRedis: repoRedis,
	}
	return repo, nil
}
