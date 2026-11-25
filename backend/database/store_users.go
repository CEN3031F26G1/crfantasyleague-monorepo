package database

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
)

//inserts or updates a Discord user.
func SaveDiscordUser(
    ctx context.Context,
    db *pgxpool.Pool,
    discordID string,
    username string,
    displayName string,
    avatarURL string,
) error {
    query := `
        INSERT INTO users (
            discord_id,
            username,
            display_name,
            avatar_url
        )
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (discord_id)
        DO UPDATE SET
            username = EXCLUDED.username,
            display_name = EXCLUDED.display_name,
            avatar_url = EXCLUDED.avatar_url,
            updated_at = CURRENT_TIMESTAMP
    `

    _, err := db.Exec(
        ctx,
        query,
        discordID,
        username,
        displayName,
        avatarURL,
    )

    return err
}
