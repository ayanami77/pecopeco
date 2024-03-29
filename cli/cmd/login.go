package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ayanami77/pecopeco-cli/api/factory/user"
	"github.com/ayanami77/pecopeco-cli/auth"
	"github.com/ayanami77/pecopeco-cli/auth/api/userinfo"
	"github.com/ayanami77/pecopeco-cli/auth/secret"
	"github.com/ayanami77/pecopeco-cli/config"
	uiutil "github.com/ayanami77/pecopeco-cli/ui/util"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login pecopeco-cli",
	Long:  "You can login the pecopeco CLI with sth account.",
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func login() {
	factory := user.CreateFactory()
	// 既にログインしていた場合は処理をはじく。
	if config.IsLogin() {
		user, err := factory.FindUser()
		if err != nil {
			uiutil.TextBlue().Println(errorMsg)
			return
		}
		uiutil.TextBlue().Printf(
			"You have already logged in as %s, %s.\nTo log in to another account, please log out of your current account first.\n",
			user.Name,
			user.Email,
		)
		return
	}

	id, secret, scopes, redirectURL, err := secret.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// OAuthによる処理
	oauth := auth.NewOAuth(
		id,
		secret,
		scopes,
		google.Endpoint,
		redirectURL,
	)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if err := oauth.Authorization(ctx); err != nil {
		fmt.Println(err)
		return
	}

	uiutil.TextGreen().Printf("Authentication complete.\n\n")

	sp := uiutil.DefaultSpinner("Logging in...")
	sp.Start()

	userinfo, err := userinfo.Get(ctx, oauth)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ログイン処理
	params := user.LoginParams{
		ID:    userinfo.ID,
		Name:  userinfo.Name,
		Email: userinfo.Email,
	}

	response, err := factory.Login(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	sp.Stop()
	uiutil.TextGreen().Printf("Successfully login as %s, %s\n", response.Name, response.Email)
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
