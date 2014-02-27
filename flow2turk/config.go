package crowdflow

import (
        "github.com/crowdmob/goamz/exp/mturk"
)

var (
        config Config
)

type Config struct {
        WithMTurk    bool
        AwsAccessKey string
        AwsSecretKey string
        AwsSandbox bool
        mturkAuth    *mturk.MTurk
}

func init() {
        config = Config{
                WithMTurk: true,
        }

}

func SetConfig(c Config) {
        config = c
}

func EnableMTurk(access_key, secret_key string, sandbox bool) {

}
