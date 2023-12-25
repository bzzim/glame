package models

type Config struct {
	WeatherAPIKey           string  `json:"WEATHER_API_KEY"`
	AppsSameTab             bool    `json:"appsSameTab"`
	BookmarksSameTab        bool    `json:"bookmarksSameTab"`
	CustomTitle             string  `json:"customTitle"`
	DaySchema               string  `json:"daySchema"`
	DefaultSearchProvider   string  `json:"defaultSearchProvider"`
	DefaultTheme            string  `json:"defaultTheme"`
	DisableAutofocus        bool    `json:"disableAutofocus"`
	DockerApps              bool    `json:"dockerApps"`
	DockerHost              string  `json:"dockerHost"`
	GreetingsSchema         string  `json:"greetingsSchema"`
	HideApps                bool    `json:"hideApps"`
	HideCategories          bool    `json:"hideCategories"`
	HideDate                bool    `json:"hideDate"`
	HideHeader              bool    `json:"hideHeader"`
	HideSearch              bool    `json:"hideSearch"`
	IsCelsius               bool    `json:"isCelsius"`
	IsKilometer             bool    `json:"isKilometer"`
	KubernetesApps          bool    `json:"kubernetesApps"`
	Lat                     float64 `json:"lat"`
	Lon                     float64 `json:"long"`
	MonthSchema             string  `json:"monthSchema"`
	PinAppsByDefault        bool    `json:"pinAppsByDefault"`
	PinCategoriesByDefault  bool    `json:"pinCategoriesByDefault"`
	SearchSameTab           bool    `json:"searchSameTab"`
	SecondarySearchProvider string  `json:"secondarySearchProvider"`
	ShowTime                bool    `json:"showTime"`
	UnpinStoppedApps        bool    `json:"unpinStoppedApps"`
	UseAmericanDate         bool    `json:"useAmericanDate"`
	UseOrdering             string  `json:"useOrdering"`
	WeatherData             string  `json:"weatherData"`
}
