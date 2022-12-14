package api

type JsonSearchUser struct {
	Type        string `json:"type"`
	Mid         int    `json:"mid"`
	Uname       string `json:"uname"`
	Usign       string `json:"usign"`
	Fans        int    `json:"fans"`
	Videos      int    `json:"videos"`
	Upic        string `json:"upic"`
	FaceNft     int    `json:"face_nft"`
	FaceNftType int    `json:"face_nft_type"`
	VerifyInfo  string `json:"verify_info"`
	Level       int    `json:"level"`
	Gender      int    `json:"gender"`
	IsUpuser    int    `json:"is_upuser"`
	IsLive      int    `json:"is_live"`
	RoomID      int    `json:"room_id"`
	Res         []struct {
		Aid          int    `json:"aid"`
		Bvid         string `json:"bvid"`
		Title        string `json:"title"`
		Pubdate      int    `json:"pubdate"`
		Arcurl       string `json:"arcurl"`
		Pic          string `json:"pic"`
		Play         string `json:"play"`
		Dm           int    `json:"dm"`
		Coin         int    `json:"coin"`
		Fav          int    `json:"fav"`
		Desc         string `json:"desc"`
		Duration     string `json:"duration"`
		IsPay        int    `json:"is_pay"`
		IsUnionVideo int    `json:"is_union_video"`
	} `json:"res"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	HitColumns     []string `json:"hit_columns"`
	IsSeniorMember int      `json:"is_senior_member"`
}

type JsonSearchRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Seid           string `json:"seid"`
		Page           int    `json:"page"`
		PageSize       int    `json:"pagesize"`
		NumResults     int    `json:"numResults"`
		NumPages       int    `json:"numPages"`
		SuggestKeyword string `json:"suggest_keyword"`
		RqtType        string `json:"rqt_type"`
		CostTime       struct {
			ParamsCheck         string `json:"params_check"`
			GetUpuserLiveStatus string `json:"get upuser live status"`
			IsRiskQuery         string `json:"is_risk_query"`
			IllegalHandler      string `json:"illegal_handler"`
			AsResponseFormat    string `json:"as_response_format"`
			AsRequest           string `json:"as_request"`
			SaveCache           string `json:"save_cache"`
			DeserializeResponse string `json:"deserialize_response"`
			AsRequestFormat     string `json:"as_request_format"`
			Total               string `json:"total"`
			MainHandler         string `json:"main_handler"`
		} `json:"cost_time"`
		ExpList struct {
			Num5508 bool `json:"5508"`
			Num6609 bool `json:"6609"`
			Num7704 bool `json:"7704"`
		} `json:"exp_list"`
		EggHit     int              `json:"egg_hit"`
		Result     []JsonSearchUser `json:"result"`
		ShowColumn int              `json:"show_column"`
		InBlackKey int              `json:"in_black_key"`
		InWhiteKey int              `json:"in_white_key"`
	} `json:"data"`
}

type JsonUserInfo struct {
	Ts   int `json:"ts"`
	Code int `json:"code"`
	Card struct {
		Mid         string  `json:"mid"`
		Name        string  `json:"name"`
		Approve     bool    `json:"approve"`
		Sex         string  `json:"sex"`
		Rank        string  `json:"rank"`
		Face        string  `json:"face"`
		Coins       float64 `json:"coins"`
		DisplayRank string  `json:"DisplayRank"`
		Regtime     int     `json:"regtime"`
		Spacesta    int     `json:"spacesta"`
		Place       string  `json:"place"`
		Birthday    string  `json:"birthday"`
		Sign        string  `json:"sign"`
		Description string  `json:"description"`
		Article     int     `json:"article"`
		Attentions  []int64 `json:"attentions"`
		Fans        int     `json:"fans"`
		Friend      int     `json:"friend"`
		Attention   int     `json:"attention"`
		LevelInfo   struct {
			NextExp      int `json:"next_exp"`
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
		} `json:"level_info"`
		Pendant struct {
			Pid    int    `json:"pid"`
			Name   string `json:"name"`
			Image  string `json:"image"`
			Expire int    `json:"expire"`
		} `json:"pendant"`
		OfficialVerify struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"official_verify"`
	} `json:"card"`
}
