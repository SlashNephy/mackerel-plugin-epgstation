package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
)

type EPGStationAPI struct {
	baseUrl *url.URL
	client  *http.Client
}

func NewEPGStationAPI(host string, port int) *EPGStationAPI {
	return &EPGStationAPI{
		baseUrl: &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		},
		client: &http.Client{},
	}
}

type EPGStationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type EPGStationStreams struct {
	Items []struct {
		StreamId    int    `json:"streamId"`
		Type        string `json:"type"`
		Mode        int    `json:"mode"`
		IsEnable    bool   `json:"isEnable"`
		ChannelId   int64  `json:"channelId"`
		VideoFileId int    `json:"videoFileId"`
		RecordedId  int    `json:"recordedId"`
		Name        string `json:"name"`
		StartAt     int    `json:"startAt"`
		EndAt       int    `json:"endAt"`
		Description string `json:"description"`
		Extended    string `json:"extended"`
	} `json:"items"`
	EPGStationError
}

func (e *EPGStationAPI) GetStreams() (*EPGStationStreams, error) {
	var result EPGStationStreams
	if err := e.get(&url.URL{Path: "/api/streams", RawQuery: "isHalfWidth=false"}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type EPGStationReserveCounts struct {
	Normal    int `json:"normal"`
	Conflicts int `json:"conflicts"`
	Skips     int `json:"skips"`
	Overlaps  int `json:"overlaps"`
	EPGStationError
}

func (e *EPGStationAPI) GetReserveCounts() (*EPGStationReserveCounts, error) {
	var result EPGStationReserveCounts
	if err := e.get(&url.URL{Path: "/api/reserves/cnts"}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type EPGStationRecording struct {
	Records []struct {
		Id          int    `json:"id"`
		RuleId      int    `json:"ruleId"`
		ProgramId   int64  `json:"programId"`
		ChannelId   int64  `json:"channelId"`
		StartAt     int    `json:"startAt"`
		EndAt       int    `json:"endAt"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Extended    string `json:"extended"`
		RawExtended struct {
		} `json:"rawExtended"`
		Genre1             int    `json:"genre1"`
		SubGenre1          int    `json:"subGenre1"`
		Genre2             int    `json:"genre2"`
		SubGenre2          int    `json:"subGenre2"`
		Genre3             int    `json:"genre3"`
		SubGenre3          int    `json:"subGenre3"`
		VideoType          string `json:"videoType"`
		VideoResolution    string `json:"videoResolution"`
		VideoStreamContent int    `json:"videoStreamContent"`
		VideoComponentType int    `json:"videoComponentType"`
		AudioSamplingRate  int    `json:"audioSamplingRate"`
		AudioComponentType int    `json:"audioComponentType"`
		IsRecording        bool   `json:"isRecording"`
		Thumbnails         []int  `json:"thumbnails"`
		VideoFiles         []struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Filename string `json:"filename"`
			Type     string `json:"type"`
			Size     int    `json:"size"`
		} `json:"videoFiles"`
		DropLog struct {
			Id            int `json:"id"`
			ErrorCnt      int `json:"errorCnt"`
			DropCnt       int `json:"dropCnt"`
			ScramblingCnt int `json:"scramblingCnt"`
		} `json:"dropLog"`
		Tags []struct {
			Id    int    `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"tags"`
		IsEncoding  bool `json:"isEncoding"`
		IsProtected bool `json:"isProtected"`
	} `json:"records"`
	Total int `json:"total"`
	EPGStationError
}

func (e *EPGStationAPI) GetRecording() (*EPGStationRecording, error) {
	var result EPGStationRecording
	if err := e.get(&url.URL{Path: "/api/recording", RawQuery: "isHalfWidth=false"}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type EPGStationEncode struct {
	RunningItems []struct {
		Id       int    `json:"id"`
		Mode     string `json:"mode"`
		Recorded struct {
			Id          int    `json:"id"`
			RuleId      int    `json:"ruleId"`
			ProgramId   int64  `json:"programId"`
			ChannelId   int64  `json:"channelId"`
			StartAt     int    `json:"startAt"`
			EndAt       int    `json:"endAt"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Extended    string `json:"extended"`
			RawExtended struct {
			} `json:"rawExtended"`
			Genre1             int    `json:"genre1"`
			SubGenre1          int    `json:"subGenre1"`
			Genre2             int    `json:"genre2"`
			SubGenre2          int    `json:"subGenre2"`
			Genre3             int    `json:"genre3"`
			SubGenre3          int    `json:"subGenre3"`
			VideoType          string `json:"videoType"`
			VideoResolution    string `json:"videoResolution"`
			VideoStreamContent int    `json:"videoStreamContent"`
			VideoComponentType int    `json:"videoComponentType"`
			AudioSamplingRate  int    `json:"audioSamplingRate"`
			AudioComponentType int    `json:"audioComponentType"`
			IsRecording        bool   `json:"isRecording"`
			Thumbnails         []int  `json:"thumbnails"`
			VideoFiles         []struct {
				Id       int    `json:"id"`
				Name     string `json:"name"`
				Filename string `json:"filename"`
				Type     string `json:"type"`
				Size     int    `json:"size"`
			} `json:"videoFiles"`
			DropLog struct {
				Id            int `json:"id"`
				ErrorCnt      int `json:"errorCnt"`
				DropCnt       int `json:"dropCnt"`
				ScramblingCnt int `json:"scramblingCnt"`
			} `json:"dropLog"`
			Tags []struct {
				Id    int    `json:"id"`
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"tags"`
			IsEncoding  bool `json:"isEncoding"`
			IsProtected bool `json:"isProtected"`
		} `json:"recorded"`
		Percent int    `json:"percent"`
		Log     string `json:"log"`
	} `json:"runningItems"`
	WaitItems []struct {
		Id       int    `json:"id"`
		Mode     string `json:"mode"`
		Recorded struct {
			Id          int    `json:"id"`
			RuleId      int    `json:"ruleId"`
			ProgramId   int64  `json:"programId"`
			ChannelId   int64  `json:"channelId"`
			StartAt     int    `json:"startAt"`
			EndAt       int    `json:"endAt"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Extended    string `json:"extended"`
			RawExtended struct {
			} `json:"rawExtended"`
			Genre1             int    `json:"genre1"`
			SubGenre1          int    `json:"subGenre1"`
			Genre2             int    `json:"genre2"`
			SubGenre2          int    `json:"subGenre2"`
			Genre3             int    `json:"genre3"`
			SubGenre3          int    `json:"subGenre3"`
			VideoType          string `json:"videoType"`
			VideoResolution    string `json:"videoResolution"`
			VideoStreamContent int    `json:"videoStreamContent"`
			VideoComponentType int    `json:"videoComponentType"`
			AudioSamplingRate  int    `json:"audioSamplingRate"`
			AudioComponentType int    `json:"audioComponentType"`
			IsRecording        bool   `json:"isRecording"`
			Thumbnails         []int  `json:"thumbnails"`
			VideoFiles         []struct {
				Id       int    `json:"id"`
				Name     string `json:"name"`
				Filename string `json:"filename"`
				Type     string `json:"type"`
				Size     int    `json:"size"`
			} `json:"videoFiles"`
			DropLog struct {
				Id            int `json:"id"`
				ErrorCnt      int `json:"errorCnt"`
				DropCnt       int `json:"dropCnt"`
				ScramblingCnt int `json:"scramblingCnt"`
			} `json:"dropLog"`
			Tags []struct {
				Id    int    `json:"id"`
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"tags"`
			IsEncoding  bool `json:"isEncoding"`
			IsProtected bool `json:"isProtected"`
		} `json:"recorded"`
		Percent int    `json:"percent"`
		Log     string `json:"log"`
	} `json:"waitItems"`
	EPGStationError
}

func (e *EPGStationAPI) GetEncode() (*EPGStationEncode, error) {
	var result EPGStationEncode
	if err := e.get(&url.URL{Path: "/api/encode", RawQuery: "isHalfWidth=false"}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type EPGStationStorages struct {
	Items []struct {
		Name      string `json:"name"`
		Available int    `json:"available"`
		Used      int    `json:"used"`
		Total     int    `json:"total"`
	} `json:"items"`
	EPGStationError
}

func (e *EPGStationAPI) GetStorages() (*EPGStationStorages, error) {
	var result EPGStationStorages
	if err := e.get(&url.URL{Path: "/api/storages"}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (e *EPGStationAPI) get(path *url.URL, result any) error {
	requestUrl := e.baseUrl.ResolveReference(path)
	request, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", "mackerel-plugin-epgstation (+https://github.com/SlashNephy/mackerel-plugin-epgstation)")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}

	return nil
}
