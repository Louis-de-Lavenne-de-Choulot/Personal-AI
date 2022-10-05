package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:5019/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState),
	)
	ch    = make(chan *spotify.Client)
	state = "abc123"

	// We'll want these variables sooner rather than later
	client      *spotify.Client
	playerState *spotify.PlayerState
)

func authSpotify() {
	//if safe.txt contains a token, unmarshal it and use it to get an authenticated client
	f, err := os.Open("safe.txt")
	if err != nil {
		println("No token found, please authenticate")
	}
	defer f.Close()
	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	if err != nil {
		println("No token found, please authenticate")

		go func() {
			//handle /callback should call completeAuth then redirect to the main page /
			http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
				completeAuth(w, r)
				http.Redirect(w, r, "/", http.StatusFound)
			})
			url := auth.AuthURL(state)
			println("Please log in to Spotify by visiting the following page in your browser:"+url, "en")

			// wait for auth to complete
			client = <-ch

			// use the client to make calls that require authorization
			user, err := client.CurrentUser(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("You are logged in as:", user.ID)
			// //save client in safe.txt
			// f, err := os.Create("safe.txt")
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// _, err = f.WriteString(client)
			// if err != nil {
			// 	fmt.Println(err)
			// 	f.Close()
			// 	return
			// }
			// err = f.Close()
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			playerState, err = client.PlayerState(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Found your %s (%s)\n", playerState.Device.Type, playerState.Device.Name)
			//close the channel to signal that auth has completed
			close(ch)
		}()
	}
	client = spotify.New(auth.Client(context.Background(), &tok))
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))

	ch <- client

	btys, err := json.Marshal(tok)
	if err != nil {
		fmt.Println(err)
		return
	}

	//save token in safe.txt
	f, err := os.Create("safe.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(string(btys))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func userInfo(userI string) {
	println("User ID: " + userI)
	userID := flag.String("user", userI, userI)

	flag.Parse()

	ctx := context.Background()

	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return
	}

	user, err := client.GetUsersPublicProfile(ctx, spotify.ID(*userID))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("Display name:", user.DisplayName)
	fmt.Println("Spotify URI:", string(user.URI))
	fmt.Println("Endpoint:", user.Endpoint)
	fmt.Println("Followers:", user.Followers.Count)
}

/*Create func for each option :

Get Several Albums
Get Album Tracks
Get Saved Albums
Save Albums
Remove Albums
Check Saved Albums
Get New Releases
Get Artist
Get Several Artists
Get Artist's Albums
Get Artist's Top Tracks
Get Artist's Related Artists
Get Show
Get Several Shows
Get Show Episodes
Get User's Saved Shows
Save Shows for Current User
Remove User's Saved Shows
Check User's Saved Shows
Get Episode
Get Several Episodes
Get User's Saved Episodes
Save Episodes for User
Remove User's Saved Episodes
Check User's Saved Episodes
Get an Audiobook
Get Several Audiobooks
Get a Chapter
Get Several Chapters
Get Track
Get Several Tracks
Get User's Saved Tracks
Save Tracks for Current User
Remove Tracks for Current User
Check User's Saved Tracks
Get Tracks' Audio Features
Get Track's Audio Features
Get Track's Audio Analysis
Get Recommendations
Search for Item
Get Current User's Profile
Get User's Top Items
Get User's Profile
Follow Playlist
Unfollow Playlist
Get Followed Artists
Follow Artists or Users
Unfollow Artists or Users
Check If User Follows Artists or Users
Check if Users Follow Playlist
Get Playlist
Change Playlist Details
Get Playlist Items
Add Items to Playlist
Update Playlist Items
Remove Playlist Items
Get Current User's Playlists
Get User's Playlists
Create Playlist
Get Featured Playlists
Get Category's Playlists
Get Playlist Cover Image
Add Custom Playlist Cover Image
Get Several Browse Categories
Get Single Browse Category
Get Available Genre Seeds
Get Playback State
Transfer Playback
Get Available Devices
Get Currently Playing Track
Start/Resume Playback
Pause Playback
Skip To Next
Skip To Previous
Seek To Position
Set Repeat Mode
Set Playback Volume
Toggle Playback Shuffle
Get Recently Played Tracks
Get the User's Queue
Add Item to Playback Queue
Get Available Markets */

