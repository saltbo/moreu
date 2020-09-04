package model

func Tables() []interface{} {
	return []interface{}{
		new(User),
		new(UserProfile),
		new(UserInvitation),
	}
}
