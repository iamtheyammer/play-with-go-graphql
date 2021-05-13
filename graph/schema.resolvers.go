package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/iamtheyammer/play-with-go-graphql/auth"
	"github.com/iamtheyammer/play-with-go-graphql/db"
	"github.com/iamtheyammer/play-with-go-graphql/graph/generated"
	"github.com/iamtheyammer/play-with-go-graphql/graph/model"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	uID := auth.UserIDFromContext(ctx)
	if uID == 0 {
		return nil, fmt.Errorf("no authentication provided")
	}

	link, err := db.CreateLink(input, 1)
	if err != nil {
		return nil, fmt.Errorf("error adding new link to database: %w", err)
	}

	return link, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	// password: length check
	if len(input.Password) < 1 {
		return "", fmt.Errorf("password must be more than zero characters")
	}

	// insert user
	userId, err := db.InsertUser(input.Username, input.Password)
	if err != nil {
		return "", fmt.Errorf("error inserting user: %w", err)
	}

	// gen jwt
	jwt, err := auth.GenerateToken(userId)
	if err != nil {
		return "", fmt.Errorf("error generating user jwt: %w", err)
	}

	return jwt, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	if len(input.Username) < 1 || len(input.Password) < 1 {
		return "", fmt.Errorf("username and password must be greater than zero characters")
	}

	// compare password
	user, err := db.GetUserByUsername(input.Username)
	if err != nil {
		return "", fmt.Errorf("error getting existing user: %w", err)
	}

	if user == nil {
		return "", fmt.Errorf("a user with that username does not exist")
	}

	now := time.Now()
	if validPw := db.ComparePasswordToHash(user.Password, input.Password); !validPw {
		return "", fmt.Errorf("invalid username/password combination")
	}
	fmt.Printf("db.ComparePasswordToHash took %dms\n", time.Since(now).Milliseconds())

	// valid, issue jwt
	jwt, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("error generating jwt: %w", err)
	}

	return jwt, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	// issue new jwt, return
	uID, err := auth.ParseToken(input.Token)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	if uID < 1 {
		uID = auth.UserIDFromContext(ctx)
	}

	if uID < 1 {
		return "", fmt.Errorf("no valid authentication present")
	}

	jwt, err := auth.GenerateToken(uID)
	if err != nil {
		return "", fmt.Errorf("error generating new token: %w", err)
	}

	return jwt, nil
}

func (r *queryResolver) Links(ctx context.Context, id *string, limit *int, offset *int) ([]*model.Link, error) {
	var intId *int

	if id != nil {
		convertedId, err := strconv.Atoi(*id)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}

		intId = &convertedId
	}

	links, err := db.ListLinks(&db.ListLinksRequest{
		ID:     intId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing links: %w", err)
	}

	return links, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
