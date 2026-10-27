package http

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/transport/http/handlers"
	appMiddleware "bn-riftstats/internal/transport/http/middleware"
)

func NewRouter(cfg *config.Config, h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Base Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Riot-Token"},
		ExposedHeaders:   []string{"Link", "Retry-After", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth Parser
	r.Use(appMiddleware.AuthMiddleware(cfg.JWTSecret))

	// Static Data & Legal Attribution Header middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-RiftStats-Patch", cfg.CurrentPatch)
			w.Header().Set("X-RiftStats-DDragon", cfg.DDragonVersion)
			next.ServeHTTP(w, r)
		})
	})

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		// 01_Regions
		r.Get("/regions", h.ListRegions)

		// 02_AccountLookup & Search
		r.Get("/accounts/by-riot-id/{region}/{gameName}/{tagLine}", h.ResolveRiotID)
		r.Get("/search/autocomplete", h.Autocomplete)

		// 03_SummonerProfile, 04_RankedStats, 05_MatchHistory
		r.Get("/summoners/{region}/by-puuid/{puuid}", h.GetSummonerProfile)
		r.Get("/summoners/{region}/by-puuid/{puuid}/ranked", h.GetRankedStats)
		r.Get("/summoners/{region}/by-puuid/{puuid}/matches", h.GetMatchHistory)
		r.Get("/summoners/{region}/by-puuid/{puuid}/mastery", h.GetSummonerMastery)

		// 06_MatchDetail & 07_MatchTimeline
		r.Get("/matches/{region}/{matchId}", h.GetMatchDetail)
		r.Get("/matches/{region}/{matchId}/timeline", h.GetMatchTimeline)

		// 08_ChampionCatalog & 09_ChampionStats
		r.Get("/champions", h.ListChampions)
		r.Get("/champions/tier-list", h.GetTierList)

		// 10_ChampionBuilds, 11_Runes, 12_Items, 13_Counters, 14_Synergies
		r.Get("/champions/{championKey}/builds", h.GetChampionBuilds)
		r.Get("/champions/{championKey}/runes", h.GetChampionRunes)
		r.Get("/champions/{championKey}/items", h.GetChampionItems)
		r.Get("/champions/{championKey}/counters", h.GetChampionCounters)
		r.Get("/champions/{championKey}/synergies", h.GetChampionSynergies)

		// 15_Leaderboards
		r.Get("/leaderboards/{region}", h.GetLeaderboards)

		// 17_GameModes & 18_Patches
		r.Get("/modes/{modeSlug}/stats", h.GetModeStats)
		r.Get("/patches", h.ListPatches)

		// 19_Spectator
		r.Get("/spectator/{region}/by-puuid/{puuid}", h.GetSpectatorGame)

		// 20_Users & 21_Favorites & 22_Auth
		r.Get("/users/me", h.GetCurrentUser)
		r.Get("/favorites", h.ListFavorites)
		r.Post("/favorites", h.CreateFavorite)
		r.Post("/auth/login", h.Login)

		// 23_Admin
		r.Get("/admin/overview", h.GetAdminOverview)

		// 24_DataIngestion & 25_Analytics
		r.Post("/ingestion/summoners/{region}/{puuid}/refresh", h.TriggerSummonerRefresh)
		r.Get("/analytics/summary", h.GetAnalyticsSummary)

		// 26_Monitoring
		r.Get("/monitoring/health", h.GetSystemHealth)
	})

	// Static SPA file server
	workDir, _ := os.Getwd()
	distDir := filepath.Join(workDir, "web", "dist")
	if info, err := os.Stat(distDir); err == nil && info.IsDir() {
		fileServer := http.FileServer(http.Dir(distDir))
		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			path := filepath.Join(distDir, req.URL.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// Fallback to index.html for SPA routes
				http.ServeFile(w, req, filepath.Join(distDir, "index.html"))
				return
			}
			fileServer.ServeHTTP(w, req)
		})
	}

	return r
}
