package permission

type (
	LoadLevel struct {
		VersionID  string  `json:"version_id"`
		MaxLatency float64 `json:"max_latency"`
	}

	Settings struct {
		Enabled               bool        `json:"enabled"`
		PermissionResolverURL string      `json:"permission_resolver_url"`
		RequestTimeout        int         `json:"request_timeout"`
		LoadLevels            []LoadLevel `json:"load_levels"`
		DefaultVersionId      string      `json:"default_version_id"`
		CriticalVersionID     string      `json:"critical_version_id"`
		CacheLifetimeSeconds  int         `json:"cache_lifetime_seconds"`
	}
)
