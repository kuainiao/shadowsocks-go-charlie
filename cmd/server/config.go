package main
import (
    "os"
    "io/ioutil"
    "encoding/json"
    "fmt"
    ss "bitbucket.org/qiuyuzhou/shadowsocks/core"
    "errors"
    log "github.com/Sirupsen/logrus"
)

type Config struct {
    Listen []string     `json:"listen"`
    Method string       `json:"method"`
    Password string     `json:"password"`
    Tokens map[string]string    `json:"tokens"`
    Timeout uint        `json:"timeout"`

    headerCipher *ss.Cipher
}

func ParseConfig(path string) (config *Config, err error) {
    file, err := os.Open(path) // For read access.
    if err != nil {
        return
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return
    }

    config = &Config{}
    if err = json.Unmarshal(data, config); err != nil {
        return nil, err
    }

    return
}

func (c *Config)GetServerSecret() string {
    return c.Password
}

func (c *Config)GetEncryptMethod() string {
    return c.Method
}

func (c *Config)NewHeaderCipher() *ss.Cipher {
    return c.headerCipher.Copy()
}

func (c *Config)Validate() (bool, error) {
    var err error
    valid := true

    if len(c.Listen) == 0 {
        log.Error("Must specify address for server")
        valid = false
    }
    if c.Password == "" {
        log.Error("Must specify password for server")
        valid = false
    }
    if c.Method == "" {
        log.Error("Must specify method for server")
        valid = false
    }

    for key, _:= range c.Tokens {
        if len(key) > ss.TOKEN_SIZE {
            log.WithField("token", key).Errorf(fmt.Sprintf("Token lenght must be less equal %v", ss.TOKEN_SIZE))
            valid = false
        }
    }

    if !valid {
        return valid, errors.New("Invalid config file")
    }

    if c.headerCipher, err = ss.NewCipher(c.Method, c.Password); err != nil {
        valid = false
        return valid, err
    }

    return valid, nil
}
