package external

import (
	"context"
	"fmt"
	"movie-festival/constants"
	"movie-festival/external/proto/tokenvalidation"
	"movie-festival/helpers"
	"movie-festival/internal/models"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	var (
		resp models.TokenData
	)

	conn, err := grpc.Dial(helpers.GetEnv("UMS_GRPC_HOST", ""), grpc.WithInsecure())
	if err != nil {
		return resp, errors.Wrap(err, "failed to dial ums grpc")
	}
	defer conn.Close()

	client := tokenvalidation.NewTokenValidationClient(conn)

	req := &tokenvalidation.TokenRequest{
		Token: token,
	}
	response, err := client.ValidateToken(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to validate token")
	}
	if response.Message != constants.SuccessMessage {
		return resp, fmt.Errorf("got response error from ums: %s", response.Message)
	}

	resp.UserID = response.Data.UserId
	resp.Username = response.Data.Username
	resp.FullName = response.Data.FullName
	resp.Email = response.Data.Email

	return resp, nil
}
