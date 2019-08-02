package shortener

import (
	"database/sql"
	"fmt"
)

func CreateLinkStore(db *sql.DB) DbLinkStore {
	return DbLinkStore{
		db,
	}
}

type DbLinkStore struct {
	db *sql.DB
}

type LinkElement struct {
	URLID           string `json:"url_id"`
	Target          string `json:"target"`
	CreatorID       int    `json:"creator_id"`
	CreatorName     string `json:"creator_name"`
	TotalVisitCount int    `json:"total_visit_count"`
	TopVisitorScore int    `json:"top_visitor_score"`
	TopVisitorID    int    `json:"top_visitor_id"`
	TopVisitorName  string `json:"top_visitor_name"`
}

func (l DbLinkStore) GetLink(id Identifier, ip_address string, user_id int) (link Link, err error) {
	getQuery := `SELECT target FROM url WHERE id=$1`

	err = l.db.QueryRow(getQuery, id).Scan(&link)

	if link == "" {
		err = ErrLinkNotFound
		return "", err
	}

	insertVisit(l.db, user_id, id, ip_address)

	return
}

func (l DbLinkStore) GetLinks(limit, offset int, orderBy string) ([]LinkElement, error) {
	data, err := l.db.Query(linksQuery, limit, offset)
	if err != nil {
		println(err.Error())
	}
	links := []LinkElement{}

	for data.Next() {
		var (
			url_id            string
			target            string
			creator_id        int
			creator_name      string
			total_visit_count int
			top_visitor_score int
			top_visitor_id    int
			top_visitor_name  string
		)
		data.Scan(&url_id, &target, &creator_id, &creator_name, &total_visit_count, &top_visitor_score, &top_visitor_id, &top_visitor_name)

		links = append(links, LinkElement{url_id, target, creator_id, creator_name, total_visit_count, top_visitor_score, top_visitor_id, top_visitor_name})
	}

	return links, nil
}

func (l DbLinkStore) DeleteLink(id Identifier) (err error) {
	deleteQuery := `DELETE FROM url WHERE id = $1 RETURNING id`
	err = l.db.QueryRow(deleteQuery, id).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrLinkNotFound
		}
		return err
	}

	return
}

func insertVisit(db *sql.DB, account_id int, url_id, ip_address string) (id int) {
	insertQuery := `INSERT INTO visit (url_id, account_id, ip_address) values ($1, $2, $3) RETURNING id`

	err := db.QueryRow(insertQuery, url_id, 1, ip_address).Scan(&id)

	if err != nil {
		fmt.Print(err)
	}

	return
}

func (l DbLinkStore) SaveLink(data *LinkPost, accountID int) error {
	insertQuery := `INSERT INTO url (id, target, account_id) values ($1, $2, $3) RETURNING id`
	var id string
	err := l.db.QueryRow(insertQuery, data.ID, data.Target, accountID).Scan(&id)

	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

var linksQuery = `SELECT
	url_id,
	target,
	account_id as creator_id,
	creator_name,
	visit_count as total_visit_count,
	top_visits as top_visitor_score,
	top_visitor as top_visitor_id,
	top_visitor_name      
FROM
(SELECT
	*          
FROM
	(SELECT
		url.id as url_id,
		url.target,
		url.account_id,
		account.name as creator_name              
	FROM
		url              
	INNER JOIN
		account                      
			ON account.id = account_id) t1          
JOIN
	(
		SELECT
			COUNT(url_id) as visit_count,
			url_id as url_id_v                  
		FROM
			visit                  
		GROUP BY
			visit.url_id             
	) t2                  
		ON t1.url_id = t2.url_id_v             
	) as leftSide      
JOIN
(
	SELECT
		url_id as url_id_r,
		cnt as top_visits,
		account_id as top_visitor,
		top_visitor_name              
	FROM
		( SELECT
			url_id,
			cnt,
			account_id,
			top_visitor_name,
			RANK() OVER (PARTITION                  
		BY
			url_id                  
		ORDER BY
			cnt DESC) AS rn                  
		FROM
			(SELECT
				url_id,
				account_id,
				account.name as top_visitor_name,
				COUNT(account_id) as cnt                      
			FROM
				visit                      
			INNER JOIN
				account 
					on account_id = account.id                      
			GROUP BY
				url_id,
				account_id,
				account.name ) t ) s                  
		WHERE
			s.rn = 1                                      ) as rightSide                      
			ON leftSide.url_id = rightSide.url_id_r     
	ORDER BY
		visit_count DESC LIMIT $1 OFFSET $2;`
