package huya

import (
	"encoding/json"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
	"github.com/pkg/errors"
)

func init() {
	extractors.Register("huya", New())
}

type extractor struct{}

// New returns a huya extractor.
func New() extractors.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	html, err := request.Get(url, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	title := "huya video"

	var videoUrl string
	quality := "normal"

	videoDesc := utils.MatchOneOf(html, `window\.HNF_GLOBAL_INIT =\s+(.+)\s+<\/script>`)

	streams := make(map[string]*extractors.Stream)

	if len(videoDesc) > 1 {
		var videoData VideoInfo
		if err := json.Unmarshal([]byte(videoDesc[1]), &videoData); err != nil {
			return nil, errors.WithStack(extractors.ErrURLParseFailed)
		}
		videoUrl = videoData.VideoData.VideoDefinitions.Value[0].SUrl
		title = videoData.VideoData.VideoTitle

		for _, v := range videoData.VideoData.VideoDefinitions.Value {
			videoUrl = v.SUrl
			quality = v.SHeight
			size, err := request.Size(videoUrl, url)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			urlData := &extractors.Part{
				URL:  videoUrl,
				Size: size,
				Ext:  "mp4",
			}
			//_ = urlData
			streams[quality] = &extractors.Stream{
				Parts:   []*extractors.Part{urlData},
				Size:    size,
				Quality: quality,
			}
		}

	} else {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	return []*extractors.Data{
		{
			Site:    "虎牙 huya.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
		},
	}, nil
}

type VideoInfo struct {
	RecommendChannelList []struct {
		Title           string `json:"title"`
		Cover           string `json:"cover"`
		Vid             int    `json:"vid"`
		SetId           int    `json:"setId"`
		ChannelItemList struct {
			BValue int `json:"_bValue"`
			Value  []struct {
				Id        string `json:"id"`
				Name      string `json:"name"`
				Classname string `json:"_classname"`
			} `json:"value"`
			Classname string `json:"_classname"`
		} `json:"channelItemList"`
		Level     int    `json:"level"`
		Url       string `json:"url"`
		Classname string `json:"_classname"`
	} `json:"recommendChannelList"`
	VideoData struct {
		Vid             int    `json:"vid"`
		LiveId          string `json:"liveId"`
		VideoTitle      string `json:"videoTitle"`
		Brief           string `json:"brief"`
		SourceClient    int    `json:"sourceClient"`
		ChannelId       string `json:"channelId"`
		VideoType       int    `json:"videoType"`
		UploadStartTime int    `json:"uploadStartTime"`
		UploadEndTime   int    `json:"uploadEndTime"`
		Covers          struct {
			Cover           string `json:"cover"`
			Cover320        string `json:"cover320"`
			Cover375        string `json:"cover375"`
			Cover640        string `json:"cover640"`
			TuijianCover    string `json:"tuijianCover"`
			OriginCover     string `json:"originCover"`
			FirstFrameCover string `json:"firstFrameCover"`
			Classname       string `json:"_classname"`
		} `json:"covers"`
		Tags        string `json:"tags"`
		PlayNums    int    `json:"playNums"`
		FavorNums   int    `json:"favorNums"`
		ShareNums   int    `json:"shareNums"`
		CommentNums int    `json:"commentNums"`
		DanmuNums   int    `json:"danmuNums"`
		VideoUrl    string `json:"videoUrl"`
		Duration    string `json:"duration"`
		UpdateTime  int    `json:"updateTime"`
		UserInfo    struct {
			Uid          int    `json:"uid"`
			UserNick     string `json:"userNick"`
			Avatar       string `json:"avatar"`
			Homepage     string `json:"homepage"`
			IsOnline     bool   `json:"isOnline"`
			LastLiveTime int    `json:"lastLiveTime"`
			RoomId       int    `json:"roomId"`
			Classname    string `json:"_classname"`
		} `json:"userInfo"`
		VideoDefinitions struct {
			BValue int `json:"_bValue"`
			Value  []struct {
				SSize       string `json:"sSize"`
				SWidth      string `json:"sWidth"`
				SHeight     string `json:"sHeight"`
				SDefinition string `json:"sDefinition"`
				SUrl        string `json:"sUrl"`
				SM3U8       string `json:"sM3u8"`
				SDefName    string `json:"sDefName"`
				SDuration   string `json:"sDuration"`
				SFormat     string `json:"sFormat"`
				STs1Url     string `json:"sTs1Url"`
				STs1Offset  string `json:"sTs1Offset"`
				ILine       int    `json:"iLine"`
				IBitrate    int    `json:"iBitrate"`
				ExtParam    struct {
					Kproto struct {
						Classname string `json:"_classname"`
					} `json:"_kproto"`
					Vproto struct {
						Classname string `json:"_classname"`
					} `json:"_vproto"`
					BKey   int `json:"_bKey"`
					BValue int `json:"_bValue"`
					Value  struct {
						Vid        string `json:"vid"`
						Client     string `json:"client"`
						Pid        string `json:"pid"`
						Bitrate    string `json:"bitrate"`
						Definition string `json:"definition"`
						Scene      string `json:"scene"`
					} `json:"value"`
					Classname string `json:"_classname"`
				} `json:"extParam"`
				Streams struct {
					BValue int `json:"_bValue"`
					Value  []struct {
						SSize      string `json:"sSize"`
						SWidth     string `json:"sWidth"`
						SHeight    string `json:"sHeight"`
						SFormat    string `json:"sFormat"`
						SUrl       string `json:"sUrl"`
						SM3U8      string `json:"sM3u8"`
						STs1Url    string `json:"sTs1Url"`
						STs1Offset string `json:"sTs1Offset"`
						IBitrate   int    `json:"iBitrate"`
						Codec      string `json:"codec"`
						ExtParam   struct {
							Kproto struct {
								Classname string `json:"_classname"`
							} `json:"_kproto"`
							Vproto struct {
								Classname string `json:"_classname"`
							} `json:"_vproto"`
							BKey   int `json:"_bKey"`
							BValue int `json:"_bValue"`
							Value  struct {
								Vid        string `json:"vid"`
								Client     string `json:"client"`
								Pid        string `json:"pid"`
								Bitrate    string `json:"bitrate"`
								Definition string `json:"definition"`
								Scene      string `json:"scene"`
							} `json:"value"`
							Classname string `json:"_classname"`
						} `json:"extParam"`
						Classname string `json:"_classname"`
					} `json:"value"`
					Classname string `json:"_classname"`
				} `json:"streams"`
				Codec     string `json:"codec"`
				Classname string `json:"_classname"`
			} `json:"value"`
			Classname string `json:"_classname"`
		} `json:"videoDefinitions"`
		SourceClientName string  `json:"sourceClientName"`
		ChannelName      string  `json:"channelName"`
		VideoTypeName    string  `json:"videoTypeName"`
		RawDuration      float64 `json:"rawDuration"`
		Topics           struct {
			BValue    int           `json:"_bValue"`
			Value     []interface{} `json:"value"`
			Classname string        `json:"_classname"`
		} `json:"topics"`
		ImplicitTags string `json:"implicitTags"`
		TemplateId   int    `json:"templateId"`
		CanPlay      int    `json:"canPlay"`
		Status       int    `json:"status"`
		PublishTime  int    `json:"publishTime"`
		Classname    string `json:"_classname"`
	} `json:"videoData"`
	MatchInfo struct {
		MatchId           int    `json:"matchId"`
		MatchDesc         string `json:"matchDesc"`
		MatchScheduleId   int    `json:"matchScheduleId"`
		MatchScheduleDesc string `json:"matchScheduleDesc"`
		Classname         string `json:"_classname"`
	} `json:"matchInfo"`
	VideoCollection []interface{} `json:"videoCollection"`
}
