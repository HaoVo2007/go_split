package user

type UserResponse struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Name          *string `json:"name"`
	Image         *string `json:"image"`
	ImagePublicID *string `json:"image_public_id"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
}

type DashboardSummaryResponse struct {
	Balance       BalanceResponse       `json:"balance"`
	Overview      OverviewResponse      `json:"overview"`
	Expenses      ExpensesResponse      `json:"expenses"`
	TopStatistics TopStatisticsResponse `json:"top_statistics"`
}

type BalanceResponse struct {
	YouOwed float64 `json:"you_owed"`
	YouPaid float64 `json:"you_paid"`
	Balance float64 `json:"balance"`
}

type OverviewResponse struct {
	TotalGroups       int `json:"total_groups"`
	TotalTransactions int `json:"total_transactions"`
	TotalFriends      int `json:"total_friends"`
}

type ExpensesResponse struct {
	TotalPaid   float64 `json:"total_paid"`
	TotalShared float64 `json:"total_shared"`
}

type TopStatisticsResponse struct {
	TopGroup  TopGroupResponse  `json:"top_group"`
	TopFriend TopFriendResponse `json:"top_friend"`
}

type TopGroupResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	TotalSpend float64 `json:"total_spend"`
}

type TopFriendResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	TotalTransactions float64 `json:"total_transactions"`
}
