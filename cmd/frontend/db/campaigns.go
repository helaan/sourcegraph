package db

import (
	"context"
	"database/sql"
	"io"

	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/types"
	"github.com/sourcegraph/sourcegraph/pkg/db/dbutil"
)

// CampaignsStore exposes methods to read and write campaigns
// from persistent storage.
type CampaignsStore struct {
	db dbutil.DB
}

// NewCampaignsStore returns a new CampaignsStore backed
// by the given db.
func NewCampaignsStore(db dbutil.DB) *CampaignsStore {
	return &CampaignsStore{db: db}
}

// CreateCampaign creates the given Campaign.
func (s *CampaignsStore) CreateCampaign(ctx context.Context, c *types.Campaign) error {
	q := createCampaignQuery(c)

	rows, err := s.db.QueryContext(ctx, q.Query(sqlf.PostgresBindVar), q.Args()...)
	if err != nil {
		return err
	}

	_, _, err = scanAll(rows, func(sc scanner) (last, count int64, err error) {
		err = scanCampaign(c, sc)
		return int64(c.ID), 1, err
	})

	return err
}

var createCampaignQueryFmtstr = `
-- source: cmd/frontend/db/campaigns.go:CreateCampaign
INSERT INTO campaigns (
	name,
	description,
	author_id,
	namespace_user_id,
	namespace_org_id,
	created_at,
	updated_at
)
VALUES (%s, %s, %s, %s, %s, %s, %s)
RETURNING
	id,
	name,
	description,
	author_id,
	namespace_user_id,
	namespace_org_id,
	created_at,
	updated_at
`

func createCampaignQuery(c *types.Campaign) *sqlf.Query {
	return sqlf.Sprintf(
		createCampaignQueryFmtstr,
		c.Name,
		c.Description,
		c.AuthorID,
		nullInt32Column(c.NamespaceUserID),
		nullInt32Column(c.NamespaceOrgID),
		c.CreatedAt,
		c.UpdatedAt,
	)
}

func nullInt32Column(n int32) *int32 {
	if n == 0 {
		return nil
	}
	return &n
}

// UpdateCampaign updates the given Campaign.
func (s *CampaignsStore) UpdateCampaign(ctx context.Context, c *types.Campaign) error {
	panic("not implemented")
}

// CountCampaigns returns the number of campaigns in the database.
func (s *CampaignsStore) CountCampaigns(ctx context.Context) (int64, error) {
	q := countCampaignsQuery

	rows, err := s.db.QueryContext(ctx, q.Query(sqlf.PostgresBindVar), q.Args()...)
	if err != nil {
		return 0, err
	}

	_, count, err := scanAll(rows, func(sc scanner) (_, count int64, err error) {
		err = sc.Scan(&count)
		return 0, count, err
	})

	return count, err
}

var countCampaignsQuery = sqlf.Sprintf(`
-- source: cmd/frontend/db/campaigns.go:CountCampaigns
SELECT COUNT(id) FROM campaigns
`)

// ListCampaignsOpts captures the query options needed for
// listing campaigns.
type ListCampaignsOpts struct {
	Limit int
}

// ListCampaigns lists Campaigns with the given filters.
func (s *CampaignsStore) ListCampaigns(ctx context.Context, opts ListCampaignsOpts) (cs []*types.Campaign, next int64, err error) {
	q := listCampaignsQuery(&opts)

	rows, err := s.db.QueryContext(ctx, q.Query(sqlf.PostgresBindVar), q.Args()...)
	if err != nil {
		return nil, 0, err
	}

	cs = make([]*types.Campaign, 0, opts.Limit)
	_, _, err = scanAll(rows, func(sc scanner) (last, count int64, err error) {
		var c types.Campaign
		if err = scanCampaign(&c, sc); err != nil {
			return 0, 0, err
		}
		cs = append(cs, &c)
		return int64(c.ID), 1, err
	})

	if len(cs) == opts.Limit {
		next = cs[len(cs)-1].ID
		cs = cs[:len(cs)-1]
	}

	return cs, next, err
}

var listCampaignsQueryFmtstr = `
-- source: cmd/frontend/db/campaigns.go:ListCampaigns
SELECT
	id,
	name,
	description,
	author_id,
	namespace_user_id,
	namespace_org_id,
	created_at,
	updated_at
FROM campaigns
ORDER BY id ASC
LIMIT %s
`

const defaultListLimit = 50

func listCampaignsQuery(opts *ListCampaignsOpts) *sqlf.Query {
	if opts.Limit == 0 {
		opts.Limit = defaultListLimit
	}
	opts.Limit++
	return sqlf.Sprintf(listCampaignsQueryFmtstr, opts.Limit)
}

// scanner captures the Scan method of sql.Rows and sql.Row
type scanner interface {
	Scan(dst ...interface{}) error
}

// a scanFunc scans one or more rows from a scanner, returning
// the last id column scanned and the count of scanned rows.
type scanFunc func(scanner) (last, count int64, err error)

func scanAll(rows *sql.Rows, scan scanFunc) (last, count int64, err error) {
	defer closeErr(rows, &err)

	last = -1
	for rows.Next() {
		var n int64
		if last, n, err = scan(rows); err != nil {
			return last, count, err
		}
		count += n
	}

	return last, count, rows.Err()
}

func closeErr(c io.Closer, err *error) {
	if e := c.Close(); err != nil && *err == nil {
		*err = e
	}
}

func scanCampaign(c *types.Campaign, s scanner) error {
	return s.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.AuthorID,
		&dbutil.NullInt32{N: &c.NamespaceUserID},
		&dbutil.NullInt32{N: &c.NamespaceOrgID},
		&c.CreatedAt,
		&c.UpdatedAt,
	)
}
