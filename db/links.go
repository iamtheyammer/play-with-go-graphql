package db

import (
	"fmt"
	"github.com/iamtheyammer/play-with-go-graphql/graph/model"
	"strconv"
)

func CreateLink(link model.NewLink, userId int) (*model.Link, error) {
	stmt, err := db.Prepare("INSERT INTO links (creator_user_id, title, address) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing insert link sql: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(userId, link.Title, link.Address)
	if err != nil {
		return nil, fmt.Errorf("error inserting link: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error retreiving last insert id: %w", err)
	}

	return &model.Link{
		ID:      strconv.FormatInt(id, 10),
		Title:   link.Title,
		Address: link.Address,
		User: &model.User{
			ID: strconv.Itoa(userId),
		},
	}, nil
}

func ListLinks(limit int) ([]*model.Link, error) {
	stmt, err := db.Prepare("SELECT users.id, users.username, links.id, links.title, links.address FROM links JOIN users ON links.creator_user_id = users.id LIMIT (?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing list links sql: %w", err)
	}

	defer stmt.Close()

	if limit == 0 {
		limit = 10
	}

	rows, err := stmt.Query(limit)
	if err != nil {
		return nil, fmt.Errorf("error opening links rows for reading: %w", err)
	}

	defer rows.Close()

	var links []*model.Link

	for rows.Next() {
		l := model.Link{}
		u := model.User{}

		err := rows.Scan(&u.ID, &u.Name, &l.ID, &l.Title, &l.Address)
		if err != nil {
			return nil, fmt.Errorf("error scanning links: %w", err)
		}

		l.User = &u

		links = append(links, &l)
	}

	return links, nil
}
