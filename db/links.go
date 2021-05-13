package db

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
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

type ListLinksRequest struct {
	ID     *int
	Limit  *int
	Offset *int
}

func ListLinks(req *ListLinksRequest) ([]*model.Link, error) {
	q := sq.Select(
		"users.id",
		"users.username",
		"links.id",
		"links.title",
		"links.address",
	).
		From("links").
		Join("users ON links.creator_user_id = users.id")

	if req.ID != nil {
		q = q.Where(sq.Eq{"links.id": *req.ID})
	}

	if req.Limit != nil {
		q = q.Limit(uint64(*req.Limit))
	} else {
		q = q.Limit(10)
	}

	if req.Offset != nil {
		q = q.Offset(uint64(*req.Offset))
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list links sql: %w", err)
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying list links sql: %w", err)
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
