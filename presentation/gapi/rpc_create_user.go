package gapi

import (
	"context"
	db "ginbank/db/sqlc"
	"ginbank/pb"
	"ginbank/service"
	"ginbank/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := service.CreateUserReq{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}

	res := s.Service.CreateUser(ctx, &arg)
	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrUniqueViolation:
			return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", res.Error.Value)
		default:
			return nil, status.Errorf(codes.Internal, "failed to create user: %s", res.Error.Value)
		}
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(res.Data.(*db.User)),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := utils.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := utils.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := utils.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
