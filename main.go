package main

import (
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
	rtmtokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtmTokenBuilder"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Token struct {
	UID      uint32            `json:"uid"`
	CHANNELS map[string]string `json:"channels"`
	RTM      string            `json:"rtm"`
}

func generateRTMToken(intUid uint32) string {
	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("CERTIFICATE")
	expireTimeInSeconds := uint32(24 * 60 * 60)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	uid := strconv.FormatUint(uint64(intUid), 10)
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtmtokenbuilder.BuildToken(appID, appCertificate, uid, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

	if err != nil {
		log.Printf("Error %+v", err)
	}
	return result
}

func generateRtcToken(uid uint32, channelName string, role rtctokenbuilder.Role) string {
	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("CERTIFICATE")
	tokenExpireTimeInSeconds := uint32(60 * 60 * 24)
	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, role, tokenExpireTimeInSeconds)
	if err != nil {
		log.Printf("Error %+v", err)
	}
	return result
}

func generateARandomUID() uint32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

func getToken(c *gin.Context) {
	uid := generateARandomUID()
	rtmToken := generateRTMToken(uid)
	channels, _ := c.GetQueryArray("channels[]")

	tokensMap := make(map[string]string)
	for _, channel := range channels {
		rtcToken := generateRtcToken(uid, channel, rtctokenbuilder.RolePublisher)
		tokensMap[channel] = rtcToken
	}

	token := Token{
		UID: uid, CHANNELS: tokensMap, RTM: rtmToken,
	}

	c.IndentedJSON(http.StatusOK, token)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.GET("/token", getToken)
	router.Run(":" + port)
}
