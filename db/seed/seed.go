package main

import (
	"api/app/models"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SomeStructWithTags struct {
	Latitude         float32 `faker:"lat"`
	Longitude        float32 `faker:"long"`
	CreditCardNumber string  `faker:"cc_number"`
	CreditCardType   string  `faker:"cc_type"`
	Email            string  `faker:"email"`
	DomainName       string  `faker:"domain_name"`
	IPV4             string  `faker:"ipv4"`
	IPV6             string  `faker:"ipv6"`
	Password         string  `faker:"password"`
	// Jwt                string  `faker:"jwt"`
	PhoneNumber        string  `faker:"phone_number"`
	MacAddress         string  `faker:"mac_address"`
	URL                string  `faker:"url"`
	UserName           string  `faker:"username"`
	TollFreeNumber     string  `faker:"toll_free_number"`
	E164PhoneNumber    string  `faker:"e_164_phone_number"`
	TitleMale          string  `faker:"title_male"`
	TitleFemale        string  `faker:"title_female"`
	FirstName          string  `faker:"first_name"`
	FirstNameMale      string  `faker:"first_name_male"`
	FirstNameFemale    string  `faker:"first_name_female"`
	LastName           string  `faker:"last_name"`
	Name               string  `faker:"name"`
	UnixTime           int64   `faker:"unix_time"`
	Date               string  `faker:"date"`
	Time               string  `faker:"time"`
	MonthName          string  `faker:"month_name"`
	Year               string  `faker:"year"`
	DayOfWeek          string  `faker:"day_of_week"`
	DayOfMonth         string  `faker:"day_of_month"`
	Timestamp          string  `faker:"timestamp"`
	Century            string  `faker:"century"`
	TimeZone           string  `faker:"timezone"`
	TimePeriod         string  `faker:"time_period"`
	Word               string  `faker:"word"`
	Sentence           string  `faker:"sentence,unique"`
	Paragraph          string  `faker:"paragraph,unique"`
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"amount"`
	AmountWithCurrency string  `faker:"amount_with_currency"`
	UUIDHypenated      string  `faker:"uuid_hyphenated"`
	UUID               string  `faker:"uuid_digit"`
	Skip               string  `faker:"-"`
	PaymentMethod      string  `faker:"oneof: cc, paypal, check, money order"` // oneof will randomly pick one of the comma-separated values supplied in the tag
	AccountID          int     `faker:"oneof: 15, 27, 61"`                     // use commas to separate the values for now. Future support for other separator characters may be added
	Price32            float32 `faker:"oneof: 4.95, 9.99, 31997.97"`
	Price64            float64 `faker:"oneof: 47463.9463525, 993747.95662529, 11131997.978767990"`
	NumS64             int64   `faker:"oneof: 1, 2"`
	NumS32             int32   `faker:"oneof: -3, 4"`
	NumS16             int16   `faker:"oneof: -5, 6"`
	NumS8              int8    `faker:"oneof: 7, -8"`
	NumU64             uint64  `faker:"oneof: 9, 10"`
	NumU32             uint32  `faker:"oneof: 11, 12"`
	NumU16             uint16  `faker:"oneof: 13, 14"`
	NumU8              uint8   `faker:"oneof: 15, 16"`
	NumU               uint    `faker:"oneof: 17, 18"`
}

func main() {
	flag.Parse()
	fmt.Printf("%v\n", flag.Args())
	db := models.DbConnection
	if err := seeds(db, flag.Args()); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	return
}

func seeds(db *gorm.DB, seedList []string) error {
	if len(seedList) > 0 {
		for _, seed := range seedList {
			switch seed {
			case "categories":
				categorySeed(db)
			case "admin_users":
				adminUserSeed(db)
			case "articles":
				articleSeed(db)
			case "artist_infos":
				artistInfoSeed(db)
			default:
				return nil
			}
		}
	} else {
		categorySeed(db)
		adminUserSeed(db)
		articleSeed(db)
		artistInfoSeed(db)
	}

	return nil
}

func adminUserSeed(db *gorm.DB) {
	// 1アカウント作成
	// for i := 1; i <= 10; i++ {
	hash, err := bcrypt.GenerateFromPassword([]byte("tama0413"), 10)
	if err != nil {
		fmt.Println(err)
	}
	pass := string(hash)

	adminUser := models.AdminUser{
		Name:      "Kenji",
		Email:     "kenji.com",
		Password:  pass,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&adminUser).Error; err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// }
	fmt.Println("Admin User seed seccessful!!")
}

func categorySeed(db *gorm.DB) {
	// ニュース1〜ニュース10まで
	for i := 1; i <= 10; i++ {
		category := models.Category{
			Name:      "ニュース" + strconv.Itoa(i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&category).Error; err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
	}
	fmt.Println("Category seed seccessful!!")
}

func articleSeed(db *gorm.DB) {
	a := SomeStructWithTags{}

	// 記事1〜記事10まで
	for i := 1; i <= 10; i++ {
		// faker作成
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}

		// id：5までは1、以降は2
		categoryID := 1
		if i > 5 {
			categoryID = 2
		}

		article := models.Article{
			Title:     a.Sentence,
			Text:      a.Paragraph,
			Category:  categoryID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&article).Error; err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
	}
	fmt.Println("Article seed seccessful!!")
}

func artistInfoSeed(db *gorm.DB) {
	// ニュース1〜ニュース10まで
	for _, artistInfoData := range artistInfoDatas {
		artistInfo := models.ArtistInfo{
			ArtistId:  artistInfoData.ArtistId,
			Name:      artistInfoData.Name,
			Url:       artistInfoData.Url,
			TwitterId: artistInfoData.TwitterId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&artistInfo).Error; err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
	}
	fmt.Println("ArtistInfo seed seccessful!!")
}

var artistInfoDatas = []models.ArtistInfo{
	{
		ArtistId:  "0zEbGW70TQHSOf4Ip1oeVn",
		Name:      "ACIDMAN",
		Url:       "http://acidman.jp/content/",
		TwitterId: "acidman_staff",
	},
	{
		ArtistId:  "3NTbOmzlj2cL86XFuDVFvZ",
		Name:      "MAN WITH A MISSION",
		Url:       "https://www.mwamjapan.info/",
		TwitterId: "MWAMofficial",
	},
}
