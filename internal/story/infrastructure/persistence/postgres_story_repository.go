package persistence

import (
	"context"
	"encoding/json"

	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/story/domain"
)

type PostgresStoryRepository struct {
	q postgres.Querier
}

func NewPostgresStoryRepository(q postgres.Querier) *PostgresStoryRepository {
	return &PostgresStoryRepository{q: q}
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// ---------- Series ----------

func (r *PostgresStoryRepository) SaveSeries(ctx context.Context, s *domain.Series) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.series (id, project_id, title, description)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description`,
		s.ID, s.ProjectID, s.Title, s.Description)
	return err
}

func (r *PostgresStoryRepository) FindSeriesByID(ctx context.Context, id string) (*domain.Series, error) {
	row := r.q.QueryRow(ctx, `SELECT id, project_id, title, description FROM story.series WHERE id = $1`, id)
	return scanSeries(row)
}

func (r *PostgresStoryRepository) FindSeriesByProject(ctx context.Context, projectID string) (*domain.Series, error) {
	row := r.q.QueryRow(ctx, `SELECT id, project_id, title, description FROM story.series WHERE project_id = $1 LIMIT 1`, projectID)
	return scanSeries(row)
}

func scanSeries(row rowScanner) (*domain.Series, error) {
	var s domain.Series
	if err := row.Scan(&s.ID, &s.ProjectID, &s.Title, &s.Description); err != nil {
		if postgres.IsNoRows(err) {
			return nil, domain.ErrSeriesNotFound
		}
		return nil, err
	}
	return &s, nil
}

// ---------- Seasons ----------

func (r *PostgresStoryRepository) SaveSeason(ctx context.Context, s *domain.Season) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.seasons (id, series_id, number, title, summary)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (id) DO UPDATE SET number = EXCLUDED.number, title = EXCLUDED.title, summary = EXCLUDED.summary`,
		s.ID, s.SeriesID, s.Number, s.Title, s.Summary)
	return err
}

