package services

import "context"

func (serv *Services) GetUserIDByOrderNumber(ctx context.Context, ordNumber string) (userID int, err error) {
	userID, _ = serv.s.GetUserIDByOrder(ctx, ordNumber)
	return userID, nil
}
