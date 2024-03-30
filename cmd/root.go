/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"github.com/B87/go-htmx-template/internal"
	"github.com/B87/go-htmx-template/internal/server"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "go-htmx-template",
	Long:  `go-htmx-template is a REST API server.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var config = internal.MustNewConfig()

var db = sqlx.MustConnect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.PGUser, config.PGPassword, "go-htmx"))

var httpServer = server.NewHttpServer(
	server.HTTPServerConfig{
		Host: "localhost",
		Port: config.ServerPort,
		CDN:  server.NewGoogleCloudBucketCDN(config.CDNBucket),
		DB:   db,
	},
)
