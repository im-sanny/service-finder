package repository

import (
	"fmt"
	"strings"
)

type QueryConfig struct {
	SelectCols string
	FromTable  string
}

func BuildPaginatedQuery(cfg QueryConfig, filters map[string]string, page, limit int) (string, string, []interface{}, error) {
	baseQuery := fmt.Sprintf(`SELECT %s FROM  %s`, cfg.SelectCols, cfg.FromTable)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", cfg.FromTable)

	var args []interface{}
	argIndex := 1
	whereClauses := []string{}

	// Filter by Location (Case-insensitive partial match)
	if loc, ok := filters["location"]; ok && loc != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("location ILIKE $%d", argIndex))
		args = append(args, "%"+loc+"%")
		argIndex++
	}

	if sid, ok := filters["service_id"]; ok && sid != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("service_id = $%d", argIndex))
		args = append(args, sid)
		argIndex++
	}

	if name, ok := filters["name"]; ok && name != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, name)
		argIndex++
	}

	if len(whereClauses) > 0 {
		whereSQL := " WHERE " + strings.Join(whereClauses, " AND")
		baseQuery += whereSQL
		countQuery += whereSQL
	}

	offset, safeLimit := CalculateOffset(page, limit)
	baseQuery += fmt.Sprintf(" ORDER BY id ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, safeLimit, offset)

	return baseQuery, countQuery, args, nil
}