func getSeveralAlbums(ids []spotify.ID) []*spotify.FullAlbum {
	albums, err := client.GetAlbums(context.Background(), ids)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return albums
}

func getAlbumTracks(id spotify.ID) *spotify.SimpleTrackPage {
	tracks, err := client.GetAlbumTracks(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return tracks
}

func getSavedAlbums() *spotify.SavedAlbumPage {
	albums, err := client.CurrentUsersAlbums(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return albums
}

func saveAlbums(ids []spotify.ID) {
	err := client.AddAlbumsToLibrary(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}

func removeAlbums(ids []spotify.ID) {
	err := client.RemoveAlbumsFromLibrary(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}

func checkSavedAlbums(ids []spotify.ID) []bool {
	albums, err := client.UserHasAlbums(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return albums
}

func getNewReleases() *spotify.SimpleAlbumPage {
	albums, err := client.NewReleases(context.Background(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return albums
}

func getArtist(id spotify.ID) *spotify.FullArtist {
	artist, err := client.GetArtist(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return artist
}

func getSeveralArtists(ids []spotify.ID) []*spotify.FullArtist {
	artists, err := client.GetArtists(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return artists
}

func getArtistsAlbums(id spotify.ID) *spotify.SimpleAlbumPage {
	albums, err := client.GetArtistAlbums(context.Background(), id, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return albums
}

func getArtistsTopTracks(id spotify.ID) []spotify.FullTrack {
	tracks, err := client.GetArtistsTopTracks(context.Background(), id, "US")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return tracks
}

func getArtistsRelatedArtists(id spotify.ID) []spotify.FullArtist {
	artists, err := client.GetRelatedArtists(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return artists
}

func getShow(id spotify.ID) *spotify.FullShow {
	show, err := client.GetShow(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return show
}

// func getSeveralShows(ids []spotify.ID) []*spotify.FullShow {
// 	shows, err := client.Sho(context.Background(), ids...)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return shows
// }

func getShowEpisodes(id string) *spotify.SimpleEpisodePage {
	episodes, err := client.GetShowEpisodes(context.Background(), id, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return episodes
}

func getSeveralEpisodes(id string) *spotify.SimpleEpisodePage {
	episodes, err := client.GetShowEpisodes(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return episodes
}

func getTrack(id spotify.ID) *spotify.FullTrack {
	track, err := client.GetTrack(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return track
}

func getSeveralTracks(ids []spotify.ID) []*spotify.FullTrack {
	tracks, err := client.GetTracks(context.Background(), ids)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return tracks
}

func getAudioFeatures(id spotify.ID) []*spotify.AudioFeatures {
	features, err := client.GetAudioFeatures(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return features
}

func getSeveralAudioFeatures(ids []spotify.ID) []*spotify.AudioFeatures {
	features, err := client.GetAudioFeatures(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return features
}

func getAudioAnalysis(id spotify.ID) *spotify.AudioAnalysis {
	analysis, err := client.GetAudioAnalysis(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return analysis
}

func getCategories() *spotify.CategoryPage {
	categories, err := client.GetCategories(context.Background(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return categories
}

func getCategory(id string) spotify.Category {
	category, err := client.GetCategory(context.Background(), id, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return spotify.Category{}
	}

	return category
}

func getCategoryPlaylists(id string) *spotify.SimplePlaylistPage {
	playlists, err := client.GetCategoryPlaylists(context.Background(), id, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return playlists
}

// func getFeaturedPlaylists() *spotify.FeaturedPlaylists {
// 	playlists, err := client.playlist(context.Background(), nil)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return playlists
// }

// func getCategorySPlaylist(id string) *spotify.SimplePlaylistPage {
// 	playlists, err := client.GetCategoryPlaylists(context.Background(), id, nil)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return playlists
// }

/* create func for each option :

Get Playlist Cover Image
Add Custom Playlist Cover Image
Get Several Browse Categories
Get Single Browse Category
Get Available Genre Seeds
Get Playback State
Transfer Playback
Get Available Devices
Get Currently Playing Track
Start/Resume Playback
Pause Playback
Skip To Next
Skip To Previous
Seek To Position
Set Repeat Mode
Set Playback Volume
Toggle Playback Shuffle
Get Recently Played Tracks
Get the User's Queue
Add Item to Playback Queue
Get Available Markets */

// func GetPlaylistCoverImage(id string) *spotify.Image {
// 	image, err := client.playlist(context.Background(), id)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return image
// }

// func AddCustomPlaylistCoverImage(id string) *spotify.Image {
// 	image, err := client.SetPlaylistImage(context.Background(), id)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return image
// }

func GetSeveralBrowseCategories() *spotify.CategoryPage {
	categories, err := client.GetCategories(context.Background(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return categories
}

func GetSingleBrowseCategory(id string) spotify.Category {
	category, err := client.GetCategory(context.Background(), id, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return spotify.Category{}
	}

	return category
}

func GetAvailableGenreSeeds() []string {
	seeds, err := client.GetAvailableGenreSeeds(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return seeds
}

func GetPlaybackState() *spotify.PlayerState {
	state, err := client.PlayerState(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return state
}

func TransferPlayback(deviceID spotify.ID) {
	err := client.TransferPlayback(context.Background(), deviceID, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func GetAvailableDevices() []spotify.PlayerDevice {
	device, err := client.PlayerDevices(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return device
}

func GetCurrentlyPlayingTrack() *spotify.CurrentlyPlaying {
	track, err := client.PlayerCurrentlyPlaying(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return track
}

func StartResumePlayback(deviceID spotify.ID) {
	err := client.PlayOpt(context.Background(), &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func PausePlayback(deviceID spotify.ID) {
	err := client.PauseOpt(context.Background(), &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func SkipToNext(deviceID spotify.ID) {
	err := client.NextOpt(context.Background(), &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func SkipToPrevious(deviceID spotify.ID) {
	err := client.PreviousOpt(context.Background(), &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func SeekToPosition(positionMs int, deviceID spotify.ID) {
	err := client.SeekOpt(context.Background(), positionMs, &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func SetRepeatMode(repeatState string, deviceID spotify.ID) {
	err := client.RepeatOpt(context.Background(), repeatState, &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func SetPlaybackVolume(volumePercent int, deviceID spotify.ID) {
	err := client.VolumeOpt(context.Background(), volumePercent, &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func TogglePlaybackShuffle(shuffleState bool, deviceID spotify.ID) {
	err := client.ShuffleOpt(context.Background(), shuffleState, &spotify.PlayOptions{DeviceID: &deviceID})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func GetRecentlyPlayedTracks() []spotify.RecentlyPlayedItem {
	recentlyPlayed, err := client.PlayerRecentlyPlayed(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return recentlyPlayed
}

// func GetSongQueuePos() *spotify.QueueTrackPage {
// 	queue, err := client.GetPlaylistsForUser(context.Background())
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return queue
// }

func AddItemToPlaybackQueue(id spotify.ID) {
	err := client.QueueSong(context.Background(), id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// func RemoveItemFromPlaybackQueue(id spotify.ID) {
// 	err := client.Son(context.Background(), id)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 	}
// }

func GetSavedAlbums() *spotify.SavedAlbumPage {
	albums, err := client.CurrentUsersAlbums(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return albums
}

func SaveAlbums(ids []spotify.ID) {
	err := client.AddAlbumsToLibrary(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func RemoveAlbums(ids []spotify.ID) {
	err := client.RemoveAlbumsFromLibrary(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func GetSavedTracks() *spotify.SavedTrackPage {
	tracks, err := client.CurrentUsersTracks(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return tracks
}

func SaveTracks(ids []spotify.ID) {
	err := client.AddTracksToLibrary(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func RemoveTracks(ids []spotify.ID) {
	err := client.RemoveTracksFromLibrary(context.Background(), ids...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func GetSavedShows() *spotify.SavedShowPage {
	shows, err := client.CurrentUsersShows(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return shows
}

// func SaveShows(ids []spotify.ID) {
// 	err := client.Show(context.Background(), ids...)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 	}
// }

// func RemoveShows(ids []spotify.ID) {
// 	err := client.RemoveShowsFromLibrary(context.Background(), ids...)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 	}
// }

// func GetSavedEpisodes() *spotify.SavedEpisodePage {
// 	episodes, err := client.Episo(context.Background())
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		return nil
// 	}

// 	return episodes
// }

// func SaveEpisodes(ids []spotify.ID) {
// 	err := client.AddEpisodesToLibrary(context.Background(), ids...)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 	}
// }

// func RemoveEpisodes(ids []spotify.ID) {
// 	err := client.RemoveEpisodesFromLibrary(context.Background(), ids...)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 	}
// }
