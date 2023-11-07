package domain_test

import (
	"context"
	"film-management/internal/user/domain"
	"film-management/internal/user/domain/mocks"
	"film-management/internal/user/domain/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockUserRepositoryBehavior func(r *mocks.MockUserRepository)
type mockAuthServiceBehavior func(r *mocks.MockAuthService)
type mockPasswordServiceBehavior func(r *mocks.MockPasswordService)

func TestService_Register(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	requireAssert := require.New(t)

	type in struct {
		user *models.User
	}

	type out struct {
		err error
	}

	tests := []struct {
		name                        string
		in                          in
		mockUserRepositoryBehavior  mockUserRepositoryBehavior
		mockPasswordServiceBehavior mockPasswordServiceBehavior
		mockAuthServiceBehavior     mockAuthServiceBehavior
		assert                      func(*in, *out)
	}{
		{
			name: "success",
			in: in{
				user: &models.User{
					Username: "test",
					Password: "12345678",
				},
			},
			mockUserRepositoryBehavior: func(r *mocks.MockUserRepository) {
				r.EXPECT().UserExistsWithUsername(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				r.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			mockPasswordServiceBehavior: func(r *mocks.MockPasswordService) {
				r.EXPECT().GeneratePasswordHash(gomock.Any()).Return("123", nil).AnyTimes()
			},
			mockAuthServiceBehavior: func(r *mocks.MockAuthService) {},
			assert: func(in *in, out *out) {
				requireAssert.NoError(out.err)
			},
		},
		{
			name: "failed to check if user exists with the same username",
			in: in{
				user: &models.User{
					Username: "test",
					Password: "12345678",
				},
			},
			mockUserRepositoryBehavior: func(r *mocks.MockUserRepository) {
				r.EXPECT().UserExistsWithUsername(gomock.Any(), gomock.Any()).Return(domain.ErrUserExistsWithUsername).AnyTimes()
			},
			mockAuthServiceBehavior:     func(r *mocks.MockAuthService) {},
			mockPasswordServiceBehavior: func(r *mocks.MockPasswordService) {},
			assert: func(in *in, out *out) {
				requireAssert.Error(out.err)
			},
		},
		{
			name: "failed to create user",
			in: in{
				user: &models.User{
					Username: "test",
					Password: "12345678",
				},
			},
			mockUserRepositoryBehavior: func(r *mocks.MockUserRepository) {
				r.EXPECT().UserExistsWithUsername(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				r.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(domain.ErrUserCreate).AnyTimes()
			},
			mockPasswordServiceBehavior: func(r *mocks.MockPasswordService) {
				r.EXPECT().GeneratePasswordHash(gomock.Any()).Return("123", nil).AnyTimes()
			},
			mockAuthServiceBehavior: func(r *mocks.MockAuthService) {},
			assert: func(in *in, out *out) {
				requireAssert.Error(out.err)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepositoryMock := mocks.NewMockUserRepository(ctrl)
			test.mockUserRepositoryBehavior(userRepositoryMock)

			authServiceMock := mocks.NewMockAuthService(ctrl)
			test.mockAuthServiceBehavior(authServiceMock)

			passwordServiceMock := mocks.NewMockPasswordService(ctrl)
			test.mockPasswordServiceBehavior(passwordServiceMock)

			userService := domain.NewService(userRepositoryMock, authServiceMock, passwordServiceMock)
			err := userService.Register(ctx, test.in.user)

			test.assert(&test.in, &out{
				err: err,
			})
		})
	}
}
