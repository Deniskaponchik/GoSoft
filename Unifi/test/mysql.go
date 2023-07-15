package main

import (
	_ "github.com/go-sql-driver/mysql"
)

//https://tutorialedge.net/golang/golang-mysql-tutorial/

/* Many ROWS.
type PC struct {
	//ID   int    `json:"id"`
	UserName string `json:"user_name"`
	Date_Mod string `json:"date_mod"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	//results, err := db.Query("SELECT id, contact login FROM glpi_db.glpi_computers where name = 'NBAB-FISHER'")
	results, err := db.Query("SELECT date_mod, contact FROM glpi_db.glpi_computers where name = 'NBAB-FISHER' ORDER BY date_mod DESC")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var pc PC
		// for each row, scan the result into our tag composite object
		err = results.Scan(&pc.Date_Mod, &pc.UserName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//log.Printf(tag.Date_Mod)
		log.Printf(pc.UserName)
		break
	}
}*/
/* Many ROWS. Полностью рабочее. !!! НЕ РЕДАКТИРОВАТЬ !!!
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Login string `json:"name"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT id, contact login FROM glpi_db.glpi_computers where name = 'NBAB-FISHER'")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.ID, &tag.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(tag.Name)
		break
	}
}*/
/*One row. Рабочее
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Login string `json:"name"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	pcname := "NBAB-FISHER"
	fmt.Println(pcname)
	var tag Tag
	err = db.QueryRow("SELECT id, contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC", pcname).Scan(&tag.ID, &tag.Name)
	// после запятой указываем значение, которое будет подставляться заметсо вопроса
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Println(tag.ID)
	log.Println(tag.Name)
}*/
/*One row. Рабочее. не менять
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Login string `json:"name"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var tag Tag
	err = db.QueryRow("SELECT id, contact login FROM glpi_db.glpi_computers where id = ?", 2).Scan(&tag.ID, &tag.Name)
	// после запятой указываем значение, которое будет подставляться заметсо вопроса
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Println(tag.ID)
	log.Println(tag.Name)
}*/
