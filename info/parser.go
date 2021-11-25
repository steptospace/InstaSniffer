package info

import (
	"InstaSniffer/api"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func UploadData(url string) (err error, ii api.ImportantInfo) {

	url = "http://www.instagram.com/" + url + "/?__a=1"
	userData := Configure{}
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return err, ii
	}

	time.Sleep(time.Millisecond * 10)

	req.Header.Set("User-Agent", "MSIE/15.0")

	res, err := spaceClient.Do(req)
	if err != nil {
		log.Error(err)
		return err, ii
	}

	// We can check  cookies
	//fmt.Println(res.Cookies())

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return err, ii
	}

	err = json.Unmarshal(body, &userData)
	if err != nil {
		return err, ii
	}

	//fmt.Println(userData.Graphql)
	accounting(infoAboutUser(userData.Graphql.User))

	return nil, infoAboutUser(userData.Graphql.User)
}

// Определится с полями какие в бд нужны и в каком виде

func accounting(data api.ImportantInfo) {
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(data.Name+".json", file, 0644)
	if err != nil {
		log.Error(err)
	}
}

func infoAboutUser(a UserInfo) (b api.ImportantInfo) {
	//Name
	b.Name = a.FullName

	//Username
	b.Username = a.Username

	// Bio
	b.Bio = a.Biography

	//Created time
	b.CreatedAt = time.Now()

	//Avatars
	b.Avatar = a.ProfilePicURLHd

	for _, j := range a.EdgeOwnerToTimelineMedia.Edges {
		mediaEdges := j.Node.EdgeMediaToCaption.Edges
		desc := ""
		if len(mediaEdges) != 0 {
			desc = mediaEdges[0].Node.Text
		}

		if j.Node.IsVideo == true {
			b.Videos = append(b.Videos, api.Media{
				Url:         j.Node.DisplayURL,
				Description: desc})
		} else {
			b.Images = append(b.Images, api.Media{
				Url:         j.Node.DisplayURL,
				Description: desc})
		}
	}
	return b
}
