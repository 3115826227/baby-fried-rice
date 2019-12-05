package model

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RspStationItem struct {
	StationNum  int
	StationName string
	ArriveTime  string
	StartTime   string
	StopMinute  int
}

type RspTrainMeta struct {
	Data []struct {
		Date             string `json:"date"`
		FromStation      string `json:"from_station"`
		StationTrainCode string `json:"station_train_code"`
		ToStation        string `json:"to_station"`
		TotalNum         string `json:"total_num"`
		TrainNo          string `json:"train_no"`
	} `json:"data"`
	Status   bool   `json:"status"`
	ErrorMsg string `json:"errorMsg"`
}

type RspStationTrainMeta struct {
	Data struct {
		Data []struct {
			StationTrainCode string `json:"station_train_code"`
			TrainNo          string `json:"train_no"`
			StartStationName string `json:"start_station_name"`
			EndStationName   string `json:"end_station_name"`
		} `json:"data"`
	} `json:"data"`
}

type RspTrainMetaWayStation struct {
	Data struct {
		DptStation  string `json:"dptStation"`
		ArrStation  string `json:"arrStation"`
		DptDate     string `json:"dptDate"`
		DptCityName string `json:"dptCityName"`
		ArrCityName string `json:"arrCityName"`
		S2SBean     struct {
			DayDifference string `json:"dayDifference"`
		} `json:"s2sBean"`
		StationItemList []struct {
			StationNo   int    `json:"stationNo"`
			StationName string `json:"stationName"`
			ArriveTime  string `json:"arriveTime"`
			StartTime   string `json:"startTime"`
			OverTime    int    `json:"overTime"`
		} `json:"stationItemList"`
		Distance int    `json:"distance"`
		TypeName string `json:"typeName"`
		Sale     string `json:"sale"`
	} `json:"data"`
}

type RspQunarTrainSeatPrice struct {
	Data struct {
		DptStation  string `json:"dptStation"`
		ArrStation  string `json:"arrStation"`
		DptDate     string `json:"dptDate"`
		DptCityName string `json:"dptCityName"`
		ArrCityName string `json:"arrCityName"`
		S2SBeanList []struct {
			DayDifference string `json:"dayDifference"`
			Seats         map[string]struct {
				Price float32 `json:"price"`
				Count int     `json:"count"`
			} `json:"seats"`
			TrainNo        string `json:"trainNo"`
			StartDate      string `json:"startDate"`
			DptStationName string `json:"dptStationName"`
			DptStationCode string `json:"dptStationCode"`
			ArrStationName string `json:"arrStationName"`
			ArrStationCode string `json:"arrStationCode"`
		} `json:"s2sBeanList"`
		Distance int    `json:"distance"`
		TypeName string `json:"typeName"`
		Sale     string `json:"sale"`
	} `json:"data"`
}

type RspZhiXingTrainMeta struct {
	TrainStopList []struct {
		StationSequence   int    `json:"StationSequence"`
		StationName       string `json:"StationName"`
		DepartureTime     string `json:"DepartureTime"`
		ArrivalTime       string `json:"ArrivalTime"`
		StopTimes         int    `json:"StopTimes"`
		DistanceFromStart int    `json:"DistanceFromStart"`
	} `json:"TrainStopList"`
}

type RspZhiXingTrainSeatPrice struct {
	ResponseBody struct {
		DepartureCity struct {
			CityName string `json:"CityName"`
		} `json:"DepartureCity"`
		ArriveCity struct {
			CityName string `json:"CityName"`
		} `json:"ArriveCity"`
		TrainItems []struct {
			TrainName    string `json:"TrainName"`
			TicketResult struct {
				TicketItems []struct {
					SeatTypeName string  `json:"SeatTypeName"`
					Price        float64 `json:"Price"`
				} `json:"TicketItems"`
			} `json:"TicketResult"`
		} `json:"TrainItems"`
	} `json:"ResponseBody"`
}

type RspTongChengYiLongTrainSeatPrice struct {
	Data struct {
		Trains []struct {
			TrainNum    string `json:"trainNum"`
			TicketState map[string]struct {
				Cn    string `json:"cn"`
				Price string `json:"price"`
				Seats string `json:"seats"`
			} `json:"ticketState"`
		} `json:"trains"`
		FromCityName string `json:"fromCityName"`
		ToCityName   string `json:"toCityName"`
	} `json:"data"`
}

type RspMeituanTrainMeta struct {
	Data struct {
		Stations []struct {
			ArriveTime  string `json:"arrive_time"`
			StartTime   string `json:"start_time"`
			StationName string `json:"station_name"`
			StationNo   string `json:"station_no"`
			StopTime    string `json:"stop_time"`
		} `json:"stations"`
	} `json:"data"`
}

type RspMeituanTrainSeatPrice struct {
	Data struct {
		FromCityName string `json:"fromCityName"`
		ToCityName   string `json:"toCityName"`
		Trains       []struct {
			Seats []struct {
				SeatPrice    float64 `json:"seat_price"`
				SeatTypeName string  `json:"seat_type_name"`
				SeatYupiao   int     `json:"seat_yupiao"`
			} `json:"seats"`
			TrainCode string `json:"train_code"`
		} `json:"trains"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type RspJindongTrainSeatPrice struct {
	Data struct {
		Value []struct {
			Seats []struct {
				Price    string `json:"price"`
				SeatName string `json:"seatName"`
			} `json:"seats"`
			TrainCode string `json:"trainCode"`
		} `json:"value"`
	} `json:"data"`
	Success bool `json:"success"`
}

type RspTriggerTrain struct {
	Train         string `json:"train"`
	StartStation  string `json:"start_station"`
	StartTime     string `json:"start_time"`
	ArriveStation string `json:"arrive_station"`
	ArriveTime    string `json:"arrive_time"`
	OverDay       int    `json:"over_day"`
}

type RspTrainMetaTrigger struct {
	Date               string            `json:"date"`
	IsTrigger          bool              `json:"is_trigger"`
	TriggerTrainNumber int               `json:"trigger_train_number"`
	TriggerTrain       []RspTriggerTrain `json:"trigger_train"`
}
