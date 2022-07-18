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
	RTMUID   string            `json:"rtmuid"`
	CHANNELS map[string]string `json:"channels"`
	RTM      string            `json:"rtm"`
}

func generateRTMToken(intUid uint32, uid string) string {
	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("CERTIFICATE")
	expireTimeInSeconds := uint32(24 * 60 * 60)
	currentTimestamp := uint32(time.Now().UTC().Unix())
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

func getNewToken(c *gin.Context) {
	uid := generateARandomUID()
	getToken(c, uid)
}

func getRefreshToken(c *gin.Context) {
	uidStr, success := c.GetQuery("uid")
	uid64, _ := strconv.ParseUint(uidStr, 10, 32)
	uid := uint32(uid64)

	if !success {
		handleGenericError("uid is a required query parameter", c)
	} else if uid == 0 {
		handleGenericError("uid must be a numeric that can be converted to uint32", c)
	} else {
		getToken(c, uid)
	}
}

func handleGenericError(msg string, c *gin.Context) {
	c.JSON(400, map[string]interface{}{"error": msg})
}

func getToken(c *gin.Context, uid uint32) {
	rtmUid := strconv.FormatUint(uint64(uid), 10)
	rtmToken := generateRTMToken(uid, rtmUid)
	channels, success := c.GetQueryArray("channels[]")

	if !success {
		handleGenericError("channels[] is a required query parameter", c)
		return
	}

	tokensMap := make(map[string]string)
	for _, channel := range channels {
		rtcToken := generateRtcToken(uid, channel, rtctokenbuilder.RolePublisher)
		tokensMap[channel] = rtcToken
	}

	token := Token{
		UID: uid, CHANNELS: tokensMap, RTM: rtmToken, RTMUID: rtmUid,
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
	router.GET("/token", getNewToken)
	router.GET("/refreshToken", getRefreshToken)
	router.Run(":" + port)
}
