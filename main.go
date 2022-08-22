package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "$TheRooney88A"
	hostname = "127.0.0.1:3306"
	dbname   = "commuter"
)

type LevelData struct {
	distance   int
	numEnemies int
	laneWidth  int
	roadBottom int
	roadTop    int
	level      int
}

type Enemy struct {
	level          int
	launchDistance int
	launchX        int
	launchY        int
	kind           string
	jam            int
	enemyIndex     int
}

func main() {
	fmt.Println(("What would you like to do. \n1 : Read\n2 : Write\n3 : Delete"))

	var choice int
	fmt.Scanf("%d", &choice)

	var db = connectToDB("")
	if db == nil {
		log.Fatal("db is null")
	}

	if choice == 1 { //read
		var innerChoice int
		fmt.Print("What level do you want to read?\n")
		fmt.Scanf("%d", &innerChoice)
		if innerChoice < 1 {
			log.Fatal("level choice must be gretaer than 1")
		}

		fmt.Println("What do you want to read \n1 : Level Data\n2 : Enemy List")
		fmt.Scanf("%d", &innerChoice)

		if innerChoice != 1 && innerChoice != 2 {
			log.Fatal("Choice must be in the list.")
		}

		if innerChoice == 1 {
			ld := readLevelData(db, innerChoice)
			fmt.Printf("%d", ld.distance)
		} else {
			readEnemyData(db, innerChoice)
		}
	} else if choice == 2 { // write
		var input string

		enterAnother := true

		for enterAnother {
			fmt.Println("Enter enemy data.\n(level|launchDistance|launchX|launchY|jam|kind)")
			fmt.Scanf("%s", &input)

			values := strings.Split(input, "|")
			fmt.Println(values)
			level, err := strconv.Atoi(values[0])
			if err != nil {
				log.Fatal(("Fatal1"))
			}
			launchDistance, err := strconv.Atoi(values[1])
			if err != nil {
				log.Fatal(("Fatal2"))
			}
			launchX, err := strconv.Atoi(values[2])
			if err != nil {
				log.Fatal(("Fatal3"))
			}
			launchY, err := strconv.Atoi(values[3])
			if err != nil {
				log.Fatal(("Fatal4"))
			}

			jam, err := strconv.Atoi(values[4])
			if err != nil {
				log.Fatal("Fatal5")
			}

			kind := values[5]

			tempEnemy := Enemy{level: level, launchDistance: launchDistance, launchX: launchX, launchY: launchY, jam: jam, kind: kind, enemyIndex: 0}

			insertEnemy(db, tempEnemy)
			enterAnother = false
		} //enterAnother

	} else if choice == 3 { //Delete
		fmt.Println("choice was 3")
	}
}

func readLevelData(db *sql.DB, levelChoice int) LevelData {
	var ld LevelData
	rows, err := db.Query("SELECT distance FROM levelData WHERE level = " + strconv.Itoa(levelChoice))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var distance int
		if err := rows.Scan(&distance); err != nil {
			log.Fatal(err)
		}
		ld.distance = distance
	}

	return ld
}

func insertLevelData(db *sql.DB, levelData LevelData) {

}

func readEnemyData(db *sql.DB, levelChoice int) {
	rows, err := db.Query("SELECT * FROM enemy;")
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var level int
		var launchDistance int
		var launchX int
		var launchY int
		var jam int
		var kind string
		var enemeyIndex int

		if err := rows.Scan(&level, &launchDistance, &launchY, &launchX, &jam, &kind, &enemeyIndex); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Level : %v Launch Distance : %v Launch X : %v Launch Y : %v Jam : %v Kind : %v Index : %v\n", level, launchDistance, launchX, launchY, jam, kind, enemeyIndex)
	}

}

func insertEnemy(db *sql.DB, enemy Enemy) {
	insertString := fmt.Sprintf("INSERT INTO enemy VALUES(%v, %v, %v, %v, %v, \"%v\",0);",
		enemy.level, enemy.launchDistance, enemy.launchX, enemy.launchY, enemy.jam, enemy.kind)

	res, err := db.Exec(insertString)
	if err != nil {
		panic(err.Error())
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Enemy succesfully added to level : %v", lastId)
}

func connectToDB(connectString string) *sql.DB {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		log.Print("error connecting to db")
		return db
	}
	//defer db.Close()

	return db

}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}
