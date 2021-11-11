package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/imDao/db"
	"database/sql"
)

type SessionRelationParam struct {
	AccountId string
	NameLike  string
	Page      int64
	PageSize  int64
}

func GetSessionRelations(param SessionRelationParam) (relations []tables.SessionUserRelation, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.SessionUserRelation{}).Where("user_id = ?", param.AccountId)
	if param.NameLike != "" {
		template = template.Where("name like '%v'%%", param.NameLike)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Find(&relations).Order("session_id desc").Error
	return
}

func GetRelationsById(sessionId int64) (relations []tables.SessionUserRelation, err error) {
	err = db.GetDB().GetDB().Where("session_id = ?", sessionId).Find(&relations).Error
	return
}

func GetRelationsBySessionIds(sessionIds []int64, accountId string) (relations []tables.SessionUserRelation, err error) {
	err = db.GetDB().GetDB().Where("session_id in (?) and user_id = ? ", sessionIds, accountId).Find(&relations).Error
	return
}

func GetRelationById(sessionId int64, accountId string) (relation tables.SessionUserRelation, err error) {
	err = db.GetDB().GetDB().Where("session_id = ? and user_id = ?", sessionId, accountId).First(&relation).Error
	return
}

func GetSessionUserCount(sessionId int64) (count int64, err error) {
	err = db.GetDB().GetDB().Model(&tables.SessionUserRelation{}).Where("session_id = ?", sessionId).Count(&count).Error
	return
}

func GetSessionById(id int64) (session tables.Session, err error) {
	err = db.GetDB().GetDB().Where("id = ?", id).Find(&session).Error
	return
}

func GetSessionsByIds(ids []int64) (sessions []tables.Session, err error) {
	err = db.GetDB().GetDB().Where("id in (?)", ids).Find(&sessions).Error
	return
}

func GetFriendSession(accountId, friend string) (sessionId int64, err error) {
	var rows *sql.Rows
	rows, err = db.GetDB().GetDB().Raw(`select a.session_id from baby_im_session_user_rel as a 
inner join baby_im_session_user_rel as b 
on a.session_id = b.session_id
left join baby_im_session as c 
on a.session_id = c.id where a.user_id = ? and b.user_id = ? and c.session_type = 0`, accountId, friend).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&sessionId); err != nil {
			return
		}
	}
	return
}
