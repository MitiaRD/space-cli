package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBuildLaunchQuery(t *testing.T) {
	tests := []struct {
		name     string
		flags    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:  "default query",
			flags: map[string]interface{}{},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with date range",
			flags: map[string]interface{}{
				"start": "2024-01-01",
				"end":   "2024-01-31",
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
					"date_utc": map[string]interface{}{
						"$gte": "2024-01-01T00:00:00.000Z",
						"$lte": "2024-01-31T23:59:59.999Z",
					},
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with failed flag",
			flags: map[string]interface{}{
				"failed": true,
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
					"success":  false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with upcoming flag",
			flags: map[string]interface{}{
				"upcoming": true,
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": true,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with limit",
			flags: map[string]interface{}{
				"limit": 10,
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
					"limit": 10,
				},
			},
		},
		{
			name: "with all flags",
			flags: map[string]interface{}{
				"start":    "2024-01-01",
				"end":      "2024-01-31",
				"failed":   true,
				"upcoming": false,
				"limit":    5,
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
					"success":  false,
					"date_utc": map[string]interface{}{
						"$gte": "2024-01-01T00:00:00.000Z",
						"$lte": "2024-01-31T23:59:59.999Z",
					},
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
					"limit": 5,
				},
			},
		},
		{
			name: "with zero limit (should not add limit)",
			flags: map[string]interface{}{
				"limit": 0,
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with negative limit (should not add limit)",
			flags: map[string]interface{}{
				"limit": -1,
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with only start date (should not add date query)",
			flags: map[string]interface{}{
				"start": "2024-01-01",
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
		{
			name: "with only end date (should not add date query)",
			flags: map[string]interface{}{
				"end": "2024-01-31",
			},
			expected: map[string]interface{}{
				"query": map[string]interface{}{
					"upcoming": false,
				},
				"options": map[string]interface{}{
					"sort": map[string]interface{}{
						"date_utc": "desc",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			cmd.Flags().String("start", "", "Start date")
			cmd.Flags().String("end", "", "End date")
			cmd.Flags().Bool("failed", false, "Show failed launches")
			cmd.Flags().Bool("upcoming", false, "Show upcoming launches")
			cmd.Flags().Int("limit", 0, "Limit number of results")

			for flag, value := range tt.flags {
				switch v := value.(type) {
				case string:
					cmd.Flags().Set(flag, v)
				case bool:
					if v {
						cmd.Flags().Set(flag, "true")
					}
				case int:
					cmd.Flags().Set(flag, fmt.Sprintf("%d", v))
				}
			}

			result := buildLaunchQuery(cmd)

			assert.Equal(t, tt.expected, result, "Query should match expected structure")
		})
	}
}
