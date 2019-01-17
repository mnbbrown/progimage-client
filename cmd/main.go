package main

import (
	"fmt"
	"github.com/mnbbrown/go-progimage/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "progimage",
	Short: "progimage cli",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [path]",
	Short: "Upload an image to progimage",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := pkg.NewClient()
		id, err := client.Upload(args[0], &pkg.UploadParams{})
		if err != nil {
			log.Printf("Failed: %s", err.Error())
			return
		}
		log.Printf("Uploaded to %s", id)
	},
}

func main() {
	rootCmd.AddCommand(uploadCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
