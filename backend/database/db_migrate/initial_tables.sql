
-- store users
CREATE TABLE users (
    discord_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    display_name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- store prof players
CREATE TABLE players (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    external_id TEXT UNIQUE,
    name TEXT NOT NULL,
    country TEXT,
    country_code VARCHAR(2),
    qualification TEXT,
    route TEXT,
    bio TEXT,
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- leagues created by users
CREATE TABLE leagues (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    owner_id TEXT NOT NULL REFERENCES users(discord_id),
    max_members INTEGER NOT NULL CHECK (max_members > 0),
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    invite_code TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- users belonging to leagues
CREATE TABLE league_members (
    league_id BIGINT NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(discord_id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (league_id, user_id)
);

-- fantasy teams created by users within leagues
CREATE TABLE fantasy_teams (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    owner_id TEXT NOT NULL REFERENCES users(discord_id),
    league_id BIGINT NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (owner_id, league_id)
);

-- tracks players drafted by fantasy teams
CREATE TABLE roster_entries(
    team_id BIGINT NOT NULL REFERENCES fantasy_teams(id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL REFERENCES players(id),
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (team_id, player_id)
);

-- record match results
CREATE TABLE matches (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    external_id TEXT UNIQUE,
    tournament_name TEXT NOT NULL,
    match_date TIMESTAMPTZ, status TEXT NOT NULL DEFAULT 'scheduled'
        CHECK (status IN ('scheduled', 'ongoing', 'completed')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- record player performance in matches
CREATE TABLE player_match_stats (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL REFERENCES players(id),
    result TEXT NOT NULL CHECK (result IN ('win', 'loss', 'draw')),
    towers_destroyed INTEGER NOT NULL DEFAULT 0
        CHECK (towers_destroyed >= 0),
    towers_lost INTEGER NOT NULL DEFAULT 0
        CHECK (towers_lost >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (match_id, player_id)
);

-- store fantasy points for players
CREATE TABLE player_fantasy_scores (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    match_id BIGINT NOT NULL
        REFERENCES matches(id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL
        REFERENCES players(id),
    fantasy_points NUMERIC(10,2) NOT NULL DEFAULT 0,
    calculated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (match_id, player_id),
    FOREIGN KEY (match_id, player_id)
        REFERENCES player_match_stats(match_id, player_id)
);

-- store fantasy points for fantasy teams
CREATE TABLE team_fantasy_scores (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    match_id BIGINT NOT NULL 
        REFERENCES matches(id) ON DELETE CASCADE,
    team_id BIGINT NOT NULL
        REFERENCES fantasy_teams(id) ON DELETE CASCADE,
    fantasy_points NUMERIC(10,2) NOT NULL DEFAULT 0,
    calculated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (match_id, team_id)
);



