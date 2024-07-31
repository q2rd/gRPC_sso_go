package auth

import (
	ssov1 "github.com/q2rd/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyIntValue    = 0
	emptyStringValue = ""
)

func basicValidation(mail, password string) error {
	if mail == emptyStringValue {
		return status.Error(codes.InvalidArgument, "Email is empty.")
	}
	if password == emptyStringValue {
		return status.Error(codes.InvalidArgument, "Password is required.")
	}
	return nil
}

func validateLoginRequest(req *ssov1.LoginRequest) error {
	if err := basicValidation(req.GetEmail(), req.GetPassword()); err != nil {
		return err
	}
	if req.GetAppId() == emptyIntValue {
		return status.Error(codes.InvalidArgument, "App id required.")
	}
	return nil
}

func validateRegisterRequest(req *ssov1.RegisterRequest) error {
	if err := basicValidation(req.GetEmail(), req.GetPassword()); err != nil {
		return err
	}
	if req.GetPassword() != req.GetPasswordConfirm() {
		return status.Error(codes.InvalidArgument, "Passwords are different. Please make sure they are the same.")
	}

	return nil
}

func validateAdminRequest(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyStringValue {
		return status.Error(codes.InvalidArgument, "User ID is required.")
	}
	return nil
}
