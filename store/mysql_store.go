package store

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // ✅ 必须匿名导入 MySQL 驱动
)

type MySQLStore struct {
	db *sql.DB
}

func InitMySQL() *sql.DB {
	dsn := "root:La+021026@tcp(192.168.31.151:3306)/shorturl?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
func NewMySQLStore(db *sql.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

func (m *MySQLStore) Save(shortUrl, longUrl, userId, customAlias string) error {
	_, err := m.db.Exec("INSERT INTO short_links (short_url, long_url, user_id, custom_alias) VALUES (?, ?, ?, ?)",
		shortUrl, longUrl, userId, customAlias)
	return err
}

func (m *MySQLStore) Get(shortUrl string) (string, error) {
	var longUrl string
	err := m.db.QueryRow("SELECT long_url FROM short_links WHERE short_url = ?", shortUrl).Scan(&longUrl)
	return longUrl, err
}

func (m *MySQLStore) Exists(shortUrl string) bool {
	var exists bool
	err := m.db.QueryRow("SELECT EXISTS(SELECT 1 FROM short_links WHERE short_url = ?)", shortUrl).Scan(&exists)
	return err == nil && exists
}

func (m *MySQLStore) IncrementVisitCount(shortUrl string) {
	_, _ = m.db.Exec("UPDATE short_links SET click_count = click_count + 1 WHERE short_url = ?", shortUrl)
}
