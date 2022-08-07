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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	metadata := s.extractMetadata(ctx)
	arg := service.LoginUserReq{
		Req: service.Req{
			ClientIP:  metadata.ClientIP,
			UserAgent: metadata.UserAgent,
		},
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	res := s.Service.LoginUser(ctx, &arg)
	if res.Error.Value != nil {
		switch res.Error.Kind {
		case utils.KindErrNotFound:
			return nil, status.Errorf(codes.NotFound, "user not found")
		case utils.KindErrUnauthorized:
			return nil, status.Errorf(codes.Unauthenticated, "user not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to find user")
		}
	}

	respData := res.Data.(*service.LoginUserRes)

	rsp := &pb.LoginUserResponse{
		User: convertUser(&db.User{
			Username:          respData.User.Username,
			FullName:          respData.User.FullName,
			Email:             respData.User.Email,
			PasswordChangedAt: respData.User.PasswordChangedAt,
			CreatedAt:         respData.User.CreatedAt,
		}),
		SessionId:             respData.SessionID.String(),
		AccessToken:           respData.AccessToken,
		RefreshToken:          respData.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(respData.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(respData.RefreshTokenExpiresAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := utils.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
