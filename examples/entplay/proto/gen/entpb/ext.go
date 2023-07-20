package entpb

import (
	ent "algorithms/examples/entplay/ent"
	user "algorithms/examples/entplay/ent/user"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ExtService implements ExtServiceServer.
type ExtService struct {
	client *ent.Client
	UnimplementedExtServiceServer
}

// TopUser returns the user with the highest ID.
func (s *ExtService) TopUser(ctx context.Context, _ *emptypb.Empty) (*User, error) {
	id := s.client.User.Query().Aggregate(ent.Max(user.FieldID)).IntX(ctx)
	user := s.client.User.GetX(ctx, id)
	return toProtoUser(user)
}

// NewExtService returns a new ExtService.
func NewExtService(client *ent.Client) *ExtService {
	return &ExtService{
		client: client,
	}
}
