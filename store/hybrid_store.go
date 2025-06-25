// store/hybrid_store.go

package store

type HybridStore struct {
	cache Store // RedisStore
	db    Store // MySQLStore
}

func NewHybridStore(cache Store, db Store) *HybridStore {
	return &HybridStore{
		cache: cache,
		db:    db,
	}
}

func (h *HybridStore) Save(shortUrl, longUrl, userId, customAlias string) error {
	// 先存数据库（保证持久化），再存缓存（保证加速访问）
	if err := h.db.Save(shortUrl, longUrl, userId, customAlias); err != nil {
		return err
	}
	_ = h.cache.Save(shortUrl, longUrl, userId, customAlias) // 缓存失败可忽略
	return nil
}

func (h *HybridStore) Get(shortUrl string) (string, error) {
	// 优先查 Redis（缓存）
	longUrl, err := h.cache.Get(shortUrl)
	if err == nil {
		return longUrl, nil
	}

	// Redis 未命中，查数据库
	longUrl, err = h.db.Get(shortUrl)
	if err != nil {
		return "", err
	}

	// 回写 Redis，加速下次访问
	_ = h.cache.Save(shortUrl, longUrl, "", "") // userId 和 customAlias 可忽略
	return longUrl, nil
}

func (h *HybridStore) Exists(shortUrl string) bool {
	// 只查数据库是否存在（缓存一般没这个功能或不可靠）
	return h.db.Exists(shortUrl)
}

// hybrid_store.go

func (h *HybridStore) IncrementVisitCount(shortUrl string) {
	if tracker, ok := h.db.(VisitTracker); ok {
		tracker.IncrementVisitCount(shortUrl)
	}
}
