CREATE ROLE animalsapp WITH LOGIN = true and PASSWORD = '1f2d1e2e67df';



CREATE KEYSPACE IF NOT EXISTS uoc_animals WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '2'};
--CREATE KEYSPACE IF NOT EXISTS uoc_animals WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

-- CREATE KEYSPACE IF NOT EXISTS media_player WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': '3'}  AND durable_writes = true AND TABLETS = {'enabled': false};
GRANT ALTER ON KEYSPACE uoc_animals TO animalsapp;
GRANT MODIFY ON KEYSPACE uoc_animals TO animalsapp;
GRANT SELECT ON KEYSPACE uoc_animals TO animalsapp;

CREATE TABLE IF NOT EXISTS uoc_animals.media (id uuid,name text,contenttype text,location text,created_at timestamp,status int,size int, PRIMARY KEY (id));

CREATE TABLE IF NOT EXISTS uoc_animals.agent (id uuid,lastname text,firstname text,type tinyint,created_at timestamp,status tinyint, PRIMARY KEY (id));

CREATE TABLE IF NOT EXISTS uoc_animals.activity (id uuid,shortcode text,description text,type text,status tinyint,created_at timestamp,updated_at timestamp, PRIMARY KEY (id));
CREATE TABLE IF NOT EXISTS uoc_animals.activity_log (id uuid,status tinyint,update_time timestamp,updated_by text,description text, PRIMARY KEY (id, update_time));


type Activity struct {
	Id          string         `json:"id"`
	ShortCode   string         `json:"shortcode" db:"shortcode"`
	Description string         `json:"description" db:"description"`
	Type        string         `json:"type"`
	Status      ActivityStatus `json:"status"`
	CreatedAt   time.Time      `json:"creationTime" db:"creation_time"`
	UpdatedAt   time.Time      `json:"updateTime" db:"update_time"`
}

type ActivityLog struct {
	Id        string         `json:"id"`
	Status    ActivityStatus `json:"status"`
	UpdatedAt time.Time      `json:"updateTime" db:"update_time"`
	UpdateBy  string         `json:"updateBy" db:"updated_by"`
}