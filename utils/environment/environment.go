package environment

import "os"

func GetPostgresDSN() string {
	return os.Getenv("POSTGRES_DSN")
}

func GetGithubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func GetRedisUri() string {
	return os.Getenv("REDIS_URI")
}

func GetWakatimeUser() string {
	return os.Getenv("WAKATIME_USER")
}

func GetWakatimeApiKey() string {
	return os.Getenv("WAKATIME_API_KEY")
}

func CheckEnvs() bool {
	if env := GetPostgresDSN(); env == "" {
		return false
	}
	if env := GetRedisUri(); env == "" {
		return false
	}
	if env := GetGithubToken(); env == "" {
		return false
	}
	if env := GetJWTSecret(); env == "" {
		return false
	}
	if env := GetWakatimeUser(); env == "" {
		return false
	}
	if env := GetWakatimeApiKey(); env == "" {
		return false
	}

	return true
}
