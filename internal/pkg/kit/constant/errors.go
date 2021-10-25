package constant

import "errors"

var (
	QueryUserByIdDaoError = errors.New("user dao by id query failed")

	NeedApplyAddFriendError = errors.New("you need apply add friend")

	NeedInviteJoinSessionError = errors.New("you need the session's origin invites you to join the session")
)

var (
	CodeNeedOriginAuditSessionMsg = "waiting the session's origin audit"
)
