package controller

import (
	"net/http"

	"github.com/synic/blog/internal/middleware"
	"github.com/synic/blog/internal/store"
	"github.com/synic/blog/internal/view"
)

type LeaderboardController struct {
	views *store.PageViewRepository
}

func NewLeaderboardController(views *store.PageViewRepository) LeaderboardController {
	return LeaderboardController{views: views}
}

func (h LeaderboardController) Show(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || !user.IsAdmin {
		view.Error(w, r, nil, http.StatusForbidden, "Forbidden", "You do not have permission to view this page.")
		return
	}

	entries, err := h.views.ViewCounts(r.Context())
	if err != nil {
		view.Error(w, r, err, 500, "Internal Server Error", "An error occurred while retrieving the leaderboard.")
		return
	}

	view.Render(w, r, view.LeaderboardView(entries))
}
