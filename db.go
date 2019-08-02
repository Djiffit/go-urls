package shortener

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "url"
)

func InitDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	fmt.Print(db.Stats())
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db")

	createTables(db)

	return db
}

var sqlStatements = []string{
	`DROP TABLE url CASCADE;`,
	`DROP TABLE account CASCADE;`,
	`DROP TABLE visit CASCADE;`,
	`CREATE TABLE account(
		ID SERIAL PRIMARY KEY,
		name varchar(100) UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		password varchar(300) NOT NULL);`,
	`CREATE TABLE url (
		ID varchar(40) PRIMARY KEY NOT NULL,
		target varchar(500) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		account_id integer,
		FOREIGN KEY (account_id) REFERENCES account(ID) ON DELETE CASCADE);`,
	`CREATE TABLE visit
		(ID SERIAL PRIMARY KEY,
		account_id integer,
		url_id varchar(40),
		ip_address varchar(100),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (url_id) REFERENCES url(ID) ON DELETE CASCADE,
		FOREIGN KEY (account_id) references account(ID) ON DELETE CASCADE);`,
	`INSERT INTO account(name, password) values ('konstaku', '34223423');`,
	`INSERT INTO url(id, target, account_id) values('testing', 'www.helsinki.fi', 1);`,
}

func createTables(db *sql.DB) {

	for _, statement := range sqlStatements {
		_, err := db.Exec(statement)
		if err != nil {
			fmt.Print(err)
		}
	}
	db.Exec(`INSERT INTO account(name, password) values ('kunsateo', '32423432');`)
	db.Exec(`INSERT INTO account(name, password) values ('sekorun', '324324234');`)
	db.Exec(`INSERT INTO account(name, password) values ('nirosek', '324234324');`)
	db.Exec(`INSERT INTO account(name, password) values ('ahmorin', '324234');`)

	parts := []string{
		"google",
		"facebook",
		"gmail",
		"reddit",
		"plex",
		"news.ycombinator",
		"mit",
		"helsinki",
		"amazon",
		"wikipedia",
		"twitter",
		"instagram",
		"alipay",
		"aliexpress",
		"microsoft",
		"taobao",
		"yahoo",
		"imdb",
	}

	ends := []string{
		"fi",
		"co.uk",
		"com",
		"xyz",
		"se",
		"dk",
		"de",
		"gov",
	}

	ips := []string{
		`81.147.35.104`,
		`154.151.115.137`,
		`57.133.51.57`,
		`95.196.52.79`,
		`100.12.107.219`,
		`56.47.31.202`,
		`251.202.202.15`,
		`109.144.233.7`,
		`126.214.10.217`,
		`238.60.116.149`,
		`158.102.178.31`,
		`201.195.98.247`,
		`199.101.35.32`,
		`247.126.15.206`,
		`28.103.10.120`,
		`34.154.175.176`,
		`72.120.117.148`,
		`95.139.27.114`,
		`174.109.138.33`,
		`47.25.96.218`,
	}

	// Get sql injected methinks ;^)
	sqlState := `INSERT INTO url(id, target, account_id) values `
	for j, part := range parts {
		for i, end := range ends {
			if j == 0 && i == 0 {
				sqlState += fmt.Sprintf(` ('%s', '%s', %d)`, part+end, `www.`+part+`.`+end, randNum(5.0))
			} else {
				sqlState += fmt.Sprintf(`, ('%s', '%s', %d)`, part+end, `www.`+part+`.`+end, randNum(5.0))
			}
		}
	}
	sqlState += ";"
	db.Query(sqlState)

	sqlState = `INSERT INTO visit(account_id, url_id, ip_address) values `
	for i := 0; i < 10000; i++ {
		if i == 0 {
			sqlState += fmt.Sprintf(` (%d, '%s', '%s') `, randNum(5.0), parts[randNum(float64(len(parts))-1)]+ends[randNum(float64(len(ends)-1))], ips[randNum(float64(len(ips)-1))])
		} else {
			sqlState += fmt.Sprintf(`, (%d, '%s', '%s') `, randNum(5.0), parts[randNum(float64(len(parts))-1)]+ends[randNum(float64(len(ends)-1))], ips[randNum(float64(len(ips)-1))])
		}
	}
	sqlState += ";"
	db.Query(sqlState)

}

func randNum(max float64) int {
	return int(math.Ceil(rand.Float64() * max))
}
