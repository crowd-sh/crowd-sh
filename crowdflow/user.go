package crowdflow

type User struct {
	Id        int64
	ApiKey    string `sql:"not null;unique"`
	Workflows []Workflow
}

func FindUserByApi(api_key string) (u User) {
	Db.Where(&User{ApiKey: api_key}).First(&u)
	return
}