func (r *PostgresStoryRepository) FindSeasonByID(ctx context.Context, id string) (*domain.Season, error) {
	row := r.q.QueryRow(ctx, `SELECT id, series_id, number, title, summary FROM story.seasons WHERE id = $1`, id)
	var s domain.Season
	if err := row.Scan(&s.ID, &s.SeriesID, &s.Number, &s.Title, &s.Summary); err != nil {
		if postgres.IsNoRows(err) {
			return nil, domain.ErrSeasonNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *PostgresStoryRepository) ListSeasonsBySeries(ctx context.Context, seriesID string) ([]*domain.Season, error) {
	rows, err := r.q.Query(ctx, `SELECT id, series_id, number, title, summary FROM story.seasons WHERE series_id = $1 ORDER BY number`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Season
	for rows.Next() {
		var s domain.Season
		if err := rows.Scan(&s.ID, &s.SeriesID, &s.Number, &s.Title, &s.Summary); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}

// ---------- Arcs ----------

func (r *PostgresStoryRepository) SaveArc(ctx context.Context, a *domain.StoryArc) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.arcs (id, season_id, title, number) VALUES ($1,$2,$3,$4)
		ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, number = EXCLUDED.number`,
		a.ID, a.SeasonID, a.Title, a.Number)
	return err
}

func (r *PostgresStoryRepository) FindArcByID(ctx context.Context, id string) (*domain.StoryArc, error) {
	var a domain.StoryArc
	err := r.q.QueryRow(ctx, `SELECT id, season_id, title, number FROM story.arcs WHERE id = $1`, id).
		Scan(&a.ID, &a.SeasonID, &a.Title, &a.Number)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrArcNotFound
	}
	return &a, err
}

func (r *PostgresStoryRepository) ListArcsBySeason(ctx context.Context, seasonID string) ([]*domain.StoryArc, error) {
	rows, err := r.q.Query(ctx, `SELECT id, season_id, title, number FROM story.arcs WHERE season_id = $1 ORDER BY number`, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StoryArc
	for rows.Next() {
		var a domain.StoryArc
		if err := rows.Scan(&a.ID, &a.SeasonID, &a.Title, &a.Number); err != nil {
			return nil, err
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}

// ---------- Episodes ----------

func (r *PostgresStoryRepository) SaveEpisode(ctx context.Context, e *domain.Episode) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.episodes (id, season_id, arc_id, number, title, summary, script, status, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8, now())
		ON CONFLICT (id) DO UPDATE SET arc_id = EXCLUDED.arc_id, number = EXCLUDED.number,
			title = EXCLUDED.title, summary = EXCLUDED.summary, script = EXCLUDED.script,
			status = EXCLUDED.status, updated_at = now()`,
		e.ID, e.SeasonID, e.ArcID, e.Number, e.Title, e.Summary, e.Script, e.Status)
	return err
}

func (r *PostgresStoryRepository) FindEpisodeByID(ctx context.Context, id string) (*domain.Episode, error) {
	var e domain.Episode
	err := r.q.QueryRow(ctx, `SELECT id, season_id, arc_id, number, title, summary, script, status FROM story.episodes WHERE id = $1`, id).
		Scan(&e.ID, &e.SeasonID, &e.ArcID, &e.Number, &e.Title, &e.Summary, &e.Script, &e.Status)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrEpisodeNotFound
	}
	return &e, err
}

func (r *PostgresStoryRepository) ListEpisodesBySeason(ctx context.Context, seasonID string) ([]*domain.Episode, error) {
	rows, err := r.q.Query(ctx, `SELECT id, season_id, arc_id, number, title, summary, script, status FROM story.episodes WHERE season_id = $1 ORDER BY number`, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Episode
	for rows.Next() {
		var e domain.Episode
		if err := rows.Scan(&e.ID, &e.SeasonID, &e.ArcID, &e.Number, &e.Title, &e.Summary, &e.Script, &e.Status); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

// ---------- Scenes ----------

func (r *PostgresStoryRepository) SaveScene(ctx context.Context, s *domain.Scene) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.scenes (id, episode_id, number, title, location_id, time_of_day, description, character_ids)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (id) DO UPDATE SET number = EXCLUDED.number, title = EXCLUDED.title,
			location_id = EXCLUDED.location_id, time_of_day = EXCLUDED.time_of_day,
			description = EXCLUDED.description, character_ids = EXCLUDED.character_ids`,
		s.ID, s.EpisodeID, s.Number, s.Title, s.LocationID, s.TimeOfDay, s.Description, mustJSON(s.CharacterIDs))
	return err
}

func (r *PostgresStoryRepository) FindSceneByID(ctx context.Context, id string) (*domain.Scene, error) {
	row := r.q.QueryRow(ctx, `SELECT id, episode_id, number, title, location_id, time_of_day, description, character_ids FROM story.scenes WHERE id = $1`, id)
	return scanScene(row)
}

func (r *PostgresStoryRepository) ListScenesByEpisode(ctx context.Context, episodeID string) ([]*domain.Scene, error) {
	rows, err := r.q.Query(ctx, `SELECT id, episode_id, number, title, location_id, time_of_day, description, character_ids FROM story.scenes WHERE episode_id = $1 ORDER BY number`, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Scene
	for rows.Next() {
		s, err := scanScene(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanScene(row rowScanner) (*domain.Scene, error) {
	var s domain.Scene
	var chars []byte
	err := row.Scan(&s.ID, &s.EpisodeID, &s.Number, &s.Title, &s.LocationID, &s.TimeOfDay, &s.Description, &chars)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrSceneNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(chars, &s.CharacterIDs)
	return &s, nil
}

// ---------- Beats ----------

func (r *PostgresStoryRepository) SaveBeat(ctx context.Context, b *domain.Beat) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.beats (id, scene_id, seq, action, dialogue, character_id)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET seq = EXCLUDED.seq, action = EXCLUDED.action,
			dialogue = EXCLUDED.dialogue, character_id = EXCLUDED.character_id`,
		b.ID, b.SceneID, b.Seq, b.Action, b.Dialogue, b.CharacterID)
	return err
}

func (r *PostgresStoryRepository) FindBeatByID(ctx context.Context, id string) (*domain.Beat, error) {
	var b domain.Beat
	err := r.q.QueryRow(ctx, `SELECT id, scene_id, seq, action, dialogue, character_id FROM story.beats WHERE id = $1`, id).
		Scan(&b.ID, &b.SceneID, &b.Seq, &b.Action, &b.Dialogue, &b.CharacterID)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrBeatNotFound
	}
	return &b, err
}

func (r *PostgresStoryRepository) ListBeatsByScene(ctx context.Context, sceneID string) ([]*domain.Beat, error) {
	rows, err := r.q.Query(ctx, `SELECT id, scene_id, seq, action, dialogue, character_id FROM story.beats WHERE scene_id = $1 ORDER BY seq`, sceneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Beat
	for rows.Next() {
		var b domain.Beat
		if err := rows.Scan(&b.ID, &b.SceneID, &b.Seq, &b.Action, &b.Dialogue, &b.CharacterID); err != nil {
			return nil, err
		}
		out = append(out, &b)
	}
	return out, rows.Err()
}

// ---------- Story graph ----------

func (r *PostgresStoryRepository) SaveGraphNode(ctx context.Context, projectID string, n *domain.StoryNode) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.graph_nodes (id, project_id, type, title, episode_id, scene_id)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, title = EXCLUDED.title,
			episode_id = EXCLUDED.episode_id, scene_id = EXCLUDED.scene_id`,
		n.ID, projectID, n.Type, n.Title, n.EpisodeID, n.SceneID)
	return err
}

func (r *PostgresStoryRepository) FindGraphNode(ctx context.Context, id string) (*domain.StoryNode, error) {
	var n domain.StoryNode
	err := r.q.QueryRow(ctx, `SELECT id, type, title, episode_id, scene_id FROM story.graph_nodes WHERE id = $1`, id).
		Scan(&n.ID, &n.Type, &n.Title, &n.EpisodeID, &n.SceneID)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrNodeNotFound
	}
	return &n, err
}

func (r *PostgresStoryRepository) ListGraphNodes(ctx context.Context, projectID string) ([]*domain.StoryNode, error) {
	rows, err := r.q.Query(ctx, `SELECT id, type, title, episode_id, scene_id FROM story.graph_nodes WHERE project_id = $1 ORDER BY id`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StoryNode
	for rows.Next() {
		var n domain.StoryNode
		if err := rows.Scan(&n.ID, &n.Type, &n.Title, &n.EpisodeID, &n.SceneID); err != nil {
			return nil, err
		}
		out = append(out, &n)
	}
	return out, rows.Err()
}

func (r *PostgresStoryRepository) SaveGraphEdge(ctx context.Context, projectID string, e *domain.StoryEdge) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.graph_edges (project_id, from_id, to_id, relation)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (project_id, from_id, to_id, relation) DO NOTHING`,
		projectID, e.FromID, e.ToID, e.Relation)
	return err
}

func (r *PostgresStoryRepository) ListGraphEdges(ctx context.Context, projectID string) ([]*domain.StoryEdge, error) {
	rows, err := r.q.Query(ctx, `SELECT from_id, to_id, relation FROM story.graph_edges WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StoryEdge
	for rows.Next() {
		var e domain.StoryEdge
		if err := rows.Scan(&e.FromID, &e.ToID, &e.Relation); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

// ---------- Plot threads ----------

func (r *PostgresStoryRepository) SavePlotThread(ctx context.Context, t *domain.PlotThread) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO story.plot_threads (id, project_id, name, description, status)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, status = EXCLUDED.status`,
		t.ID, t.ProjectID, t.Name, t.Description, t.Status)
	return err
}

func (r *PostgresStoryRepository) ListPlotThreads(ctx context.Context, projectID string) ([]*domain.PlotThread, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, name, description, status FROM story.plot_threads WHERE project_id = $1 ORDER BY id`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.PlotThread
	for rows.Next() {
		var t domain.PlotThread
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Name, &t.Description, &t.Status); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}
