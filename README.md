# Personal AI

## install

```bash
git clone
cd personal-ai
# create empty token.json
touch token.json

Add your API keys to keys file
echo """
export const WOLFRAM_APP_ID="XXXXXX"
export const WIT_AI_ACCESS_TOKEN="XXXXXXXX"

export const SPOTIFY_ID="XXXXXXXXXX"
export const SPOTIFY_SECRET="XXXXXXXXXX"

export const TOKEN_FILE="token.json"

export const Target_Playlist="fav-playlist"

export const SPOTIFY_REDIRECT_URI="http://localhost:8000/spotify/callback"
""" > keys
source keys
go run .
```
