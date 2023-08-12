package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// InstagramPostMediaData is pre-processing data from Picuki scraping
type InstagramPostMediaData struct {
	URL string `json:"url"`
}

// InstagramPostData is pre-processing data from Picuki scraping
type InstagramPostData struct {
	ID        string                   `json:"id"`
	Text      string                   `json:"text"`
	MediaList []InstagramPostMediaData `json:"mediaList"`
	CreatedAt string                   `json:"createdAt"`
}

// UserInstagramPosts is pre-processing data from Picuki scraping
type UserInstagramPosts struct {
	Username string              `json:"username"`
	Posts    []InstagramPostData `json:"posts"`
}

// LoadInstagramPostList load data from remote JSON
func LoadInstagramPostList(cacheData *[]UserInstagramPosts) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-instagram-posts/instagram-posts.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			postList := []UserInstagramPosts{}
			json.Unmarshal(bodyBuf, &postList)
			*cacheData = postList
		}
	}
}

// FindUserInstagramPosts finds user Instagram posts by specified user
func FindUserInstagramPosts(postList []UserInstagramPosts, instagramUser string) UserInstagramPosts {
	return funk.Find(postList, func(instagramPosts UserInstagramPosts) bool {
		return instagramPosts.Username == instagramUser
	}).(UserInstagramPosts)
}

// FindInstagramPost finds Instagram post by specified ID
func FindInstagramPost(posts []InstagramPostData, picukiMediaID string) InstagramPostData {
	return funk.Find(posts, func(post InstagramPostData) bool {
		return post.ID == picukiMediaID
	}).(InstagramPostData)
}
