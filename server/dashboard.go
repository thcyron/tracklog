package server

import (
	"net/http"

	"github.com/thcyron/tracklog"
)

const dashboardRecentLogsCount = 5

type dashboardData struct {
	Logs []dashboardLog
}

type dashboardLog struct {
	ID       int
	Name     string
	Date     string
	Duration string
	Distance string
	Tags     []string
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)

	user := ctx.User()
	if user == nil {
		s.redirectToSignIn(w, r)
		return
	}

	logs, err := s.db.RecentUserLogs(user, dashboardRecentLogsCount)
	if err != nil {
		panic(err)
	}

	data := new(dashboardData)
	for _, log := range logs {
		data.Logs = append(data.Logs, dashboardLog{
			ID:       log.ID,
			Name:     log.Name,
			Date:     log.Start.Format(logTimeFormat),
			Duration: tracklog.Duration(log.Duration).String(),
			Distance: tracklog.Distance(log.Distance).String(),
			Tags:     log.Tags,
		})
	}

	ctx.SetTitle("Dashboard")
	ctx.Breadcrumb().Add("Dashboard", "", true)
	ctx.SetActiveTab("dashboard")
	ctx.SetData(data)

	s.render(w, r, "dashboard")
}
