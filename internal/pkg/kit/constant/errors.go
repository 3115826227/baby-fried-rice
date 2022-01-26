package constant

import "errors"

var (
	QueryUserByIdDaoError = errors.New("user dao by id query failed")

	NeedApplyAddFriendError = errors.New("you need apply add friend")

	NeedInviteJoinSessionError = errors.New("you need the session's origin invites you to join the session")

	UnOriginInviteJoinSessionError = errors.New("you are not the session's origin, can't invite user join the session")

	UnOriginInviteRemoveSessionError = errors.New("you are not the session's origin, can't make user remove the session")
)
