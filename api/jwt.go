package api

import "time"

const JwtAccessTokenLifetime = 15 * time.Minute
const JwtRefreshTokenLifetime = 24 * time.Hour
