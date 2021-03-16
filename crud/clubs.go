package crud

import (
	"github.com/Krishap-s/keats-backend/models"
	"github.com/Krishap-s/keats-backend/pgdb"
	"github.com/Krishap-s/keats-backend/schemas"
	"github.com/google/uuid"
)

// CreateUser creates a club in the database or returns an error
func CreateClub(objIn *schemas.ClubCreate) (*models.Club, error) {
	db := pgdb.GetDB()
	uid, err := uuid.Parse(objIn.HostID)
	if err != nil {
		return nil, err
	}
	club := &models.Club{
		ClubName: objIn.ClubName,
		PageSync: objIn.PageSync,
		FileURL:  objIn.FileURL,
		Private:  objIn.Private,
		PageNo:   objIn.PageNo,
		HostID:   uid,
	}

	_, err = db.Model(club).
		Returning("*").
		Insert()
	if err != nil {
		return nil, err
	}

	clubuser := &models.ClubUser{
		ClubID: club.ID,
		UserID: uid,
	}
	_, err = db.Model(clubuser).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return club, nil
}

// UpdateClub updates a club in the database or returns an error
func UpdateClub(objIn *schemas.ClubUpdate) (*models.Club, error) {
	db := pgdb.GetDB()
	uid, err := uuid.Parse(objIn.ID)
	if err != nil {
		return nil, err
	}
	club := &models.Club{
		ID:       uid,
		ClubName: objIn.ClubName,
		FileURL:  objIn.FileURL,
		Private:  objIn.Private,
		PageSync: objIn.PageSync,
	}

	_, err = db.Model(club).Returning("*").WherePK().UpdateNotZero()
	if err != nil {
		return nil, err
	}

	return club, nil
}

// GetClub gets a club from the database or returns an error
func GetClub(id string) (*models.Club, error) {
	db := pgdb.GetDB()
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	club := &models.Club{
		ID: cid,
	}
	err = db.Model(club).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return club, nil
}

// CreateClubUser creates a clubuser record in the database
func CreateClubUser(ClubId string, UserId string) (*models.ClubUser, error) {
	db := pgdb.GetDB()
	cid, err := uuid.Parse(ClubId)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(UserId)
	if err != nil {
		return nil, err
	}
	clubuser := &models.ClubUser{
		ClubID: cid,
		UserID: uid,
	}
	_, err = db.Model(clubuser).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return clubuser, err
}

// GetClubUser get clubuser records from database
func GetClubUser(ClubId string) ([]*models.User, error) {
	db := pgdb.GetDB()
	cid, err := uuid.Parse(ClubId)
	if err != nil {
		return nil, err
	}
	var users []*models.User
	err = db.Model(&users).
		ColumnExpr("\"user\".\"id\" , \"user\".\"username\", \"user\".\"profile_pic\"").
		Join("INNER JOIN club_users as cu").
		JoinOn("cu.user_id = \"user\".\"id\"").
		Where("cu.club_id = ?", cid).
		Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteClubUser deletes clubuser record from database
func DeleteClubUser(ClubId string, UserId string) (*models.ClubUser, error) {
	db := pgdb.GetDB()
	cid, err := uuid.Parse(ClubId)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(UserId)
	if err != nil {
		return nil, err
	}
	clubuser := &models.ClubUser{
		ClubID: cid,
		UserID: uid,
	}
	_, err = db.Model(clubuser).Where("user_id = ?user_id and club_id = ?club_id").Delete()
	if err != nil {
		return nil, err
	}
	var users []*models.User
	err = db.Model(&users).
		ColumnExpr("\"user\".\"id\" , \"user\".\"username\", \"user\".\"profile_pic\"").
		Join("INNER JOIN club_users as cu").
		JoinOn("cu.user_id = \"user\".\"id\"").
		Where("cu.club_id = ?", cid).
		Select()
	if err != nil {
		return nil, err
	}
	// Reset Host ID to someone else if host themselves is leaving
	club := &models.Club{
		ID: cid,
	}
	err = db.Model(club).WherePK().Select()
	if err != nil {
		return nil, err
	}
	if club.HostID == uid {
		if len(users) != 0 {
			club.HostID = users[0].ID
		} else {
			club.HostID = uuid.Nil
		}
		_, err = db.Model(club).WherePK().Update()
		if err != nil {
			return nil, err
		}
	}
	return clubuser, err
}