package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sparrc/go-ping"
	"log"
	"os"
	"time"
)

var (
	ip_address string
	id         int
)

func main() {
	start := time.Now()

	db_name := goDotEnvVariable("DB_NAME")
	db_user := goDotEnvVariable("DB_USER")
	db_password := goDotEnvVariable("DB_PASSWORD")
	db_hostname := goDotEnvVariable("DB_HOSTNAME")
	db_port := goDotEnvVariable("DB_PORT")

	db_string := db_user + ":" + db_password + "@tcp(" + db_hostname + ":" + db_port + ")/" + db_name

	conn, _ := sql.Open("mysql", db_string)

	rows, _ := conn.Query("select id,ip_address from sites")

	for rows.Next() {

		if err := rows.Scan(&id, &ip_address); err != nil {
			log.Fatal(err)
		}

		if len(ip_address) != 0 {
			go pingHost(conn, ip_address, id)
			time.Sleep(150 * time.Millisecond)
		}
	}

	elapsedTime := time.Since(start)

	fmt.Println("\nTime: " + elapsedTime.String())

	defer conn.Close()
}

func pingHost(conn *sql.DB, ip_address string, site_id int) {
	pinger, _ := ping.NewPinger(ip_address)

	pinger.Count = 3
	pinger.Interval = 1
	pinger.Run()
	fmt.Printf(".")
	average := pinger.Statistics().MaxRtt.Microseconds()

	raw := `INSERT INTO sites_ping_responses (site_id, first, second, third, average) VALUES (?,?,?,?,?)`

	conn.Exec(raw, id, 0, 0, 0, average)
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
