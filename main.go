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
	UID     uint32 `json:"uid"`
	CHANNEL string `json:"channel"`
	RTC     string `json:"token2"`
	RTM     string `json:"rtm"`
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

func generateRtcToken(int_uid uint32, channelName string, role rtctokenbuilder.Role) string {
	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("CERTIFICATE")
	tokenExpireTimeInSeconds := uint32(60 * 60 * 24)
	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, int_uid, role, tokenExpireTimeInSeconds)
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
	channels, _ := c.GetQueryArray("channels[]")
	tokens := make([]Token, len(channels))
	for idx, channel := range channels {
		rtcToken := generateRtcToken(uid, channel, rtctokenbuilder.RolePublisher)
		rtmToken := generateRTMToken(uid)
		sampleToken := Token{
			UID: uid, RTC: rtcToken, CHANNEL: channel, RTM: rtmToken,
		}
		tokens[idx] = sampleToken
	}

	c.IndentedJSON(http.StatusOK, tokens)
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
